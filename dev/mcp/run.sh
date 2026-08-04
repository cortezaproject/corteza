#!/usr/bin/env bash
#
# Launcher for the developer MCP.
#
# The client speaks JSON-RPC over stdout, so nothing but protocol may be written
# there: `go build` chatter, module downloads and compiler warnings all go to
# stderr, and the binary is built to a file rather than run through `go run`,
# which would put its own diagnostics in the way and pay compile cost on every
# start.
set -euo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bin="${here}/build/dev-mcp"

mkdir -p "${here}/build"

# Rebuild when any source is newer than the binary. Cheap to check, and it means
# a developer editing a tool does not have to remember to rebuild — restarting
# the MCP client is enough.
if [[ ! -x "${bin}" ]] || [[ -n "$(find "${here}" -name '*.go' -newer "${bin}" -print -quit 2>/dev/null)" ]]; then
	(cd "${here}" && go build -o "${bin}" .) 1>&2
fi

exec "${bin}" "$@"
