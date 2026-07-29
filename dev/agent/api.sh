#!/usr/bin/env bash
# Authenticated curl wrapper for the local dev API.
#
# Usage: api.sh METHOD PATH [extra curl args...]
#   api.sh GET '/system/users/?limit=5'
#   api.sh POST /compose/namespace/ -d '{"name":"Agent Test","slug":"agent-test"}'
#
# Notes:
#  - list endpoints need the trailing slash (/system/users/ — without it: 404)
#  - the API returns errors as {"error":{...}} with HTTP 200; this script
#    detects that, prints the message to stderr and exits non-zero
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

method="${1:?usage: api.sh METHOD PATH [curl args...]}"
path="${2:?usage: api.sh METHOD PATH [curl args...]}"
shift 2

token="$("$AGENT_DIR/token.sh")"

body=$(mktemp)
trap 'rm -f "$body"' EXIT

status=$(curl -s -m 30 -o "$body" -w '%{http_code}' -X "$method" \
  -H "Authorization: Bearer $token" \
  -H 'Content-Type: application/json' \
  "$HUMAN_API$path" "$@")

if errmsg=$(json_get error.message <"$body" 2>/dev/null); then
  echo "API error (HTTP $status): $errmsg" >&2
  exit 1
fi
python3 -m json.tool <"$body" 2>/dev/null || {
  cat "$body"
  echo
}

if ((status >= 400)); then
  echo "HTTP $status" >&2
  exit 1
fi
