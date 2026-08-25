#!/usr/bin/env bash
#
# Launch Claude Code with a token for THIS checkout's Human MCP server.
#
# `.mcp.json` can interpolate an environment variable into the Authorization
# header but cannot run a command to produce one, so the token has to exist
# before Claude Code starts. Minting it here rather than in a shell profile
# keeps it scoped to the checkout you launched from — token.sh resolves the API
# port through stack.sh, so a worktree gets its own server's token.
#
#   dev/claude.sh [claude flags...]
#   make claude -- [claude flags...]
#
# A dev server that is down is not a reason to refuse: Claude Code is useful
# without the MCP server, and every other route to Human re-authenticates per
# call anyway.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if tok="$("$HERE/agent/token.sh" 2>/dev/null)" && [[ -n "$tok" ]]; then
  export HUMAN_MCP_TOKEN="$tok"
else
  echo "note: no token minted (dev server down?) — human-local MCP will not connect" >&2
fi

exec claude "$@"
