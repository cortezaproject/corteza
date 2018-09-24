## builder image

FROM golang:1.10-alpine AS builder

WORKDIR /go/src/github.com/crusttech/crust

ARG APP
ENV APP=$APP
ENV CGO_ENABLED=0

COPY . .

RUN mkdir /build; go build -o /build/$APP cmd/$APP/*.go


## target image

FROM alpine:3.7

ENV PATH="/crust:{$PATH}"
WORKDIR /crust

# @todo copy crm/types, migrations

COPY --from=builder /build/* /crust/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG APP
ENV APP=$APP

ENTRYPOINT /crust/$APP
