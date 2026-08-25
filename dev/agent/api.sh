#!/usr/bin/env bash
# Authenticated curl wrapper for the local dev API.
#
# Usage: api.sh [--json] METHOD PATH [extra curl args...]
#   api.sh GET '/system/users/?limit=5'
#   api.sh POST /compose/namespace/ -d '{"name":"Agent Test","slug":"agent-test"}'
#   api.sh --json GET /nope/          # every outcome is JSON on stdout
#
# Notes:
#  - --json (or HUMAN_API_JSON=1) guarantees a JSON object on stdout whatever
#    happens, so a caller piping into a parser never dies on prose. The
#    diagnostics below are good to read and fatal to parse: without it a
#    failure prints to stderr and leaves stdout empty, which reaches python as
#    "Expecting value: line 1 column 1" — a parse error standing in for a
#    perfectly clear message. Exit codes are unchanged either way.
#  - list endpoints need the trailing slash (/system/users/ — without it: 404)
#  - the API returns errors as {"error":{...}} with HTTP 200; this script
#    detects that, prints the message to stderr and exits non-zero
#  - a restarting server answers with its boot banner rather than JSON; that
#    exits 75 with a "retry shortly" message instead of printing the banner
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

json_mode="${HUMAN_API_JSON:-}"
if [[ "${1:-}" == "--json" ]]; then
  json_mode=1
  shift
fi

# One diagnostic, in whichever shape the caller can consume.
die_diag() {
  local kind="$1" code="$2"
  shift 2
  if [[ -n "$json_mode" ]]; then
    python3 -c 'import json,sys; print(json.dumps({"error":{"kind":sys.argv[1],"message":" ".join(sys.argv[2:])}}))' \
      "$kind" "$@"
  else
    local line
    for line in "$@"; do echo "$line" >&2; done
  fi
  exit "$code"
}

method="${1:?usage: api.sh [--json] METHOD PATH [curl args...]}"
path="${2:?usage: api.sh [--json] METHOD PATH [curl args...]}"
shift 2

token="$("$AGENT_DIR/token.sh")"

body=$(mktemp)
trap 'rm -f "$body"' EXIT

status=$(curl -s -m 30 -o "$body" -w '%{http_code}' -X "$method" \
  -H "Authorization: Bearer $token" \
  -H 'Content-Type: application/json' \
  "$HUMAN_API$path" "$@")

if errmsg=$(json_get error.message <"$body" 2>/dev/null); then
  # the API reports its own errors inside a 200; hand the body back verbatim
  # in json mode so the caller reads the real error object, not a rendering
  if [[ -n "$json_mode" ]]; then
    cat "$body"
    echo
    exit 1
  fi
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
    die_diag restarting 75 \
      "the dev server is restarting: it answered with its boot banner, not a response." \
      "Retry in a few seconds — dev_server_status reports when it is back up."
  fi

  # No route matched, so the server fell through to the SPA and served the
  # webapp's index.html with HTTP 200. Dumping a page of markup reads as a
  # broken endpoint; the cause is almost always a missing service prefix,
  # because the JS clients carry theirs in the client's base URL and the paths
  # in api-clients/*.ts are written relative to it.
  if grep -qiE '<!doctype html|<html' "$body"; then
    die_diag no_route 44 \
      "no API route matched $path — the server returned the webapp's HTML, not a response." \
      "Paths here are absolute and need the service prefix: /system, /compose," \
      "/automation, /federation. So '/agents/' is '/system/agents/'." \
      "Collection endpoints also need the trailing slash."
  fi

  if [[ -n "$json_mode" ]]; then
    die_diag not_json 1 "server returned a non-JSON body (HTTP $status): $(head -c 200 "$body")"
  fi
  cat "$body"
  echo
fi

if ((status >= 400)); then
  if [[ -n "$json_mode" ]]; then
    exit 1
  fi
  echo "HTTP $status" >&2
  exit 1
fi
