#!/usr/bin/env bash
# Delete all agent-created disposable data: every compose namespace whose
# slug starts with "agent-". Never touches unprefixed data. Re-seed fixtures
# afterwards with seed.sh.
#
# --purge additionally HARD-deletes all soft-deleted agent-* namespaces
# (repeated test cycles leave slug-sharing corpses that pollute lookups).
# Runs the dev-only server CLI purge with ENVIRONMENT=dev — an unset
# ENVIRONMENT counts as production and the command refuses; localhost has
# already been enforced by common.sh.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

purge=0
[[ "${1:-}" == "--purge" ]] && purge=1

token="$("$AGENT_DIR/token.sh")"

curl -sf -m 15 -H "Authorization: Bearer $token" \
  "$HUMAN_API/compose/namespace/?query=agent-&limit=200" |
  python3 -c '
import json, sys
for ns in json.load(sys.stdin)["response"]["set"]:
    if ns.get("slug", "").startswith("agent-"):
        print(ns["namespaceID"], ns["slug"])
' |
  while read -r id slug; do
    curl -sf -m 15 -X DELETE -H "Authorization: Bearer $token" \
      "$HUMAN_API/compose/namespace/$id" >/dev/null
    echo "deleted namespace $slug (ID $id)"
  done

if ((purge)); then
  ENVIRONMENT=dev server_cli compose namespaces purge 2>/dev/null | grep -v '"level":"warn"' || true
fi

echo "cleanup done"
