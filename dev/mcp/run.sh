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
stale() {
	[[ ! -x "${bin}" ]] && return 0
	[[ -n "$(find "${here}" -name '*.go' -not -path '*/vendor/*' -newer "${bin}" -print -quit 2>/dev/null)" ]]
}

# The shared packages are vendored, so editing server/pkg/mcpkit does not reach
# this binary until they are re-vendored. Without this check a change there is
# built, tested and apparently working on the server side while the developer
# MCP silently keeps running the old copy.
revendored=0
shared="${here}/../../server/pkg/mcpkit"
if [[ -d "${shared}" ]] && [[ -n "$(find "${shared}" -name '*.go' -newer "${here}/vendor/modules.txt" -print -quit 2>/dev/null)" ]]; then
	(cd "${here}" && GOFLAGS=-mod=mod go mod vendor) 1>&2
	revendored=1
fi

# A refreshed vendor tree always means a rebuild: stale() deliberately ignores
# vendor/ so that routine module churn does not force one, which leaves this as
# the only thing that notices the shared packages moved underneath us.
if [[ "${revendored}" == 1 ]] || stale; then
	(cd "${here}" && go build -o "${bin}" .) 1>&2
fi

exec "${bin}" "$@"
