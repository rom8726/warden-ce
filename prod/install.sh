#!/bin/bash

# Warden Platform Installer
# Bash version of the Go installer

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration variables
INSTALL_DIR="/opt/warden"
DOCKER_REGISTRY=""
PLATFORM_VERSION="latest"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        print_error "This installer must be run as root (sudo)"
        exit 1
    fi
}

# Function to print welcome message
print_welcome() {
    echo "================================================="
    echo "       Welcome to Warden Platform Installer      "
    echo "================================================="
    echo "This installer will set up the Warden platform on your system."
    echo "It will create necessary directories and configuration files."
    echo
}

# Function to read user input with validation
read_input() {
    local prompt="$1"
    local validation_func="$2"
    local input=""
    
    while true; do
        echo -n "$prompt: "
        read -r input
        input=$(echo "$input" | xargs) # trim whitespace
        
        if [[ -n "$input" ]]; then
            if [[ -n "$validation_func" ]]; then
                if $validation_func "$input"; then
                    break
                fi
            else
                break
            fi
        else
            print_error "Input cannot be empty. Please try again."
        fi
    done
    
    echo "$input"
}

# Function to read yes/no input
read_yes_no() {
    local prompt="$1"
    local input=""
    
    while true; do
        echo -n "$prompt (y/n): "
        read -r input
        input=$(echo "$input" | tr '[:upper:]' '[:lower:]')
        
        if [[ "$input" == "y" || "$input" == "yes" ]]; then
            return 0
        elif [[ "$input" == "n" || "$input" == "no" ]]; then
            return 1
        else
            print_error "Please enter 'y' or 'n'"
        fi
    done
}

