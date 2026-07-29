#!/usr/bin/env bash
# Shared config for the agent dev toolkit. Sourced by the other scripts.
# The toolkit is LOCAL-ONLY: it refuses to talk to non-localhost API hosts.

AGENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$AGENT_DIR/../.." && pwd)"
SERVER_DIR="$REPO_DIR/server"
STATE_DIR="$AGENT_DIR/.state"

HUMAN_API="${HUMAN_API:-http://localhost:1043/api}"
HUMAN_AUTH="${HUMAN_AUTH:-${HUMAN_API%/api}/auth}"
HUMAN_BASE="${HUMAN_API%/api}"

AGENT_EMAIL="agent@local.dev"
AGENT_HANDLE="agent-dev"
AGENT_CLIENT="dev-agent"
AGENT_SCOPE="profile api"

case "$HUMAN_API" in
  http://localhost[:/]* | http://127.0.0.1[:/]* | "http://[::1]"[:/]*) ;;
  *)
    echo "agent toolkit is local-only; refusing to touch $HUMAN_API" >&2
    exit 1
    ;;
esac

mkdir -p "$STATE_DIR"

command -v python3 >/dev/null || {
  echo "agent toolkit needs python3 (for JSON handling)" >&2
  exit 1
}

# json_get DOTTED.PATH — read JSON on stdin, print value at path; exit 1 if absent.
json_get() {
  python3 -c '
import json, sys
path = sys.argv[1].split(".")
try:
    d = json.load(sys.stdin)
    for p in path:
        d = d[int(p)] if isinstance(d, list) else d[p]
except Exception:
    sys.exit(1)
if d is None:
    sys.exit(1)
print(d if not isinstance(d, (dict, list)) else json.dumps(d))
' "$1"
}

server_bin() {
  if [[ -x "$SERVER_DIR/build/gin-bin" ]]; then
    echo "$SERVER_DIR/build/gin-bin"
  else
    ls -t "$SERVER_DIR"/build/human-server-* 2>/dev/null | head -1
  fi
}

# Run the server binary as a CLI against the local dev DB (loads server/.env).
server_cli() {
  local bin
  bin="$(server_bin)"
  if [[ -z "$bin" ]]; then
    echo "no server binary in server/build — run 'cd server && make watch' (or 'make build')" >&2
    return 1
  fi
  (cd "$SERVER_DIR" && "$bin" --env-file .env "$@")
}
