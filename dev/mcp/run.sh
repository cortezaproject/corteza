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

# A failed rebuild keeps the previous binary serving: the client gets the tools
# it had, and the compiler output is on stderr and in build/build.log. Without
# a previous binary there is nothing to fall back to.
if stale; then
	log="${here}/build/build.log"
	if (cd "${here}" && go build -o "${bin}.new" .) 2>"${log}"; then
		mv -f "${bin}.new" "${bin}"
	else
		cat "${log}" >&2
		rm -f "${bin}.new"
		if [[ -x "${bin}" ]]; then
			echo "dev-mcp: rebuild failed; serving the previous binary from $(date -r "${bin}" '+%F %T')" >&2
		else
			echo "dev-mcp: build failed and no previous binary at ${bin}" >&2
			exit 1
		fi
	fi
fi

exec "${bin}" "$@"
