# All-in-one Docker image (server + webapp)
# Expects:
#   .docker-build/server/        - pre-built server package
#   .docker-build/webapp-merged/ - pre-built web apps
FROM ubuntu:22.04

ARG SASS_VERSION=1.85.1
ARG SASS_URL=https://github.com/sass/dart-sass/releases/download/${SASS_VERSION}/dart-sass-${SASS_VERSION}-linux-x64.tar.gz

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /opt
RUN curl -sOL $SASS_URL && tar -xzf dart-sass-${SASS_VERSION}-linux-x64.tar.gz && rm dart-sass-*.tar.gz

RUN mkdir -p /human/webapp
COPY .docker-build/server /human
COPY .docker-build/webapp-merged /human/webapp

WORKDIR /human

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
ENV HTTP_WEBAPP_ENABLED "true"
ENV HTTP_WEBAPP_BASE_DIR "/human/webapp"
ENV PATH "/opt/dart-sass:/human/bin:${PATH}"

VOLUME /data

EXPOSE 80

ENTRYPOINT ["./bin/human-server"]

CMD ["serve-api"]
