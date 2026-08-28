#!/usr/bin/env bash
#
# The queue of work deferred out of a task, so it survives the turn that found
# it. Two tiers:
#
#   local   this checkout's todo — the worktree's, or the session's on the
#           primary. Read back automatically, at triage and at the end of a
#           task.
#   global  everything else, owned by nobody. Read only when somebody asks
#           for it.
#
#   backlog.sh add TEXT [--why W] [--files F,F] [--task T] [--global]
#   backlog.sh promote ID                                   local -> global
#   backlog.sh list [--global|--both|--orphaned] [--closed] [--files F]
#   backlog.sh show ID                                      one item in full
#   backlog.sh done ID [--note N]                           close it
#   backlog.sh drop ID [--note N]                           close it as not-doing
#   backlog.sh release OWNER                                a whole todo -> global
#   backlog.sh owner                                        who owns local items here
#
# An item is local when doing it would change a file the current task is
# already changing, or when it follows directly from that change. Everything
# else is global: filed, named in the report, and out of the way until the
# pool is asked for.
#
# A worktree owns its todo, so the queue survives the session that opened it
# and dies with `worktree.sh rm`, which promotes what is still open. On the
# primary there is no such boundary, so the session is the owner.
#
# Every scope this file did not write is global — the shared pool, and the
# per-session items from before the tiers existed. Sessions that have since
# ended cannot act on them, and the pool is where unowned work belongs.
#
# `--orphaned` is the recovery path: local items whose worktree is gone or
# whose session has ended without promoting them.
#
# An item records where it came from — the task that deferred it, the session,
# the commit HEAD was on. Six weeks later that provenance is the difference
# between a todo you can act on and a sentence nobody can place.
#
set -euo pipefail

AGENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$AGENT_DIR/common.sh"

BACKLOG="$STATE_DIR/backlog.jsonl"
WT_DIR="$STATE_DIR/worktrees"
touch "$BACKLOG"

# Who owns a local item here: the worktree this checkout is, else the session.
# A registry entry's path is the checkout it was made for, so a worktree
# answers with its own name however many sessions have driven it.
owner_now() {
  local f path
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    path="$(json_get path <"$f" 2>/dev/null || true)"
    if [[ "$path" == "$REPO_DIR" ]]; then
      echo "wt:$(json_get name <"$f")"
      return
    fi
  done
  echo "s:$AGENT_SESSION"
}

die() {
  bad "backlog: $*"
  exit 1
}

py() { python3 -c "$1" "${@:2}"; }

cmd_add() {
  local text="${1:-}" why="" files="" task="" scope="local"
  shift || true
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --global)
        scope="global"
        shift
        ;;
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
  [[ -n "$text" ]] || die "usage: backlog.sh add TEXT [--why W] [--files F,F] [--task T] [--global]"

  local head branch
  head="$(git -C "$REPO_DIR" rev-parse --short HEAD 2>/dev/null || echo '')"
  branch="$(git -C "$REPO_DIR" rev-parse --abbrev-ref HEAD 2>/dev/null || echo '')"

  # Append-only, one JSON object per line: two sessions writing at once
  # interleave lines, they do not corrupt each other's.
  py '
import json, os, secrets, sys, time
text, why, files, task, session, head, branch, path, scope, owner = sys.argv[1:11]
item = {
    # Random, not a truncated clock: two sessions queueing in the same
    # millisecond got the same id, and list() folds same-id records together,
    # so the second item vanished silently.
    "id": "b" + secrets.token_hex(4),
    "text": text, "why": why or None,
    "files": [f for f in files.split(",") if f] or None,
    "task": task or None, "status": "open", "scope": scope,
    "owner": owner if scope == "local" else None,
    "session": session, "commit": head or None, "branch": branch or None,
    "ts": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
}
with open(path, "a") as fh:
    fh.write(json.dumps(item) + "\n")
print("%s %s" % (item["id"], scope))
' "$text" "$why" "$files" "$task" "$AGENT_SESSION" "$head" "$branch" "$BACKLOG" "$scope" "$(owner_now)"
}

# Out of this checkout's todo and into the pool. Append-only, like a close.
cmd_promote() {
  [[ -n "${1:-}" ]] || die "usage: backlog.sh promote ID"
  py '
import json, sys, time
path, want, session = sys.argv[1:4]
if not any(json.loads(l)["id"] == want for l in open(path) if l.strip()):
    sys.exit("backlog: no item %s" % want)
with open(path, "a") as fh:
    fh.write(json.dumps({
        "id": want, "scope": "global", "owner": None, "promoted_by": session,
        "promoted": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
    }) + "\n")
print("%s global" % want)
' "$BACKLOG" "$1" "$AGENT_SESSION"
}

