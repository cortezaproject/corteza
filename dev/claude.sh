#!/usr/bin/env bash
#
# Launch Claude Code for THIS checkout.
#
# The `human-local` MCP server (Human's own /api/mcp, ~150 configurator tools)
# is opt-in: `HUMAN_MCP=1` or a leading `--human-mcp` adds it, with a freshly
# minted token. Without it a session reaches the dev server through
# dev/agent/api.sh and mcp.py, which authenticate per call and record what they
# create in the cleanup ledger.
#
# `.mcp.json` can interpolate an environment variable but cannot run a command
# to produce one, so both the URL and the token have to exist before Claude
# Code starts. Exporting them here rather than from a shell profile keeps them
# scoped to the checkout you launched from — stack.sh resolves the API port
# from this checkout's own server/.env, so a worktree reaches its own server
# with its own server's token.
#
# The URL matters as much as the token: the config's default is the primary's
# 1043, and inside a worktree that is a live server on a different database,
# holding a token the shared secret makes valid. Calls then succeed against the
# wrong stack, and the session reads and writes the primary's data believing it
# is isolated.
#
#   dev/claude.sh [--human-mcp] [claude flags...]
#   make claude [MCP=1] -- [claude flags...]
#
# A dev server that is down is not a reason to refuse: Claude Code is useful
# without the MCP server, and every other route to Human re-authenticates per
# call anyway.
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=dev/agent/stack.sh
source "$HERE/agent/stack.sh"
want_mcp="${HUMAN_MCP:-}"
if [[ "${1:-}" == "--human-mcp" ]]; then
  want_mcp=1
  shift
fi

mcp_args=()
if [[ -n "$want_mcp" ]]; then
  export HUMAN_MCP_URL="$HUMAN_API/mcp"
  if tok="$("$HERE/agent/token.sh" 2>/dev/null)" && [[ -n "$tok" ]]; then
    export HUMAN_MCP_TOKEN="$tok"
  else
    echo "note: no token minted (dev server down?) — human-local MCP will not connect" >&2
    echo "      (it would have used $HUMAN_MCP_URL)" >&2
  fi
  mcp_args=(--mcp-config "$HERE/human-local.mcp.json")
fi

exec "${CLAUDE_BIN:-claude}" ${mcp_args[@]+"${mcp_args[@]}"} "$@"
