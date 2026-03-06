# Webapp-only Docker image
# Build context is repo root (.)
# Expects .docker-build/webapp-merged/ to contain the pre-built web apps
FROM nginx:stable-alpine

WORKDIR /usr/share/nginx/html

COPY .docker-build/webapp-merged /usr/share/nginx/html
COPY client/web/nginx.conf /etc/nginx/nginx.conf
COPY client/web/entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 80

HEALTHCHECK --interval=30s --start-period=10s --timeout=30s \
  CMD wget --quiet --tries=1 --spider "http://127.0.0.1:80/config.js" || exit 1

ENTRYPOINT ["/entrypoint.sh"]
