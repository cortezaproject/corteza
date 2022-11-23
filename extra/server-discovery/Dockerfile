# deploy stage
FROM ubuntu:20.04

RUN apt-get -y update \
 && apt-get -y install \
    ca-certificates \
    curl \
 && rm -rf /var/lib/apt/lists/*

ARG VERSION=2022.3.1-discovery
ARG SERVER_VERSION=${VERSION}
ARG CORTEZA_SERVER_PATH=https://releases.cortezaproject.org/files/corteza-server-discovery-${SERVER_VERSION}-linux-amd64.tar.gz
RUN mkdir /tmp/server
ADD $CORTEZA_SERVER_PATH /tmp/server

VOLUME /data

RUN tar -zxvf "/tmp/server/$(basename $CORTEZA_SERVER_PATH)" -C / && \
    rm -rf "/tmp/server" && \
    mv /corteza-server /corteza

WORKDIR /corteza

#HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
#    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

ENV STORAGE_PATH "/data"
ENV HTTP_ADDR "0.0.0.0:80"
ENV CORREDOR_ENABLED "false"
ENV HTTP_WEBAPP_ENABLED "false"
ENV PATH "/corteza/bin:${PATH}"

EXPOSE 80

ENTRYPOINT ["./bin/corteza-server"]

#CMD ["serve-api"]
