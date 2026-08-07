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

# --- created-resource ledger -------------------------------------------------
#
# Everything the toolkit creates on the dev server is recorded here, and
# cleanup.sh deletes ONLY what it finds in this file. Nothing else is ever a
# candidate: not a slug pattern, not a label, not an author. A marker carried
# in the data can be renamed, cleared by an update that omits the field, or
# imitated — a ledger cannot, and if it is lost the worst case is that scratch
# survives rather than that something real is deleted.
#
# .state/ is gitignored, so the ledger is local to this machine.
LEDGER="$STATE_DIR/created.jsonl"

# The session that created a resource, so a session can clean up after itself
# without touching a concurrent one's work. Falls back to a constant when the
# harness does not export one — those entries are still cleanable with --all.
AGENT_SESSION="${CLAUDE_CODE_SESSION_ID:-unknown}"

# ledger_record KIND ID SLUG — append one created resource.
ledger_record() {
  python3 -c '
import json, os, sys
kind, rid, slug, session, path = sys.argv[1:6]
with open(path, "a") as fh:
    fh.write(json.dumps({
        "kind": kind, "id": rid, "slug": slug,
        "session": session, "ts": __import__("time").strftime("%Y-%m-%dT%H:%M:%S%z"),
    }) + "\n")
' "$1" "$2" "${3:-}" "$AGENT_SESSION" "$LEDGER"
}
