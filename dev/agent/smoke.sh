#!/usr/bin/env bash
# Answer "what state is the world in" in one call:
# server up? which version? token mintable? who am I? API reachable with auth?
set -uo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

# ok/bad/note come from tty.sh; only the exit code is this script's own.
fail=0
failed() {
  bad "$@"
  fail=1
}

echo "$(paint "$C_BOLD" "$G_SECTION smoke") $(paint "$C_DIM" "$HUMAN_API")"

ver=$(curl -sf -m 5 "$HUMAN_BASE/version" | json_get response.version 2>/dev/null)
if [[ -n "${ver:-}" ]]; then
  ok "server up (version: $ver)"
else
  # The watcher rebuilds on any .go write and the API is down for a moment while
  # it swaps the binary; "start it" is the wrong instruction then, and waiting is
  # the right one. Matched by this checkout's own binary path, because another
  # checkout on this machine runs a watcher of the same name.
  if pgrep -f "devwatch .*-bin $SERVER_DIR/build/dev-bin" >/dev/null 2>&1; then
    failed "server is building at $HUMAN_BASE — the watcher is up but the binary is not serving yet."
    note "Retry in ~15s; dev/agent/logs.sh -n 20 shows the build."
  else
    failed "server not reachable at $HUMAN_BASE — start it: cd server && make watch"
  fi
  exit 1
fi

tok=$("$AGENT_DIR/token.sh" 2>/dev/null)
if [[ -n "${tok:-}" ]]; then
  ok "token minted"
else
  failed "cannot mint token — run dev/agent/bootstrap.sh"
  exit 1
fi

who=$(curl -sf -m 5 -H "Authorization: Bearer $tok" "$HUMAN_AUTH/oauth2/info" 2>/dev/null)
if [[ -n "${who:-}" ]]; then
  sub=$(json_get sub <<<"$who" 2>/dev/null || echo '?')
  usr=$(curl -sf -m 5 -H "Authorization: Bearer $tok" "$HUMAN_API/system/users/$sub" 2>/dev/null)
  ok "identity: $(json_get response.name <<<"$usr" 2>/dev/null || echo '?') <$(json_get response.email <<<"$usr" 2>/dev/null || echo '?')> (userID $sub)"
else
  failed "token rejected by $HUMAN_AUTH/oauth2/info"
fi

if "$AGENT_DIR/api.sh" GET '/system/users/?limit=1' >/dev/null 2>&1; then
  ok "authenticated API access works"
else
  failed "authenticated API call failed (GET /system/users/)"
fi

exit $fail
