FROM golang:1.10-alpine AS builder

WORKDIR /go/src/github.com/crusttech/crust

ENV CGO_ENABLED=0

COPY . .

RUN mkdir /build; \
    find cmd -type d -mindepth 1 -maxdepth 1 | while read CMDPATH; do \
        go build -o /build/$(basename ${CMDPATH}) ${CMDPATH}/*.go; \
    done;

FROM alpine:3.7

ENV PATH="/crust:{$PATH}"
WORKDIR /crust

# @todo copy crm/types, migrations

COPY --from=builder ./build/* /crust/

CMD ["/bin/echo", "Run of the crust commands: sam, crm"]