# Every open local item of one owner, handed to the pool. `worktree.sh rm`
# calls this: the checkout is going, and its todo would go with it.
cmd_release() {
  [[ -n "${1:-}" ]] || die "usage: backlog.sh release OWNER"
  py '
import json, sys, time
path, owner, session = sys.argv[1:4]
items = {}
for line in open(path):
    line = line.strip()
    if line:
        it = json.loads(line)
        items.setdefault(it["id"], {}).update(it)
stale = [i for i in items.values()
         if i.get("status") == "open" and i.get("scope") == "local"
         and i.get("owner") == owner]
with open(path, "a") as fh:
    for i in stale:
        fh.write(json.dumps({
            "id": i["id"], "scope": "global", "owner": None,
            "promoted_by": session,
            "promoted": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
        }) + "\n")
print(len(stale))
' "$BACKLOG" "$1" "$AGENT_SESSION"
}

# Later lines win, so a close is just another append.
cmd_list() {
  local closed="" files="" view="local"
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --closed)
        closed=1
        shift
        ;;
      --global | --both | --orphaned)
        view="${1#--}"
        shift
        ;;
      --files)
        files="$2"
        shift 2
        ;;
      *) die "unknown flag $1" ;;
    esac
  done

  # What an owner string can still refer to: a worktree the registry knows, a
  # session with a process. Anything else owns items nobody will ever pick up.
  local live wts f
  live="$(pgrep -af 'claude' 2>/dev/null | grep -oE '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | sort -u | tr '\n' ',' || true)"
  live="${live}${AGENT_SESSION}"
  wts=""
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] && wts="$wts$(basename "$f" .json),"
  done

  py '
import json, os, sys
path, closed, want, view, me, live, wts = sys.argv[1:8]
live = set(x for x in live.split(",") if x)
wts = set(x for x in wts.split(",") if x)
items = {}
for line in open(path):
    line = line.strip()
    if not line:
        continue
    it = json.loads(line)
    items.setdefault(it["id"], {}).update(it)
items = list(items.values())
if not closed:
    items = [i for i in items if i.get("status") == "open"]
if want:
    items = [i for i in items if any(want in f for f in (i.get("files") or []))]


# Only what this script wrote as local is local. Every other scope — the old
# shared pool, the per-session items from before the tiers — is global.
def local(i):
    return i.get("scope") == "local"


def mine(i):
    return local(i) and i.get("owner") == me


def orphaned(i):
    if not local(i):
        return False
    owner = i.get("owner") or ""
    if owner.startswith("wt:"):
        return owner[3:] not in wts
    return owner[2:] not in live


if view == "global":
    items = [i for i in items if not local(i)]
elif view == "orphaned":
    items = [i for i in items if orphaned(i)]
elif view == "both":
    items = [i for i in items if mine(i) or not local(i)]
else:
    items = [i for i in items if mine(i)]

if not items:
    hint = {
        "global": "nothing in the global pool",
        "orphaned": "no items abandoned by a finished worktree or session",
        "both": "nothing queued",
    }.get(view, "nothing on the local todo")
    print(hint if not want else "%s against %s" % (hint, want))
    raise SystemExit
# Colour arrives from tty.sh, already empty when it is switched off, so there
# is no second place deciding whether this output is decorated.
DIM, BOLD, RESET = (os.environ.get(k, "") for k in ("C_DIM", "C_BOLD", "C_RESET"))
CYAN, GREEN = os.environ.get("C_CYAN", ""), os.environ.get("C_GREEN", "")

for i in sorted(items, key=lambda x: x["ts"]):
    mark, colour = {
        "open": ("[ ]", ""),
        "done": ("[x]", GREEN),
        "dropped": ("[-]", DIM),
    }.get(i["status"], ("[?]", DIM))
    tag = "me" if mine(i) else ("··" if local(i) else "  ")
    print("%s%s %s %s%s%s  %s%s" % (
        colour, mark, tag, CYAN, i["id"], RESET + colour, i["text"], RESET))
    if i.get("files"):
        print("%s        files: %s%s" % (DIM, ", ".join(i["files"]), RESET))
    if i.get("task"):
        print("%s        from:  %s (%s)%s" % (DIM, i["task"], i["ts"][:10], RESET))
' "$BACKLOG" "$closed" "$files" "$view" "$(owner_now)" "$live" "$wts"
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
  promote)
    shift
    cmd_promote "$@"
    ;;
  release)
    shift
    cmd_release "$@"
    ;;
  owner)
    owner_now
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
    sed -n '2,20p' "${BASH_SOURCE[0]}" | sed 's/^# \?//'
    exit 1
    ;;
esac
