#!/usr/bin/env bash
#
# Where THIS checkout's dev stack listens.
#
# One resolver for every script that has to reach the local server or webapp.
# A worktree runs its own server, against its own database, on its own ports —
# so a script carrying the literal 1043 talks to the primary instead. The token
# is shared and the user IDs match, so the call succeeds and the check reports
# a pass against code the worktree never ran.
#
# Precedence, highest first:
#   1. the environment          — an explicit HUMAN_API/HUMAN_WEBAPP always wins
#   2. this checkout's own files — server/.env, .env.e2e, the worktree registry
#   3. the primary's ports       — 1043 / 3001 / 5173
#
#   source dev/agent/stack.sh      sets the variables
#   eval "$(dev/agent/stack.sh)"   the same, from a subshell
#   dev/agent/stack.sh             prints KEY=value for a non-bash caller
#
# Which ports a slot gets is worktree.sh's business; this only reads back what
# worktree.sh (or dev/setup.sh) wrote.

STACK_ROOT="${STACK_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"

# The LAST uncommented assignment wins, because that is the one the server
# gets: godotenv parses the whole file into a map before applying it, so a key
# set twice takes its final value. Reading the first match instead reports a
# port nothing is listening on the moment someone appends an override.
_stack_last() { sed -nE "$2" "$1" 2>/dev/null | tail -1; }

# HTTP_ADDR is a bind address — ":1043", "127.0.0.1:1043", quoted or not.
_stack_api_port() {
  _stack_last "$1/server/.env" \
    's/^[[:space:]]*HTTP_ADDR=["'\'']?[^:"'\'' ]*:([0-9]+).*/\1/p'
}

_stack_api_base() {
  _stack_last "$1/server/.env" \
    's#^[[:space:]]*HTTP_API_BASE_URL=["'\'']?(/[^"'\'' ]*).*#\1#p'
}

# .env.e2e first: every worktree gets one naming its own vite port, and the
# literal in playwright.config.ts is that file's fallback rather than its value.
# The registry sits between them for the one moment .env.e2e does not exist yet
# — dev/setup.sh generating it — where the config literal would hand a lane the
# primary's webapp.
_stack_webapp() { # _stack_webapp ROOT API_PORT
  local unify="$1/client/web/unify" url port
  url="$(_stack_last "$unify/.env.e2e" \
    's#^[[:space:]]*E2E_BASE_URL=["'\'']?(https?://[^"'\'' ]+).*#\1#p')"
  [[ -n "$url" ]] && {
    echo "$url"
    return
  }
  port="$(_stack_registry "$1" "$2" vite)"
  [[ -n "$port" ]] && {
    echo "http://localhost:$port"
    return
  }
  _stack_last "$unify/playwright.config.ts" \
    "s#.*E2E_BASE_URL \|\| '([^']+)'.*#\1#p"
}

# gin's proxy port never reaches an env file — worktree.sh passes it on the
# command line — so the worktree registry is the only record of it, and of the
# vite port before .env.e2e exists to name it.
#
# The entry has to agree with server/.env about the API port before anything in
# it is worth reading: a checkout repointed after it was registered would
# otherwise be handed ports in front of somebody else's stack.
_stack_registry() { # _stack_registry ROOT API_PORT KEY
  local f
  for f in "$1"/dev/agent/.state/worktrees/*.json; do
    [[ -f "$f" ]] || continue
    grep -q "\"path\": \"$1\"" "$f" || continue
    grep -q "\"api\": ${2:-0}," "$f" || continue
    sed -nE "s/.*\"$3\":[[:space:]]*([0-9]+).*/\\1/p" "$f" | head -1
    return
  done
}

_stack_port="$(_stack_api_port "$STACK_ROOT")"
_stack_base="$(_stack_api_base "$STACK_ROOT")"
_stack_gin="$(_stack_registry "$STACK_ROOT" "$_stack_port" gin)"

# Falling back rather than failing keeps every script usable on a checkout that
# has not been set up yet — `make setup` itself runs before server/.env exists.
HUMAN_API="${HUMAN_API:-http://localhost:${_stack_port:-1043}${_stack_base:-/api}}"
HUMAN_BASE="${HUMAN_API%"${_stack_base:-/api}"}"
HUMAN_AUTH="${HUMAN_AUTH:-$HUMAN_BASE/auth}"
HUMAN_WEBAPP="${HUMAN_WEBAPP:-$(_stack_webapp "$STACK_ROOT" "$_stack_port")}"
HUMAN_WEBAPP="${HUMAN_WEBAPP:-http://localhost:5173}"
HUMAN_GIN="${HUMAN_GIN:-http://localhost:${_stack_gin:-3001}}"
unset _stack_port _stack_base _stack_gin

# Executed rather than sourced: print what a non-bash caller should read.
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  cat <<EOF
HUMAN_API=$HUMAN_API
HUMAN_BASE=$HUMAN_BASE
HUMAN_AUTH=$HUMAN_AUTH
HUMAN_WEBAPP=$HUMAN_WEBAPP
HUMAN_GIN=$HUMAN_GIN
EOF
fi
