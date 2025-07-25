FROM nginx:stable-alpine

# Build arguments for versioning
ARG TOOL_VERSION
ARG TOOL_BUILD_TIME

RUN apk add --no-cache gettext bash

# Set version labels
LABEL version="${TOOL_VERSION}"
LABEL build-time="${TOOL_BUILD_TIME}"

COPY etc/nginx /etc/nginx/
COPY entrypoint.sh /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
