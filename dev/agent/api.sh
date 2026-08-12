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
#  - a restarting server answers with its boot banner rather than JSON; that
#    exits 75 with a "retry shortly" message instead of printing the banner
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

# Record namespace creates so cleanup.sh knows what this session made. Only
# creates (POST to the collection, no ID in the path) qualify — an update is a
# POST to .../namespace/{id} and must not be logged as a new resource.
if [[ "$method" == "POST" && "$path" =~ ^/compose/namespace/?(\?.*)?$ ]]; then
  if ns_id=$(json_get response.namespaceID <"$body" 2>/dev/null); then
    ledger_record namespace "$ns_id" "$(json_get response.slug <"$body" 2>/dev/null || true)"
  fi
fi
if ! python3 -m json.tool <"$body" 2>/dev/null; then
  # A restarting dev server answers on the socket with its boot banner rather
  # than a response. Printing that and exiting 0 makes every caller that pipes
  # into a JSON parser fail with "Expecting value: line 1 column 1", which reads
  # as a malformed request instead of "it is coming back in a few seconds".
  if grep -qE 'Human server initializing|^(PASS|FAIL) [A-Za-z]' "$body"; then
    echo "the dev server is restarting: it answered with its boot banner, not a response." >&2
    echo "Retry in a few seconds — dev_server_status reports when it is back up." >&2
    exit 75
  fi

  cat "$body"
  echo
fi

if ((status >= 400)); then
  echo "HTTP $status" >&2
  exit 1
fi
