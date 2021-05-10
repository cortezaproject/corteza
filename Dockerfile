# run stage
FROM golang:1.16-alpine as run-stage
RUN apk add build-base --no-cache
WORKDIR /corteza
COPY . ./


# build stage
FROM run-stage as build-stage
RUN CGO_ENABLED=1 GOOS="linux" GOARCH="amd64" go build -o corteza-server cmd/corteza/main.go


# deploy stage
FROM golang:1.16-alpine as deploy-stage
ENV STORAGE_PATH "/data"
ENV CORREDOR_ADDR "corredor:80"
ENV HTTP_ADDR "0.0.0.0:80"
WORKDIR /corteza
VOLUME /data
COPY --from=build-stage /corteza/corteza-server ./bin
COPY provision ./provision
HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 \
    CMD curl --silent --fail --fail-early http://127.0.0.1:80/healthcheck || exit 1
EXPOSE 80
ENTRYPOINT ["./corteza-server"]
CMD ["serve-api"]
