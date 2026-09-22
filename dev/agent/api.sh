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
    local line first=1
    for line in "$@"; do
      if [[ -n "$first" ]]; then
        bad "$line"
        first=""
      else
        note "$line"
      fi
    done
  fi
  exit "$code"
}

method="${1:?usage: api.sh [--json] METHOD PATH [curl args...]}"
path="${2:?usage: api.sh [--json] METHOD PATH [curl args...]}"
shift 2

# A JSON body handed over as a bare argument is the easy mistake here, because
# every other tool in this toolkit takes one that way. curl reads a non-flag
# argument as a URL, so the request goes out with NO body and the server
# accepts it: a create lands as an empty record and reports success. Refuse
# instead — the body has to reach curl behind -d.
prev=""
for arg in "$@"; do
  case "$arg" in
    '{'* | '['*)
      case "$prev" in
        -d | --data | --data-raw | --data-binary | --data-ascii | --json) ;;
        *)
          die_diag usage 2 \
            "JSON body passed as a bare argument; curl reads it as a URL and sends no body." \
            "Put it behind -d:  api.sh $method $path -d '$arg'"
          ;;
      esac
      ;;
  esac
  prev="$arg"
done

token="$("$AGENT_DIR/token.sh")"

body=$(mktemp)
trap 'rm -f "$body"' EXIT

status=$(curl -s -m 30 -o "$body" -w '%{http_code}' -X "$method" \
  -H "Authorization: Bearer $token" \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  "$HUMAN_API$path" "$@")

if errmsg=$(json_get error.message <"$body" 2>/dev/null); then
  # the API reports its own errors inside a 200; hand the body back verbatim
  # in json mode so the caller reads the real error object, not a rendering
  if [[ -n "$json_mode" ]]; then
    cat "$body"
    echo
    exit 1
  fi
  bad "API error (HTTP $status): $errmsg"
  exit 1
fi

# An automation write reports its own refusal in `issues` beside a 200 and a
# valid body, so the response parses, the exit code is 0 and a TAQ the runtime
# has rejected is indistinguishable from a working one. Treat a severity=error
# issue as what it is: a failed write.
#
# The row is stored all the same — issues are a lint report on a saved draft,
# not a rejected write — so it holds its handle, and a retried create answers
# "handle not unique" instead of naming the real fault. The ID is reported
# below so the retry is an update.
if issue=$(python3 - "$body" <<'PY' 2>/dev/null
import json, sys

try:
    d = json.load(open(sys.argv[1]))
except Exception:
    sys.exit(1)
r = d.get("response")
issues = (r or {}).get("issues") if isinstance(r, dict) else None
for i in issues or []:
    if isinstance(i, dict) and i.get("severity") == "error":
        msg = f"{i.get('code', 'issue')}: {i.get('message', '')}".strip()
        rid = r.get("automationID") or r.get("workflowID") or r.get("triggerID")
        print(f"{msg}\t{rid or ''}")
        sys.exit(0)
sys.exit(1)
PY
); then
  if [[ -n "$json_mode" ]]; then
    cat "$body"
    echo
  else
    stored_id="${issue##*$'\t'}"
    issue="${issue%%$'\t'*}"
    bad "API refused this automation (HTTP $status): $issue"
    note "The body is valid and the status is 200 — the refusal is in response.issues."
    if [[ -n "$stored_id" ]]; then
      note "It was stored anyway as ID $stored_id and holds its handle: fix the spec and PUT/POST to that ID. Repeating the create answers \"handle not unique\"."
    fi
  fi
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

  # Some endpoints — the automation ones especially — report failure as
  # text/plain beside HTTP 200 instead of the {"error":{...}} envelope. Without
  # this the body is printed and the exit code stays 0, so a refused write
  # reads as a success to anything checking $?.
  if head -c 200 "$body" | grep -qE '^Error: '; then
    die_diag plain_error 1 "$(head -1 "$body")" \
      "(reported as text/plain beside HTTP $status, not the usual error envelope)"
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
  bad "HTTP $status"
  exit 1
fi
