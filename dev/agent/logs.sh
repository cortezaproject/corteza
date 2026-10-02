#!/usr/bin/env bash
# Read the dev server log (server/build/dev.log, written by `make watch`).
#
# Usage: logs.sh [-n LINES] [-w] [-a] [-f] [PATTERN]
#   logs.sh                 # errors and build markers in the last 400 lines
#   logs.sh -w              # … warnings too
#   logs.sh -a              # the raw tail, last 100 lines
#   logs.sh -n 500 error    # last 500 raw lines matching "error" (case-insensitive)
#   logs.sh -f              # follow
#
# The default view is for reading after a failed call: one line per error-level
# entry (time, level, logger, msg, error), the [devwatch] build/serve markers
# that place it, and a footer counting what was hidden. The raw tail is behind
# -a because most of it is a restart's warn banner, repeated per rebuild.
#
# Identical entries (level, logger, msg, error) print once with a ×N count, and
# stamps carry the date when the window spans more than one day.
#
# HUMAN_DEV_LOG points the script at another checkout's log, e.g. the primary's
# server/build/dev.log from a worktree whose own server is not running.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

LOG="${HUMAN_DEV_LOG:-$SERVER_DIR/build/dev.log}"

lines=""
follow=0
raw=0
warn=0
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
  -a)
    raw=1
    shift
    ;;
  -w)
    warn=1
    shift
    ;;
  *)
    pattern="$1"
    shift
    ;;
  esac
done

if [[ ! -f "$LOG" ]]; then
  bad "no dev log at $LOG"
  note "this checkout's server has not run under 'make watch' (it tees its output there)"
  note "another checkout's log: HUMAN_DEV_LOG=<checkout>/server/build/dev.log $0"
  exit 1
fi

if ((follow)); then
  exec tail -f "$LOG"
elif [[ -n "$pattern" ]]; then
  tail -n "${lines:-100}" "$LOG" | grep -i --color="$GREP_COLOUR_WHEN" -- "$pattern" || {
    warn "no match for '$pattern' in last ${lines:-100} lines"
    exit 1
  }
elif ((raw)); then
  tail -n "${lines:-100}" "$LOG"
else
  tail -n "${lines:-400}" "$LOG" | python3 -c '
import json, sys, time

show_warn = sys.argv[1] == "1"
window = sys.argv[2]
loud = {"error", "dpanic", "panic", "fatal"}
hidden = {"warn": 0, "info": 0, "debug": 0}
shown = 0
out = []            # (ts, text) in first-seen order; repeats fold into their first line
seen = {}           # entry key -> index into out
counts_by = {}

def ts_of(e):
    try:
        return float(e.get("ts"))
    except (TypeError, ValueError):
        return None

for raw in sys.stdin:
    line = raw.rstrip("\n")
    if not line:
        continue
    if line.startswith("[devwatch]"):
        out.append((None, line))
        continue
    try:
        e = json.loads(line)
    except ValueError:
        # a panic trace or a bare "Error:" has no envelope and is always news
        out.append((None, line))
        shown += 1
        continue
    level = str(e.get("level", "")).lower()
    if level in loud or (show_warn and level == "warn"):
        parts = [level.upper()]
        if e.get("logger"):
            parts.append(e["logger"] + ":")
        parts.append(str(e.get("msg", "")))
        if e.get("error"):
            parts.append("— " + str(e["error"]))
        key = " ".join(parts)
        if key in seen:
            counts_by[key] += 1
        else:
            seen[key] = len(out)
            counts_by[key] = 1
            out.append((ts_of(e), key))
            shown += 1
    elif level in hidden:
        hidden[level] += 1
    else:
        hidden.setdefault(level or "other", 0)
        hidden[level or "other"] += 1

stamps = [t for t, _ in out if t is not None]
fmt = "%m-%d %H:%M:%S" if stamps and max(stamps) - min(stamps) > 86400 else "%H:%M:%S"
for t, text in out:
    if t is None and text.startswith("[devwatch]"):
        print(text)
        continue
    n = counts_by.get(text, 1)
    prefix = time.strftime(fmt, time.localtime(t)) + " " if t is not None else ""
    print(prefix + (f"×{n} " if n > 1 else "") + text)

counts = ", ".join(f"{n} {k}" for k, n in hidden.items() if n)
if shown == 0:
    print(f"no errors in the last {window} lines" + (f" ({counts} hidden)" if counts else ""))
elif counts:
    print(f"— hidden: {counts}; -w shows warnings, -a the raw tail")
' "$warn" "${lines:-400}"
fi
