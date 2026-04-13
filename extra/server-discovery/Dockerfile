# Discovery server Docker image
# Expects:
#   .docker-build/server/ - pre-built discovery server package
FROM ubuntu:22.04

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /human

COPY .docker-build/server /human

VOLUME /data

WORKDIR /human

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

ENV STORAGE_PATH "/data"
ENV HTTP_ADDR "0.0.0.0:80"
ENV CORREDOR_ENABLED "false"
ENV HTTP_WEBAPP_ENABLED "false"
ENV PATH "/human/bin:${PATH}"

EXPOSE 80

ENTRYPOINT ["./bin/human-server"]
