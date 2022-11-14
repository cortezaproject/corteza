# build-stage
FROM alpine:3 as build-stage

ARG VERSION=2022.9.0
ARG CORTEZA_CLIENT_WEB_AIO_PKG=https://releases.cortezaproject.org/files/corteza-webapp-${VERSION}.tar.gz

# Download package
ADD $CORTEZA_CLIENT_WEB_AIO_PKG /tmp

RUN apk update && apk add --no-cache file

RUN rm -rf /html; mkdir /html
WORKDIR /html


RUN file "/tmp/$(basename $CORTEZA_CLIENT_WEB_AIO_PKG)" | grep -q 'gzip' && \
    tar zxvf "/tmp/$(basename $CORTEZA_CLIENT_WEB_AIO_PKG)" -C /html

# deploy-stage
FROM nginx:stable-alpine

WORKDIR /usr/share/nginx/html

COPY --from=build-stage /html /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
#COPY CONTRIBUTING.* DCO LICENSE README.* ./
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 80

HEALTHCHECK --interval=30s --start-period=10s --timeout=30s \
  CMD wget --quiet --tries=1 --spider "http://127.0.0.1:80/config.js" || exit 1

ENTRYPOINT ["/entrypoint.sh"]
