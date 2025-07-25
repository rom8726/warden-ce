FROM postgres:16

# Build arguments for versioning
ARG TOOL_VERSION
ARG TOOL_BUILD_TIME

# Install PgBouncer and gettext-base for envsubst
RUN apt-get update && apt-get install -y \
    pgbouncer \
    postgresql-client \
    gettext-base \
    && rm -rf /var/lib/apt/lists/*

# Set version labels
LABEL version="${TOOL_VERSION}"
LABEL build-time="${TOOL_BUILD_TIME}"

# Copy entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Create pgbouncer user
RUN useradd -r -s /bin/false pgbouncer

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["pgbouncer", "-u", "pgbouncer", "/etc/pgbouncer/pgbouncer.ini"]
