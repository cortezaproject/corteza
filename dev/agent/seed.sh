#!/usr/bin/env bash
# Import dev fixtures into the local dev server.
#
# Usage: seed.sh [--force] [fixture...]
#   fixture = directory name under dev/fixtures/ (default: all)
#   --force = delete the fixture's namespace first, then re-import
#
# Convention: each fixture lives in dev/fixtures/<slug>/def.yaml where <slug>
# is also the compose namespace slug (agent- prefixed). Because record CSV
# imports are not idempotent, an existing namespace is skipped unless --force.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

FIXTURES_DIR="$REPO_DIR/dev/fixtures"

force=0
picks=()
for a in "$@"; do
  case "$a" in
  --force) force=1 ;;
  *) picks+=("$a") ;;
  esac
done
[[ ${#picks[@]} -gt 0 ]] || {
  for d in "$FIXTURES_DIR"/*/; do picks+=("$(basename "$d")"); done
}

token="$("$AGENT_DIR/token.sh")"

ns_id() { # ns_id SLUG — print namespaceID or fail
  curl -sf -m 10 -H "Authorization: Bearer $token" \
    "$HUMAN_API/compose/namespace/?slug=$1&limit=1" |
    json_get response.set.0.namespaceID
}

for slug in "${picks[@]}"; do
  def="$FIXTURES_DIR/$slug/def.yaml"
  [[ -f "$def" ]] || {
    echo "no such fixture: $slug ($def missing)" >&2
    exit 1
  }

  if id=$(ns_id "$slug" 2>/dev/null); then
    if ((force)); then
      curl -sf -m 15 -X DELETE -H "Authorization: Bearer $token" \
        "$HUMAN_API/compose/namespace/$id" >/dev/null
      echo "deleted existing namespace $slug (ID $id)"
    else
      echo "fixture $slug already seeded (namespace ID $id) — use --force to re-import"
      continue
    fi
  fi

  # import a staged copy of the fixture DIRECTORY — sibling CSVs are only
  # registered as record datasource providers on directory decode, and
  # ui.json must be excluded (the decoder parses .json as YAML and errors)
  stage="$STATE_DIR/import-$slug"
  rm -rf "$stage"
  mkdir -p "$stage"
  cp "$FIXTURES_DIR/$slug"/* "$stage"/
  rm -f "$stage/ui.json"
  server_cli import --skip-existing "$stage"
  id=$(ns_id "$slug") || {
    echo "import ran but namespace $slug not found" >&2
    exit 1
  }

  # CLI import writes straight to the store, bypassing the running server's
  # DAL registry — record queries would fail with "model does not exist".
  # A no-op module update via REST (POST, not PUT!) triggers DalModelReplace
  # and registers each model on the live server.
  mids=$(curl -sf -m 15 -H "Authorization: Bearer $token" \
    "$HUMAN_API/compose/namespace/$id/module/?limit=200" |
    python3 -c '
import json, sys
for m in json.load(sys.stdin)["response"]["set"]:
    print(m["moduleID"])
')
  for mid in $mids; do
    curl -sf -m 15 -H "Authorization: Bearer $token" \
      "$HUMAN_API/compose/namespace/$id/module/$mid" |
      json_get response |
      curl -sf -m 15 -X POST -H "Authorization: Bearer $token" \
        -H 'Content-Type: application/json' -d @- \
        "$HUMAN_API/compose/namespace/$id/module/$mid" >/dev/null || {
      echo "DAL re-registration failed for module $mid" >&2
      exit 1
    }
  done

  # presentation layer (charts + pages) is built via REST — envoy YAML
  # cannot resolve page-block refs, see pagebuild.py
  if [[ -f "$FIXTURES_DIR/$slug/ui.json" ]]; then
    python3 "$AGENT_DIR/pagebuild.py" "$slug" "$FIXTURES_DIR/$slug/ui.json"
  fi

  # Refresh the handle → ID map so checks can look resources up instead of
  # spending a request per kind rediscovering what seeding just created.
  "$AGENT_DIR/ids.sh" "$slug" >/dev/null || true

  echo "seeded fixture $slug (namespace ID $id)"
done
