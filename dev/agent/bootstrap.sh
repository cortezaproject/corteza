#!/usr/bin/env bash
# Idempotent, self-healing setup of the agent dev toolkit:
#  1. verify the local dev server is up
#  2. ensure the agent@local.dev user exists (envoy import) with super-admin
#  3. ensure the dev_agent client_credentials auth client exists and is
#     correctly configured (managed via REST — envoy YAML cannot set
#     validGrant / resolve impersonateUser reliably)
#  4. cache the server-generated client secret in .state/ (gitignored),
#     verify the oauth2 flow end to end, then run the smoke check
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

section "agent toolkit bootstrap"
echo "  $(paint "$C_DIM" "$HUMAN_API")"

if ! curl -sf -m 5 "$HUMAN_BASE/version" >/dev/null; then
  bad "dev server not reachable at $HUMAN_BASE"
  note "start it: cd server && make watch"
  exit 1
fi

# `--skip-existing` keys on the YAML identifier, which IS the handle, so it
# cannot recognise a user provisioned under an earlier handle and tries to
# create a second one with the same email. Probe by email instead: the CLI
# reads the database directly and needs no token.
if ! server_cli auth jwt "$AGENT_EMAIL" >/dev/null 2>&1; then
  server_cli import --skip-existing "$AGENT_DIR/seed/dev-agent.yaml"
fi
server_cli roles useradd super-admin "$AGENT_EMAIL" >/dev/null 2>&1 ||
  note "'roles useradd super-admin' reported existing membership (ok)"
ok "user $AGENT_EMAIL ensured (super-admin)"

# Bootstrap admin token straight from the CLI (the JWT lands on stderr).
admin_tok=$(server_cli auth jwt "$AGENT_EMAIL" --scope profile --scope api 2>&1 |
  grep -Eom1 'ey[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+')
[[ -n "$admin_tok" ]] || {
  bad "could not mint bootstrap JWT via server CLI"
  exit 1
}

capi() { # capi METHOD PATH [curl args...]
  local method="$1" path="$2"
  shift 2
  curl -sf -m 15 -X "$method" -H "Authorization: Bearer $admin_tok" \
    -H 'Content-Type: application/json' "$HUMAN_API$path" "$@"
}

agent_uid=$(capi GET "/system/users/?email=$AGENT_EMAIL" | json_get response.set.0.userID) || {
  bad "cannot resolve userID of $AGENT_EMAIL"
  exit 1
}

# The seed key is the handle; rename a box provisioned under the hyphenated one.
agent_handle=$(capi GET "/system/users/$agent_uid" | json_get response.handle) || agent_handle=""
if [[ "$agent_handle" != "agent_dev" ]]; then
  capi PUT "/system/users/$agent_uid" \
    -d "{\"email\":\"$AGENT_EMAIL\",\"name\":\"Dev Agent (local tooling)\",\"handle\":\"agent_dev\"}" >/dev/null
  ok "user handle $agent_handle → agent_dev"
fi

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

# A box provisioned under the old hyphenated handle keeps its numeric ID (and
# therefore its cached secret) — the update path below renames it in place.
if [[ -z "$client_id" ]]; then
  client_id=$(capi GET "/system/auth/clients/?handle=dev-agent" |
    json_get response.set.0.authClientID) || client_id=""
  [[ -n "$client_id" ]] && echo "renaming auth client dev-agent -> $AGENT_CLIENT (ID $client_id)"
fi

if [[ -z "$client_id" ]]; then
  client_id=$(client_payload | capi POST /system/auth/clients/ -d @- |
    json_get response.authClientID)
  ok "auth client $AGENT_CLIENT created $(paint "$C_DIM" "(ID $client_id)")"
else
  client_payload | capi PUT "/system/auth/clients/$client_id" -d @- >/dev/null
  ok "auth client $AGENT_CLIENT updated $(paint "$C_DIM" "(ID $client_id)")"
fi

umask 077
capi GET "/system/auth/clients/$client_id/secret" | json_get response >"$STATE_DIR/secret"
[[ -s "$STATE_DIR/secret" ]] || {
  bad "could not read client secret"
  exit 1
}
ok "client secret cached in .state/secret"

# One identity for both paths: the auth client impersonates this user for API
# calls, and the same user logs into the webapp for browser checks. Nothing in
# the client_credentials flow cares that it also holds a password — it loads
# the impersonated user by ID (auth/handlers/handle_oauth2.go).
# Password lives in .state/ui-password (gitignored).
if [[ ! -f "$STATE_DIR/ui-password" ]] || [[ -n "${RESET_UI_PASSWORD:-}" ]]; then
  openssl rand -hex 12 >"$STATE_DIR/ui-password"
