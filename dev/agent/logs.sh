#!/usr/bin/env bash
# Read the dev server log (server/build/dev.log, written by `make watch`).
#
# Usage: logs.sh [-n LINES] [-f] [PATTERN]
#   logs.sh                 # last 100 lines
#   logs.sh -n 500 error    # last 500 lines matching "error" (case-insensitive)
#   logs.sh -f              # follow
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

LOG="$SERVER_DIR/build/dev.log"

lines=100
follow=0
pattern=""
while [[ $# -gt 0 ]]; do
  case "$1" in
  -n)
    lines="$2"
    shift 2
    ;;
  -f)
    follow=1
    shift
    ;;
  *)
    pattern="$1"
    shift
    ;;
  esac
done

if [[ ! -f "$LOG" ]]; then
  echo "no dev log at $LOG — (re)start the dev server with 'cd server && make watch'" >&2
  echo "(the watch target tees its output there; an older session may predate this)" >&2
  exit 1
fi

if ((follow)); then
  exec tail -f "$LOG"
elif [[ -n "$pattern" ]]; then
  tail -n "$lines" "$LOG" | grep -i --color=never -- "$pattern" || {
    echo "(no match for '$pattern' in last $lines lines)" >&2
    exit 1
  }
else
  tail -n "$lines" "$LOG"
fi
