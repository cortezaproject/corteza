#!/usr/bin/env bash
# Delete the namespaces THIS SESSION created, and nothing else.
#
# Usage: cleanup.sh [--all | --session ID] [--purge]
#
# The only deletion candidates are entries in .state/created.jsonl, written by
# api.sh and mcp.py as they create things. A session cleans up after itself,
# which is what makes this safe to run while real data is on the server: a
# namespace this session never created is not a candidate, whoever made it and
# whatever it is called.
#
# Deliberately NOT used as signals: slug prefixes (renameable), labels (any
# update that omits the field silently clears them) and authorship (an agent
# creating something because it was asked to does not make it disposable).
#
#   --all          every session's entries, not just this one's
#   --session ID   one specific session's entries
#   --purge        additionally HARD-delete all soft-deleted namespaces via the
#                  dev-only server CLI. This is NOT ledger-scoped — it hits
#                  every soft-deleted namespace on the server, including ones
#                  someone else deleted, so it is opt-in.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

purge=0
scope="$AGENT_SESSION"

while (($#)); do
  case "$1" in
    --purge) purge=1 ;;
    --all) scope="" ;;
    --session)
      shift
      scope="${1:?--session needs an ID}"
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 1
      ;;
  esac
  shift
done

if [[ ! -f "$LEDGER" ]]; then
  echo "nothing recorded as created (no $LEDGER)"
  exit 0
fi

token="$("$AGENT_DIR/token.sh")"
deleted=0

# Read the ledger up front: the loop rewrites it, and entries are unique by ID
# so a namespace recorded twice is only deleted once.
while read -r id slug; do
  [[ -z "$id" ]] && continue

  status=$(curl -s -m 15 -o /dev/null -w '%{http_code}' -X DELETE \
    -H "Authorization: Bearer $token" \
    "$HUMAN_API/compose/namespace/$id" || echo 000)

  case "$status" in
    2*) echo "deleted namespace ${slug:-$id} (ID $id)" ;;
    404) echo "already gone: ${slug:-$id} (ID $id)" ;;
    *) echo "could not delete ${slug:-$id} (ID $id): HTTP $status" >&2 ;;
  esac
  deleted=$((deleted + 1))
done < <(python3 -c '
import json, sys
scope, path = sys.argv[1], sys.argv[2]
seen = set()
for line in open(path):
    line = line.strip()
    if not line:
        continue
    try:
        e = json.loads(line)
    except json.JSONDecodeError:
        continue
    if e.get("kind") != "namespace" or e["id"] in seen:
        continue
    if scope and e.get("session") != scope:
        continue
    seen.add(e["id"])
    print(e["id"], e.get("slug", ""))
' "$scope" "$LEDGER")

# Drop the entries we just handled; anything out of scope stays for its own
# session to clean up.
python3 -c '
import json, sys
scope, path = sys.argv[1], sys.argv[2]
kept = []
for line in open(path):
    line = line.strip()
    if not line:
        continue
    try:
        e = json.loads(line)
    except json.JSONDecodeError:
        continue
    if scope and e.get("session") != scope:
        kept.append(line)
open(path, "w").write("\n".join(kept) + ("\n" if kept else ""))
' "$scope" "$LEDGER"

if ((purge)); then
  ENVIRONMENT=dev server_cli compose namespaces purge 2>/dev/null | grep -v '"level":"warn"' || true
fi

echo "cleanup done (${deleted} recorded namespace(s) in scope${scope:+, session $scope})"
