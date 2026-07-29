#!/usr/bin/env bash
# Answer "what state is the world in" in one call:
# server up? which version? token mintable? who am I? API reachable with auth?
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

fail=0
ok() { echo "  ✓ $*"; }
bad() {
  echo "  ✗ $*" >&2
  fail=1
}

echo "smoke: $HUMAN_API"

ver=$(curl -sf -m 5 "$HUMAN_BASE/version" | json_get response.version 2>/dev/null)
if [[ -n "${ver:-}" ]]; then
  ok "server up (version: $ver)"
else
  bad "server not reachable at $HUMAN_BASE — start it: cd server && make watch"
  exit 1
fi

tok=$("$AGENT_DIR/token.sh" 2>/dev/null)
if [[ -n "${tok:-}" ]]; then
  ok "token minted"
else
  bad "cannot mint token — run dev/agent/bootstrap.sh"
  exit 1
fi

who=$(curl -sf -m 5 -H "Authorization: Bearer $tok" "$HUMAN_AUTH/oauth2/info" 2>/dev/null)
if [[ -n "${who:-}" ]]; then
  sub=$(json_get sub <<<"$who" 2>/dev/null || echo '?')
  usr=$(curl -sf -m 5 -H "Authorization: Bearer $tok" "$HUMAN_API/system/users/$sub" 2>/dev/null)
  ok "identity: $(json_get response.name <<<"$usr" 2>/dev/null || echo '?') <$(json_get response.email <<<"$usr" 2>/dev/null || echo '?')> (userID $sub)"
else
  bad "token rejected by $HUMAN_AUTH/oauth2/info"
fi

if "$AGENT_DIR/api.sh" GET '/system/users/?limit=1' >/dev/null 2>&1; then
  ok "authenticated API access works"
else
  bad "authenticated API call failed (GET /system/users/)"
fi

exit $fail
