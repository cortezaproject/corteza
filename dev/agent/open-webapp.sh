#!/usr/bin/env bash
#
# Open this checkout's webapp in a browser.
#
# The URL comes from stack.sh rather than from a scan of listening ports: with
# several worktrees up, a scan cannot tell this checkout's vite from another
# lane's, and it picks whichever answers first.
#
# The opener is chosen by platform. WSL is checked before xdg-open because a
# WSL box usually has xdg-open and no Linux browser for it to reach, so trying
# it first succeeds silently and opens nothing.
set -euo pipefail

# shellcheck source=dev/agent/stack.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/stack.sh"

echo "webapp   $HUMAN_WEBAPP"

if [[ -n "${WSL_DISTRO_NAME:-}" ]] || grep -qi microsoft /proc/version 2>/dev/null; then
  openers=(wslview cmd.exe xdg-open)
elif [[ "$(uname -s)" == Darwin ]]; then
  openers=(open)
else
  openers=(xdg-open)
fi

for opener in "${openers[@]}"; do
  command -v "$opener" >/dev/null 2>&1 || continue
  if [[ "$opener" == cmd.exe ]]; then
    "$opener" /c start "$HUMAN_WEBAPP" </dev/null >/dev/null 2>&1
  else
    "$opener" "$HUMAN_WEBAPP" >/dev/null 2>&1
  fi
  exit 0
done

echo "no browser opener found (${openers[*]}) — open the URL above yourself" >&2
