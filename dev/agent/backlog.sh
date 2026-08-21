#!/usr/bin/env bash
#
# The queue of work deferred out of a task, so it survives the turn that found
# it.
#
# It lives in the shared .state, which every worktree symlinks — one queue for
# every session on this machine, not one per checkout. It is gitignored, so it
# never lands in a commit and never leaves this machine.
#
#   backlog.sh add TEXT [--why W] [--files F,F] [--task T]  queue an item
#   backlog.sh list [--all] [--files F]                     open items (or every item)
#   backlog.sh show ID                                      one item in full
#   backlog.sh done ID [--note N]                           close it
#   backlog.sh drop ID [--note N]                           close it as not-doing
#
# An item records where it came from — the task that deferred it, the session,
# the commit HEAD was on. Six weeks later that provenance is the difference
# between a todo you can act on and a sentence nobody can place.
#
set -euo pipefail

AGENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$AGENT_DIR/common.sh"

BACKLOG="$STATE_DIR/backlog.jsonl"
touch "$BACKLOG"

die() {
  echo "backlog: $*" >&2
  exit 1
}

py() { python3 -c "$1" "${@:2}"; }

cmd_add() {
  local text="${1:-}" why="" files="" task=""
  shift || true
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --why)
        why="$2"
        shift 2
        ;;
      --files)
        files="$2"
        shift 2
        ;;
      --task)
        task="$2"
        shift 2
        ;;
      *) die "unknown flag $1" ;;
    esac
  done
  [[ -n "$text" ]] || die "usage: backlog.sh add TEXT [--why W] [--files F,F] [--task T]"

  local head branch
  head="$(git -C "$REPO_DIR" rev-parse --short HEAD 2>/dev/null || echo '')"
  branch="$(git -C "$REPO_DIR" rev-parse --abbrev-ref HEAD 2>/dev/null || echo '')"

  # Append-only, one JSON object per line: two sessions writing at once
  # interleave lines, they do not corrupt each other's.
  py '
import json, os, secrets, sys, time
text, why, files, task, session, head, branch, path = sys.argv[1:9]
item = {
    # Random, not a truncated clock: two sessions queueing in the same
    # millisecond got the same id, and list() folds same-id records together,
    # so the second item vanished silently.
    "id": "b" + secrets.token_hex(4),
    "text": text, "why": why or None,
    "files": [f for f in files.split(",") if f] or None,
    "task": task or None, "status": "open",
    "session": session, "commit": head or None, "branch": branch or None,
    "ts": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
}
with open(path, "a") as fh:
    fh.write(json.dumps(item) + "\n")
print(item["id"])
' "$text" "$why" "$files" "$task" "$AGENT_SESSION" "$head" "$branch" "$BACKLOG"
}

# Later lines win, so a close is just another append.
cmd_list() {
  local all="" files=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --all)
        all=1
        shift
        ;;
      --files)
        files="$2"
        shift 2
        ;;
      *) die "unknown flag $1" ;;
    esac
  done
  py '
import json, sys
path, show_all, want = sys.argv[1:4]
items = {}
for line in open(path):
    line = line.strip()
    if not line:
        continue
    it = json.loads(line)
    items.setdefault(it["id"], {}).update(it)
items = list(items.values())
if not show_all:
    items = [i for i in items if i.get("status") == "open"]
if want:
    items = [i for i in items if any(want in f for f in (i.get("files") or []))]
if not items:
    print("nothing queued" if not want else "nothing queued against %s" % want)
    raise SystemExit
for i in sorted(items, key=lambda x: x["ts"]):
    mark = {"open": "[ ]", "done": "[x]", "dropped": "[-]"}.get(i["status"], "[?]")
    print("%s %s  %s" % (mark, i["id"], i["text"]))
    if i.get("files"):
        print("        files: %s" % ", ".join(i["files"]))
    if i.get("task"):
        print("        from:  %s (%s)" % (i["task"], i["ts"][:10]))
' "$BACKLOG" "$all" "$files"
}

cmd_show() {
  [[ -n "${1:-}" ]] || die "usage: backlog.sh show ID"
  py '
import json, sys
path, want = sys.argv[1:3]
item = {}
for line in open(path):
    line = line.strip()
    if line and json.loads(line)["id"] == want:
        item.update(json.loads(line))
if not item:
    sys.exit("backlog: no item %s" % want)
print(json.dumps(item, indent=2))
' "$BACKLOG" "$1"
}

cmd_close() {
  local status="$1" id="${2:-}" note=""
  shift 2 || true
  [[ "${1:-}" == "--note" ]] && note="$2"
  [[ -n "$id" ]] || die "usage: backlog.sh $status ID [--note N]"
  py '
import json, sys, time
path, want, status, note, session = sys.argv[1:6]
if not any(json.loads(l)["id"] == want for l in open(path) if l.strip()):
    sys.exit("backlog: no item %s" % want)
with open(path, "a") as fh:
    fh.write(json.dumps({
        "id": want, "status": status, "note": note or None,
        "closed_by": session,
        "closed": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
    }) + "\n")
print("%s %s" % (want, status))
' "$BACKLOG" "$id" "$status" "$note" "$AGENT_SESSION"
}

case "${1:-}" in
  add)
    shift
    cmd_add "$@"
    ;;
  list)
    shift
    cmd_list "$@"
    ;;
  show)
    shift
    cmd_show "$@"
    ;;
  done)
    shift
    cmd_close done "$@"
    ;;
  drop)
    shift
    cmd_close dropped "$@"
    ;;
  *)
    sed -n '2,19p' "${BASH_SOURCE[0]}" | sed 's/^# \?//'
    exit 1
    ;;
esac
