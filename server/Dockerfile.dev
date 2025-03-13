# deploy stage
FROM ubuntu:20.04

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

ARG VERSION=2022.9
ARG SASS_VERSION=1.69.5
ARG SERVER_VERSION=${VERSION}
ARG CORTEZA_SERVER_PATH=https://releases.cortezaproject.org/files/corteza-server-${SERVER_VERSION}-linux-amd64.tar.gz
ARG SASS_URL=https://github.com/sass/dart-sass/releases/download/${SASS_VERSION}/dart-sass-${SASS_VERSION}-linux-x64.tar.gz

RUN mkdir /tmp/server
ADD $CORTEZA_SERVER_PATH /tmp/server

VOLUME /data

RUN tar -zxvf "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" -C / && \
    rm -rf "/tmp/server" && \
    mv /corteza-server /corteza

WORKDIR /opt

RUN curl -sOL $SASS_URL
RUN tar -xzf dart-sass-${SASS_VERSION}-linux-x64.tar.gz

WORKDIR /corteza

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
ENV HTTP_WEBAPP_ENABLED "false"
ENV PATH "/opt/dart-sass:/corteza/bin:${PATH}"

EXPOSE 80

ENTRYPOINT ["./bin/corteza-server"]

CMD ["serve-api"]