fi
capi POST "/system/users/$agent_uid/password" \
  -d "{\"password\":\"$(cat "$STATE_DIR/ui-password")\"}" >/dev/null
ok "browser login ready: $(paint "$C_CYAN" "$AGENT_EMAIL") $(paint "$C_DIM" "· password in dev/agent/.state/ui-password")"

# Converge a box provisioned while UI login was a separate user.
old_ui=$(capi GET "/system/users/?email=agent-ui@local.dev" | json_get response.set.0.userID) || old_ui=""
if [[ -n "$old_ui" ]]; then
  capi DELETE "/system/users/$old_ui" >/dev/null 2>&1 || true
  ok "removed superseded user agent-ui@local.dev $(paint "$C_DIM" "(ID $old_ui)")"
fi

# Read-only browser-login user, for checking what a user WITHOUT permission
# sees: the admin sidebar drops what it cannot reach, and an editor opened by
# deep link renders read-only. The role allows a deliberate SUBSET — users and
# roles, read only — so a gate that has stopped working shows up as entries
# that should have gone and fields that should have been disabled.
#
# It has no user group, so everything not allowed below resolves to deny.
ro_email="agent-ro@local.dev"
# Filtered client-side, and guarded below: the fixture's rules and members are
# written to whatever this resolves to, so resolving to the wrong role would
# grant them to it.
ro_role=$(capi GET "/system/roles/?limit=500" | json_pick response.set handle agent_readonly roleID) || ro_role=""
if [[ -z "$ro_role" ]]; then
  ro_role=$(capi POST /system/roles/ \
    -d '{"name":"Dev Agent read-only (RBAC fixture)","handle":"agent_readonly"}' |
    json_get response.roleID)
  ok "role agent_readonly created $(paint "$C_DIM" "(ID $ro_role)")"
fi

# Guard: never write the fixture's rules or members onto another role.
ro_role_handle=$(capi GET "/system/roles/$ro_role" | json_get response.handle) || ro_role_handle=""
[[ "$ro_role_handle" == "agent_readonly" ]] || {
  bad "refusing to configure role $ro_role — its handle is '$ro_role_handle', not agent_readonly"
  exit 1
}

# `corteza::system/` is the component resource; `corteza::system:component` is
# refused. Application access is what admits the section at all — without it the
# shell bounces every admin route to /?denied=admin.
capi PATCH "/system/permissions/$ro_role/rules" -d '{"rules":[
  {"resource":"corteza::system/","operation":"users.search","access":"allow"},
  {"resource":"corteza::system/","operation":"roles.search","access":"allow"},
  {"resource":"corteza::system:user/*","operation":"read","access":"allow"},
  {"resource":"corteza::system:role/*","operation":"read","access":"allow"},
  {"resource":"corteza::system:application/*","operation":"read","access":"allow"},
  {"resource":"corteza::system:application/*","operation":"access","access":"allow"}
]}' >/dev/null
ok "role agent_readonly rules applied"

ro_uid=$(capi GET "/system/users/?email=$ro_email" | json_get response.set.0.userID) || ro_uid=""
if [[ -z "$ro_uid" ]]; then
  ro_uid=$(capi POST /system/users/ \
    -d "{\"email\":\"$ro_email\",\"name\":\"Dev Agent RO (read-only browser login)\",\"handle\":\"agent_ro\"}" |
    json_get response.userID)
  ok "user $ro_email created"
fi
capi POST "/system/roles/$ro_role/member/$ro_uid" >/dev/null 2>&1 || true
if [[ ! -f "$STATE_DIR/ro-password" ]] || [[ -n "${RESET_UI_PASSWORD:-}" ]]; then
  openssl rand -hex 12 >"$STATE_DIR/ro-password"
fi
capi POST "/system/users/$ro_uid/password" \
  -d "{\"password\":\"$(cat "$STATE_DIR/ro-password")\"}" >/dev/null
ok "read-only login ready: $(paint "$C_CYAN" "$ro_email") $(paint "$C_DIM" "· password in dev/agent/.state/ro-password")"

rm -f "$STATE_DIR/token" "$STATE_DIR/token-exp"
if ! "$AGENT_DIR/token.sh" >/dev/null; then
  bad "oauth client_credentials flow failed for client '$AGENT_CLIENT'"
  exit 1
fi
ok "oauth client_credentials flow verified"

exec "$AGENT_DIR/smoke.sh"
