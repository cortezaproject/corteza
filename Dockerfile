# build-stage
FROM alpine:3 as build-stage

# use docker build --build-arg VERSION=2021.9.0 .
ARG VERSION=2022.9.0
ARG SASS_VERSION=1.69.5
ARG SERVER_VERSION=${VERSION}
ARG WEBAPP_VERSION=${VERSION}
ARG CORTEZA_SERVER_PATH=https://releases.cortezaproject.org/files/corteza-server-${SERVER_VERSION}-linux-amd64.tar.gz
ARG CORTEZA_WEBAPP_PATH=https://releases.cortezaproject.org/files/corteza-webapp-${WEBAPP_VERSION}.tar.gz
ARG SASS_URL=https://github.com/sass/dart-sass/releases/download/${SASS_VERSION}/dart-sass-${SASS_VERSION}-linux-x64.tar.gz

RUN mkdir /tmp/server
RUN mkdir /tmp/webapp

ADD $CORTEZA_SERVER_PATH /tmp/server
ADD $CORTEZA_WEBAPP_PATH /tmp/webapp

RUN apk update && apk add --no-cache file
RUN apk add curl

RUN file "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" | grep -q 'gzip' && \
    tar zxvf "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" -C / || \
    cp -a "/tmp/server" /

RUN mv /corteza-server /corteza

WORKDIR /corteza

RUN rm -rf /corteza/webapp

RUN file "/tmp/webapp/$(basename $CORTEZA_WEBAPP_PATH)" | grep -q 'gzip' && \
    mkdir /corteza/webapp && tar zxvf "/tmp/webapp/$(basename $CORTEZA_WEBAPP_PATH)" -C /corteza/webapp || \
    cp -a "/tmp/webapp" /corteza/webapp

WORKDIR /tmp

RUN curl -sOL $SASS_URL
RUN tar -xzf dart-sass-${SASS_VERSION}-linux-x64.tar.gz

# deploy-stage
FROM ubuntu:20.04

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
ENV HTTP_WEBAPP_ENABLED "true"
ENV HTTP_WEBAPP_BASE_DIR "/corteza/webapp"
ENV PATH "/opt/dart-sass:/corteza/bin:${PATH}"

WORKDIR /corteza

VOLUME /data

COPY --from=build-stage /corteza ./
COPY --from=build-stage /tmp/dart-sass /opt/dart-sass

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

EXPOSE 80

ENTRYPOINT ["./bin/corteza-server"]

CMD ["serve-api"]
