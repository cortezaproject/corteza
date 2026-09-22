# Build context is the repository root:
#   docker build -f corredor/Dockerfile -t human-corredor .
#
# The pnpm workspace root is the image root (/), so the package keeps its
# /corredor path and the protobuf definitions sit at /def/protobuf.
FROM node:22-slim

ENV PATH=/corredor/node_modules/.bin:$PATH

# GRPC's env-vars
# https://github.com/grpc/grpc/blob/master/doc/environment_variables.md
ENV GRPC_VERBOSITY=ERROR

ENV CORREDOR_ADDR=0.0.0.0:80
ENV CORREDOR_ENVIRONMENT=prod
ENV CORREDOR_LOG_PRETTY=false
ENV CORREDOR_LOG_LEVEL=info
ENV CORREDOR_LOG_ENABLED=true
ENV CORREDOR_SERVER_CERTIFICATES_ENABLED=false
ENV CORREDOR_SCRIPTS_AUTO_UPDATE_DEPENDENCIES=false
ENV CORREDOR_EXT_SERVER_SCRIPTS_WATCH=false
ENV CORREDOR_EXT_CLIENT_SCRIPTS_WATCH=false
ENV CORREDOR_EXT_SEARCH_PATHS=/corredor/usr/*:/corredor/usr
# This assumes that container will be part of the docker-compose setup were Human is
# ran under "server" service and can be directly accessed via internal docker network
ENV CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE="http://server/api/{service}"

# git for extensions with git dependencies, netcat for the healthcheck
RUN apt-get update \
 && apt-get -y install --no-install-recommends git netcat-openbsd \
 && rm -rf /var/lib/apt/lists/*

ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
RUN corepack enable

WORKDIR /

# Workspace manifests, then the packages Corredor needs
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml ./
COPY lib/js lib/js
COPY def/protobuf def/protobuf
COPY corredor corredor

RUN pnpm install --frozen-lockfile --filter @planetcrust/human-corredor...

WORKDIR /corredor

# Client & server scripts location for user scripts & extensions
VOLUME /corredor/usr

# TLS certificates should be placed here if CORREDOR_SERVER_CERTIFICATES_ENABLED
VOLUME /corredor/certs

HEALTHCHECK --interval=30s --start-period=1m --timeout=30s --retries=3 CMD nc -z -v localhost 80

CMD ["tsx", "src/server.ts"]
