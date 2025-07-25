#!/usr/bin/env bash

NGINX_PATH=/etc/nginx

# Check required environment variables
if [ -z "$DOMAIN" ]; then
    echo "Error: DOMAIN environment variable is required"
    exit 1
fi

if [ -z "$SSL_CERT" ]; then
    echo "Error: SSL_CERT environment variable is required"
    exit 1
fi

if [ -z "$SSL_CERT_KEY" ]; then
    echo "Error: SSL_CERT_KEY environment variable is required"
    exit 1
fi

# shellcheck disable=SC2016,SC2044

for conflist in $(find $NGINX_PATH/templates -name "*.template" -exec basename {} \;); do
    envsubst '${DOMAIN},${BACKEND_URL},${FRONTEND_URL},${INGEST_SERVER_URL},${PG_ACL_LIST},${SSL_CERT},${SSL_CERT_KEY},${SECURE_LINK_MD5}' < "$NGINX_PATH"/templates/"$conflist" > "$NGINX_PATH"/conf.d/"${conflist%.*}".conf
done

nginx -g "daemon off;"
