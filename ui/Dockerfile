FROM node:24.3 AS build

ARG VITE_VERSION
ARG VITE_BUILD_TIME

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

# Build with default values - will be overridden at runtime
RUN VITE_VERSION=${VITE_VERSION:-dev} \
    VITE_BUILD_TIME=${VITE_BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)} \
    npm run build

FROM nginx:stable-alpine

COPY --from=build /app/dist /usr/share/nginx/html

COPY ./nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Copy startup script
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
