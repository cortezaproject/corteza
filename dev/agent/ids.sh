#!/usr/bin/env bash
# Print a handle → ID map for a namespace, and cache it for scripts to read.
#
# Usage: ids.sh [SLUG...]          (default: every fixture under dev/fixtures)
#   ids.sh catalogue               refresh one namespace
#   ids.sh --path catalogue        print where the cache lives and exit
#
# Why: a browser or API check needs real IDs, and hunting them costs one
# request per resource kind plus a guess at each endpoint's prefix. This does
# it once and writes .state/ids.json:
#
#   { "catalogue": { "namespaceID": "…",
#                    "module": { "catalogue_field": "…" },
#                    "page":   { "field_catalogue_record": "…" },
#                    "chart":  { "cat_chart_bar": "…" },
#                    "record": { "catalogue_field": "…" } } }
#
# Pages and charts are keyed by handle, falling back to a slugified title for
# the many pages that have no handle. `record` holds one arbitrary record ID
# per module — enough to open a record screen, which is all a UI check needs.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

IDS_FILE="$STATE_DIR/ids.json"

if [[ "${1:-}" == "--path" ]]; then
  echo "$IDS_FILE"
  exit 0
fi

picks=("$@")
if [[ ${#picks[@]} -eq 0 ]]; then
  for d in "$REPO_DIR/dev/fixtures"/*/; do picks+=("$(basename "$d")"); done
fi

token="$("$AGENT_DIR/token.sh")"

fetch() { # fetch PATH — authenticated GET, empty on failure
  curl -sf -m 20 -H "Authorization: Bearer $token" "$HUMAN_API$1" || true
}

for slug in "${picks[@]}"; do
  ns=$(fetch "/compose/namespace/?slug=$slug&limit=1")
  ns_id=$(printf '%s' "$ns" | json_get response.set.0.namespaceID 2>/dev/null || true)
  if [[ -z "$ns_id" ]]; then
    echo "no namespace with slug $slug" >&2
    continue
  fi

  modules=$(fetch "/compose/namespace/$ns_id/module/?limit=200")
  pages=$(fetch "/compose/namespace/$ns_id/page/?limit=200")
  charts=$(fetch "/compose/namespace/$ns_id/chart/?limit=200")

  # One record per module, fetched only for modules that have any.
  records="{}"
  for line in $(printf '%s' "$modules" | python3 -c '
import json, sys
try:
    for m in json.load(sys.stdin)["response"]["set"]:
        print(m["moduleID"] + ":" + (m.get("handle") or m["moduleID"]))
except Exception:
    pass
'); do
    mid="${line%%:*}"
    handle="${line#*:}"
    rid=$(fetch "/compose/namespace/$ns_id/module/$mid/record/?limit=1" |
      json_get response.set.0.recordID 2>/dev/null || true)
    [[ -n "$rid" ]] && records=$(printf '%s' "$records" |
      python3 -c 'import json,sys;d=json.load(sys.stdin);d[sys.argv[1]]=sys.argv[2];print(json.dumps(d))' \
        "$handle" "$rid")
  done

  MODULES="$modules" PAGES="$pages" CHARTS="$charts" RECORDS="$records" \
    python3 -c '
import json, os, re, sys

slug, ns_id, path = sys.argv[1], sys.argv[2], sys.argv[3]

def load(env):
    try:
        return json.loads(os.environ[env])["response"]["set"] or []
    except Exception:
        return []

def key(item, id_field):
    # Handle first; a page usually has none, so its title is slugified into one.
    h = item.get("handle") or ""
    if not h:
        h = re.sub(r"[^a-z0-9]+", "_", (item.get("title") or item.get("name") or "").lower()).strip("_")
    return h or item[id_field]

out = {"namespaceID": ns_id}
for env, id_field, name in (
    ("MODULES", "moduleID", "module"),
    ("PAGES", "pageID", "page"),
    ("CHARTS", "chartID", "chart"),
):
    out[name] = {key(i, id_field): i[id_field] for i in load(env)}
out["record"] = json.loads(os.environ["RECORDS"])

try:
    with open(path) as fh:
        all_ids = json.load(fh)
except Exception:
    all_ids = {}

all_ids[slug] = out
with open(path, "w") as fh:
    json.dump(all_ids, fh, indent=2, sort_keys=True)

counts = ", ".join(f"{len(out[k])} {k}" for k in ("module", "page", "chart", "record"))
print(f"{slug}: namespace {ns_id} ({counts}) -> {path}")
' "$slug" "$ns_id" "$IDS_FILE"
done
