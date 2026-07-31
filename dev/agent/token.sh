#!/usr/bin/env bash
# Print a valid bearer token for the local dev API (cached until near expiry).
# Prefers the oauth2 client_credentials flow via the dev-agent auth client;
# falls back to minting a JWT through the server CLI.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

TOKEN_FILE="$STATE_DIR/token"
EXP_FILE="$STATE_DIR/token-exp"

now=$(date +%s)
if [[ -f "$TOKEN_FILE" && -f "$EXP_FILE" ]] && (($(cat "$EXP_FILE") > now + 60)); then
  # trust but verify — CLI-minted tokens may live shorter than assumed
  if curl -sf -m 5 -H "Authorization: Bearer $(cat "$TOKEN_FILE")" \
    "$HUMAN_AUTH/oauth2/info" >/dev/null 2>&1; then
    cat "$TOKEN_FILE"
    exit 0
  fi
  rm -f "$TOKEN_FILE" "$EXP_FILE"
fi

save() { # $1=token $2=expires-in
  printf '%s' "$1" >"$TOKEN_FILE"
  chmod 600 "$TOKEN_FILE"
  echo $((now + $2)) >"$EXP_FILE"
  printf '%s\n' "$1"
}

mint_oauth() {
  [[ -f "$STATE_DIR/secret" ]] || return 1
  local resp tok exp
  resp=$(curl -sf -m 10 "$HUMAN_AUTH/oauth2/token" \
    -d grant_type=client_credentials \
    -d "client_id=$AGENT_CLIENT" \
    --data-urlencode "client_secret=$(cat "$STATE_DIR/secret")" \
    --data-urlencode "scope=$AGENT_SCOPE") || return 1
  tok=$(json_get access_token <<<"$resp") || return 1
  exp=$(json_get expires_in <<<"$resp") || exp=7200
  save "$tok" "$exp"
}

# note: the CLI prints the JWT to stderr (cobra default), hence 2>&1
mint_cli() {
  local tok
  tok=$(server_cli auth jwt "$AGENT_EMAIL" --scope profile --scope api 2>&1 |
    grep -Eom1 'ey[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+') || return 1
  [[ -n "$tok" ]] || return 1
  save "$tok" 7000
}

mint_oauth || mint_cli || {
  echo "failed to mint a token — run dev/agent/bootstrap.sh first" >&2
  exit 1
}
