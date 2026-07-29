#!/usr/bin/env bash
# Delete all agent-created disposable data: every compose namespace whose
# slug starts with "agent-". Never touches unprefixed data. Re-seed fixtures
# afterwards with seed.sh.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

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

echo "cleanup done"
