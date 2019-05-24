#!/usr/bin/env sh

# This file is intended to be used from Dockerfile,
# presumably from cortezaproject/corteza-server-builder.

set -eu

BUILD_TIME=${BUILD_TIME:-$(date +%FT%T%z)}
GIT_TAG=${GIT_TAG:-$(git describe --always --tags)}
APP=${1}
DST=${2:-"/bin/corteza-server-${APP}"}

LDFLAGS=""
LDFLAGS="${LDFLAGS} -X github.com/cortezaproject/corteza-server/internal/version.BuildTime=${BUILD_TIME}"
LDFLAGS="${LDFLAGS} -X github.com/cortezaproject/corteza-server/internal/version.Version=${GIT_TAG}"

go build -ldflags "${LDFLAGS}" -o $DST ./cmd/$APP/*.go
