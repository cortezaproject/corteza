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

# The shared packages come from ../../server through the replace directive and
# are compiled from source, so a change there reaches this binary as soon as it
# is rebuilt. They are watched alongside this module's own sources: without that
# a change to server/pkg/mcpkit is built, tested and apparently working on the
# server side while the developer MCP silently keeps running the old copy.
shared="${here}/../../server/pkg/mcpkit"

# Rebuild when any source is newer than the binary. Cheap to check, and it means
# a developer editing a tool does not have to remember to rebuild — restarting
# the MCP client is enough.
stale() {
	[[ ! -x "${bin}" ]] && return 0

	local roots=("${here}")
	[[ -d "${shared}" ]] && roots+=("${shared}")

	[[ -n "$(find "${roots[@]}" -name '*.go' -newer "${bin}" -print -quit 2>/dev/null)" ]]
}

if stale; then
	(cd "${here}" && go build -o "${bin}" .) 1>&2
fi

exec "${bin}" "$@"
