# build stage
FROM golang:1.16-alpine as build-stage

ENV GOFLAGS='-mod=readonly'
ENV BUILD_OS=linux
ENV BUILD_ARCH=amd64
ENV BUILD_VERSION=latest

RUN apk add build-base --no-cache

WORKDIR /corteza

COPY . ./

RUN make release-clean release


# deploy stage
FROM alpine:3

ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
ENV HTTP_WEBAPP_ENABLED "false"
ENV PATH "/corteza/bin:${PATH}"

WORKDIR /corteza

VOLUME /data

COPY --from=build-stage /corteza/build/pkg/corteza-server ./

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1

EXPOSE 80

ENTRYPOINT ["./bin/corteza-server"]

CMD ["serve-api"]
