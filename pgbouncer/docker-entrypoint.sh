#!/bin/bash

set -e

# Default values
PGBOUNCER_USER=${PGBOUNCER_USER:-"postgres"}
PGBOUNCER_PASSWORD=${PGBOUNCER_PASSWORD:-"password"}
PGBOUNCER_DB_USER=${PGBOUNCER_DB_USER:-"postgres"}
PGBOUNCER_DB_PASSWORD=${PGBOUNCER_DB_PASSWORD:-"password"}

# Function to get password hash from PostgreSQL
get_password_hash() {
    local host=$1
    local port=$2
    local user=$3
    local password=$4
    local db=$5
    
    # Get the actual password hash from PostgreSQL
    # We need to connect as a superuser or the user itself to read pg_authid
    HASHED_PASSWORD=$(PGPASSWORD="$password" psql -h "$host" -p "$port" -U "$user" -d "$db" -t -c "SELECT rolpassword FROM pg_authid WHERE rolname = '$user';" 2>/dev/null | tr -d '[:space:]')
    
    if [ $? -eq 0 ] && [ -n "$HASHED_PASSWORD" ]; then
        # Return the hash in PgBouncer format: "username" "hash"
        echo "\"$user\" \"$HASHED_PASSWORD\""
    else
        echo "Failed to get password hash from PostgreSQL" >&2
        exit 1
    fi
}

# Function to generate pgbouncer.ini from template
generate_config() {
    local template_file="/etc/pgbouncer/pgbouncer.ini.template"
    local config_file="/etc/pgbouncer/pgbouncer.ini"
    
    if [ -f "$template_file" ]; then
        # Use envsubst to substitute environment variables
        envsubst < "$template_file" > "$config_file"
    else
        # Generate default config if no template
        cat > "$config_file" << EOF
[databases]
* = host=${DB_HOST} port=${DB_PORT} dbname=${DB_NAME}

[pgbouncer]
listen_addr = 0.0.0.0
listen_port = ${LISTEN_PORT}
auth_type = scram-sha-256
auth_file = /etc/pgbouncer/userlist.txt
pool_mode = transaction
max_client_conn = ${MAX_CLIENT_CONN}
default_pool_size = ${DEFAULT_POOL_SIZE}
min_pool_size = ${MIN_POOL_SIZE}
reserve_pool_size = ${RESERVE_POOL_SIZE}
reserve_pool_timeout = 5
max_db_connections = ${MAX_DB_CONNECTIONS}
max_user_connections = ${MAX_USER_CONNECTIONS}
server_reset_query = DISCARD ALL
server_check_query = select 1
server_check_delay = 30
ignore_startup_parameters = ${IGNORE_STARTUP_PARAMETERS}
idle_transaction_timeout = 0
tcp_keepalive = 1
tcp_keepidle = 1
tcp_keepintvl = 1
tcp_keepcnt = 5
log_connections = 1
log_disconnections = 1
log_pooler_errors = 1
stats_period = 60
verbose = ${VERBOSE}
EOF
    fi
}

# Main execution
echo "Starting PgBouncer setup..."

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; do
    echo "PostgreSQL is not ready yet, waiting..."
    sleep 2
done

echo "PostgreSQL is ready!"

# Generate userlist.txt with password hash
echo "Generating userlist.txt..."
get_password_hash "$DB_HOST" "$DB_PORT" "$DB_USER" "$DB_PASSWORD" "$DB_NAME" > /etc/pgbouncer/userlist.txt

echo "Userlist.txt generated:"
cat /etc/pgbouncer/userlist.txt

# Generate pgbouncer.ini
echo "Generating pgbouncer.ini..."
generate_config

echo "PgBouncer configuration generated:"
cat /etc/pgbouncer/pgbouncer.ini

# Set proper permissions
chown -R pgbouncer:pgbouncer /etc/pgbouncer/
chmod 600 /etc/pgbouncer/userlist.txt
chmod 644 /etc/pgbouncer/pgbouncer.ini

echo "Starting PgBouncer..."
exec "$@"