# Function to validate email
validate_email() {
    local email="$1"
    if [[ "$email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]] && [[ ! "$email" =~ \.\..*@ ]]; then
        return 0
    else
        print_error "Please enter a valid email address"
        return 1
    fi
}

# Function to generate random string
generate_random_string() {
    local length="$1"
    LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c "$length"
}

# Function to create directories
create_directories() {
    print_info "Creating installation directories..."
    
    local directories=(
        "$INSTALL_DIR"
        "$INSTALL_DIR/nginx/ssl"
        "$INSTALL_DIR/secrets"
    )
    
    for dir in "${directories[@]}"; do
        print_info "Creating directory: $dir"
        mkdir -p "$dir"
        chmod 600 "$dir"
    done
}

# Function to generate SSL certificate
generate_ssl_certificate() {
    local domain="$1"
    local key_path="$INSTALL_DIR/nginx/ssl/nginx_key.pem"
    local cert_path="$INSTALL_DIR/nginx/ssl/nginx_cert.pem"
    
    print_info "Generating self-signed SSL certificate for $domain..."
    
    # Generate private key
    openssl genrsa -out "$key_path" 2048 2>/dev/null
    
    # Generate certificate
    openssl req -new -x509 -key "$key_path" -out "$cert_path" -days 3650 -subj "/C=US/ST=State/L=City/O=Warden/CN=$domain" 2>/dev/null
    
    # Set proper permissions
    chmod 600 "$key_path"
    chmod 644 "$cert_path"
    
    print_success "SSL certificate generated successfully"
}

# Function to render template
render_template() {
    local template_file="$1"
    local output_file="$2"
    local temp_file="$3"
    
    # Create temporary file with template content
    cat > "$temp_file" << 'EOF'
#!/bin/bash
# Template renderer
EOF
    
    # Add template content
    cat "$template_file" >> "$temp_file"
    
    # Make it executable and run
    chmod +x "$temp_file"
    "$temp_file" > "$output_file"
    
    # Clean up
    rm -f "$temp_file"
}

# Function to collect user input
collect_user_input() {
    print_info "=== Platform Configuration ==="
    
    # Get admin email
    ADMIN_EMAIL=$(read_input "Enter administrator email" "validate_email")
    
    # Get domain
    DOMAIN=$(read_input "Enter domain for the platform")
    FRONTEND_URL="https://$DOMAIN"
    
    # Ask about SSL certificate
    if read_yes_no "Do you have an existing SSL certificate for this domain?"; then
        HAS_EXISTING_SSL_CERT=true
        print_info "You will need to place your SSL certificate and key files at:"
        print_info "  - Certificate: $INSTALL_DIR/nginx/ssl/nginx_cert.pem"
        print_info "  - Key: $INSTALL_DIR/nginx/ssl/nginx_key.pem"
        print_info "You will be reminded about this at the end of installation."
    else
        HAS_EXISTING_SSL_CERT=false
        print_info "A self-signed SSL certificate will be generated for you at the end of installation."
    fi
    
    print_info "=== SMTP Server Configuration ==="
    
    # Get SMTP server details
    MAILER_ADDR=$(read_input "Enter SMTP server address (including port)")
    MAILER_USER=$(read_input "Enter SMTP user")
    MAILER_PASSWORD=$(read_input "Enter SMTP password")
    MAILER_FROM=$(read_input "Enter email address for sending emails (from)")
    
    # TLS option
    if read_yes_no "Use TLS for SMTP connection?"; then
        MAILER_USE_TLS=true
    else
        MAILER_USE_TLS=false
    fi
}

# Function to generate secrets
generate_secrets() {
    print_info "Generating secure passwords and keys..."
    
    PG_PASSWORD=$(generate_random_string 12)
    CH_PASSWORD=$(generate_random_string 12)
    SECRET_KEY=$(generate_random_string 32)
    JWT_SECRET_KEY=$(generate_random_string 32)
    ADMIN_TMP_PASSWORD=$(generate_random_string 12)
    SECURE_LINK_MD5=""
    
    print_success "Generated secure passwords and keys"
}

# Function to create platform.env
create_platform_env() {
    local platform_env_file="$INSTALL_DIR/platform.env"
    
    cat > "$platform_env_file" << EOF
DOCKER_REGISTRY=$DOCKER_REGISTRY
DOMAIN=$DOMAIN
DB_PASSWORD=$PG_PASSWORD
CLICKHOUSE_DB_PASSWORD=$CH_PASSWORD
SSL_CERT=nginx_cert.pem
SSL_CERT_KEY=nginx_key.pem
SECURE_LINK_MD5=$SECURE_LINK_MD5
PLATFORM_VERSION=$PLATFORM_VERSION
EOF
    
    print_success "Created $platform_env_file"
}

# Function to create config.env
create_config_env() {
    local config_env_file="$INSTALL_DIR/config.env"
    
    cat > "$config_env_file" << EOF
WARDEN_LOGGER_LEVEL=info

WARDEN_FRONTEND_URL=$FRONTEND_URL

WARDEN_API_SERVER_ADDR=:8080
WARDEN_API_SERVER_READ_TIMEOUT=15s
WARDEN_API_SERVER_WRITE_TIMEOUT=30s
WARDEN_API_SERVER_IDLE_TIMEOUT=60s

WARDEN_TECH_SERVER_ADDR=:8081
WARDEN_TECH_SERVER_READ_TIMEOUT=15s
WARDEN_TECH_SERVER_WRITE_TIMEOUT=30s
WARDEN_TECH_SERVER_IDLE_TIMEOUT=60s

WARDEN_SECRET_KEY=$SECRET_KEY

# JWT
WARDEN_JWT_SECRET_KEY=$JWT_SECRET_KEY
WARDEN_ACCESS_TOKEN_TTL=3h
WARDEN_REFRESH_TOKEN_TTL=168h
WARDEN_RESET_PASSWORD_TTL=8h

# PostgreSQL
WARDEN_POSTGRES_HOST=warden-pgbouncer
WARDEN_POSTGRES_DATABASE=warden
WARDEN_POSTGRES_PASSWORD=$PG_PASSWORD
WARDEN_POSTGRES_PORT=6432
WARDEN_POSTGRES_USER=warden
WARDEN_POSTGRES_MIGRATIONS_DIR=/migrations/postgresql
WARDEN_POSTGRES_MAX_IDLE_CONN_TIME=5m
WARDEN_POSTGRES_MAX_CONNS=20
WARDEN_POSTGRES_CONN_MAX_LIFETIME=10m

# PostgreSQL for migrations (direct connection)
WARDEN_POSTGRES_MIGRATION_HOST=warden-postgresql
WARDEN_POSTGRES_MIGRATION_PORT=5432

# Redis
WARDEN_REDIS_HOST=warden-redis
WARDEN_REDIS_PORT=6379
WARDEN_REDIS_PASSWORD=
WARDEN_REDIS_DB=0

# ClickHouse
WARDEN_CLICKHOUSE_HOST=warden-clickhouse
WARDEN_CLICKHOUSE_PORT=9000
WARDEN_CLICKHOUSE_DATABASE=warden
WARDEN_CLICKHOUSE_USER=default
WARDEN_CLICKHOUSE_PASSWORD=$CH_PASSWORD
WARDEN_CLICKHOUSE_TIMEOUT=10s
WARDEN_CLICKHOUSE_MIGRATIONS_DIR=/migrations/clickhouse

# Kafka
WARDEN_KAFKA_BROKERS=warden-kafka:9092
WARDEN_KAFKA_CLIENT_ID=app-warden
WARDEN_KAFKA_CONSUMER_GROUP_ID=warden-consumer
WARDEN_KAFKA_VERSION=3.9.0
WARDEN_KAFKA_TIMEOUT=10s

# Mailer
WARDEN_MAILER_ADDR=$MAILER_ADDR
WARDEN_MAILER_USER=$MAILER_USER
WARDEN_MAILER_PASSWORD=$MAILER_PASSWORD
WARDEN_MAILER_FROM=$MAILER_FROM
WARDEN_MAILER_ALLOW_INSECURE=false
WARDEN_MAILER_USE_TLS=$MAILER_USE_TLS
WARDEN_MAILER_CERT_FILE="/opt/warden/secrets/mailer_cert.pem"
WARDEN_MAILER_KEY_FILE="/opt/warden/secrets/mailer_key.pem"

# For emails
WARDEN_LOGO_URL=https://warden-project.tech/logo.png

# Envelope consumer LRU cache
WARDEN_CACHE_ENABLED=true
WARDEN_CACHE_RELEASE_CACHE_SIZE=10000
WARDEN_CACHE_ISSUE_CACHE_SIZE=10000
WARDEN_CACHE_ISSUE_RELEASE_CACHE_SIZE=10000

# Ingest server rate limit
WARDEN_RATE_LIMIT_RATE_LIMIT=100
WARDEN_RATE_LIMIT_RPS_WINDOW=10s
WARDEN_RATE_LIMIT_STATS_REFRESH_INTERVAL=1s

# Issue notificator
WARDEN_ISSUE_NOTIFICATOR_WORKER_COUNT=5

# User notificator
WARDEN_USER_NOTIFICATOR_WORKER_COUNT=5

# Admin user
WARDEN_ADMIN_EMAIL=$ADMIN_EMAIL
WARDEN_ADMIN_TMP_PASSWORD=$ADMIN_TMP_PASSWORD
EOF
    
    print_success "Created $config_env_file"
}

# Function to copy docker-compose.yml
copy_docker_compose() {
    local docker_compose_file="$INSTALL_DIR/docker-compose.yml"
    
    cp "$(dirname "$0")/docker-compose.yml" "$docker_compose_file"
    
    print_success "Created $docker_compose_file"
}

# Function to create Makefile
create_makefile() {
    local makefile="$INSTALL_DIR/Makefile"
    
    cat > "$makefile" << 'EOF'
_COMPOSE=docker compose -f docker-compose.yml --project-name warden --env-file platform.env

.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: up
up: ## Up the environment in docker compose
	${_COMPOSE} up -d

.PHONY: down
down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

.PHONY: pull
pull: ## Pull images from remote Docker registry
	${_COMPOSE} pull
EOF
    
    print_success "Created $makefile"
}

# Function to print final information
print_final_info() {
    echo
    print_success "Installation completed successfully!"
    print_info "The platform has been installed in $INSTALL_DIR"
    print_info "You can check and modify settings in $INSTALL_DIR/platform.env and $INSTALL_DIR/config.env"
    print_info "A Makefile has been created in $INSTALL_DIR with commands for starting and stopping the platform"
    
    echo
    print_info "ADMIN LOGIN INFORMATION:"
    print_info "  - Email: $ADMIN_EMAIL"
    print_info "  - Temporary Password: $ADMIN_TMP_PASSWORD"
    print_info "Please use these credentials to log in to the platform. You will be prompted to change the password on first login."
    
    if [[ "$HAS_EXISTING_SSL_CERT" == true ]]; then
        echo
        print_warning "REMINDER: Don't forget to place your SSL certificate and key files at:"
        print_info "  - Certificate: $INSTALL_DIR/nginx/ssl/nginx_cert.pem"
        print_info "  - Key: $INSTALL_DIR/nginx/ssl/nginx_key.pem"
    fi
    
    if [[ "$MAILER_USE_TLS" == true ]]; then
        echo
        print_warning "REMINDER: Don't forget to place your email TLS certificate and key files at:"
        print_info "  - Certificate: $INSTALL_DIR/secrets/mailer_cert.pem"
        print_info "  - Key: $INSTALL_DIR/secrets/mailer_key.pem"
    fi
}

# Main installation function
main() {
    # Check if running as root
    check_root
    
    # Print welcome message
    print_welcome
    
    # Inform about installation directory
    print_info "The platform will be installed in the $INSTALL_DIR directory."
    echo
    
    # Ask for confirmation to proceed
    if ! read_yes_no "Do you want to continue with the installation?"; then
        print_info "Installation cancelled."
        exit 0
    fi
    
    # Collect user input
    collect_user_input
    
    # Generate passwords and other required values
    generate_secrets
    
    # Create required directories
    create_directories
    
    # Create configuration files
    create_platform_env
    create_config_env
    copy_docker_compose
    create_makefile
    
    # Handle SSL certificate based on user's choice
    if [[ "$HAS_EXISTING_SSL_CERT" == false ]]; then
        generate_ssl_certificate "$DOMAIN"
    fi
    
    # Print final information
    print_final_info
}

# Run main function
main "$@"
