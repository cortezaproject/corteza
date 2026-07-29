#!/usr/bin/env bash
# Idempotent, self-healing setup of the agent dev toolkit:
#  1. verify the local dev server is up
#  2. ensure the agent-dev user exists (envoy import) with the super-admin role
#  3. ensure the dev-agent client_credentials auth client exists and is
#     correctly configured (managed via REST — envoy YAML cannot set
#     validGrant / resolve impersonateUser reliably)
#  4. cache the server-generated client secret in .state/ (gitignored),
#     verify the oauth2 flow end to end, then run the smoke check
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

echo "== agent toolkit bootstrap ($HUMAN_API)"

if ! curl -sf -m 5 "$HUMAN_BASE/version" >/dev/null; then
  echo "dev server not reachable at $HUMAN_BASE — start it: cd server && make watch" >&2
  exit 1
fi

server_cli import --skip-existing "$AGENT_DIR/seed/dev-agent.yaml"
server_cli roles useradd super-admin "$AGENT_EMAIL" >/dev/null 2>&1 ||
  echo "note: 'roles useradd super-admin' reported existing membership (ok)"
echo "user $AGENT_EMAIL ensured (super-admin)"

# Bootstrap admin token straight from the CLI (the JWT lands on stderr).
admin_tok=$(server_cli auth jwt "$AGENT_EMAIL" --scope profile --scope api 2>&1 |
  grep -Eom1 'ey[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+')
[[ -n "$admin_tok" ]] || {
  echo "could not mint bootstrap JWT via server CLI" >&2
  exit 1
}

capi() { # capi METHOD PATH [curl args...]
  local method="$1" path="$2"
  shift 2
  curl -sf -m 15 -X "$method" -H "Authorization: Bearer $admin_tok" \
    -H 'Content-Type: application/json' "$HUMAN_API$path" "$@"
}

agent_uid=$(capi GET "/system/users/?email=$AGENT_EMAIL" | json_get response.set.0.userID) || {
  echo "cannot resolve userID of $AGENT_EMAIL" >&2
  exit 1
}

client_payload() {
  cat <<EOF
{
  "handle": "$AGENT_CLIENT",
  "meta": {
    "name": "Dev Agent (local tooling)",
    "description": "Local-only client_credentials client used by agent dev tooling (dev/agent/*). Do not create on shared or production instances."
  },
  "validGrant": "client_credentials",
  "scope": "$AGENT_SCOPE",
  "enabled": true,
  "trusted": true,
  "security": { "impersonateUser": "$agent_uid" }
}
EOF
}

client_id=$(capi GET "/system/auth/clients/?handle=$AGENT_CLIENT" |
  json_get response.set.0.authClientID) || client_id=""

if [[ -z "$client_id" ]]; then
  client_id=$(client_payload | capi POST /system/auth/clients/ -d @- |
    json_get response.authClientID)
  echo "auth client $AGENT_CLIENT created (ID $client_id)"
else
  client_payload | capi PUT "/system/auth/clients/$client_id" -d @- >/dev/null
  echo "auth client $AGENT_CLIENT updated (ID $client_id)"
fi

umask 077
capi GET "/system/auth/clients/$client_id/secret" | json_get response >"$STATE_DIR/secret"
[[ -s "$STATE_DIR/secret" ]] || {
  echo "could not read client secret" >&2
  exit 1
}
echo "client secret cached in .state/secret"

rm -f "$STATE_DIR/token" "$STATE_DIR/token-exp"
if ! "$AGENT_DIR/token.sh" >/dev/null; then
  echo "oauth client_credentials flow failed for client '$AGENT_CLIENT'" >&2
  exit 1
fi
echo "oauth client_credentials flow verified"

exec "$AGENT_DIR/smoke.sh"
