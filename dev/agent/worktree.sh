#!/usr/bin/env bash
#
# One task, one worktree, one whole flow.
#
# Every session that needs to build and verify independently gets its own
# checkout, its own server, its own webapp and its own database. Nothing it
# does is visible to another session until it commits, and nothing another
# session does can change what it is testing.
#
# The primary checkout is slot 0 and is never touched by this script.
#
#   worktree.sh new  NAME [--base REF]   create checkout, DB, ports, env files
#   worktree.sh up   [NAME]              start this worktree's server + webapp
#   worktree.sh down [NAME]              stop them
#   worktree.sh list                     every slot, with what is running
#   worktree.sh info [NAME]              ports, DB and paths for one worktree
#   worktree.sh land NAME [--keep]       rebase onto main, merge, remove
#   worktree.sh rm   NAME [--keep-branch]  stop, drop DB, remove checkout
#
# Ports are derived from the slot, so two worktrees can never collide:
#   API 1043+slot*100 · gin 3001+slot*100 · vite 5173+slot
#
set -euo pipefail

AGENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$AGENT_DIR/common.sh"

MAX_SLOTS=8
WT_DIR="$STATE_DIR/worktrees"
mkdir -p "$WT_DIR"

# The main checkout — the one holding .git, slot 0, and the real .state.
primary_repo() {
  git -C "$REPO_DIR" worktree list --porcelain | awk '/^worktree /{print $2; exit}'
}

# Where worktrees are created: a sibling dir, never inside the repo.
worktrees_root() { echo "$(dirname "$(primary_repo)")/human-worktrees"; }

slot_api() { echo $((1043 + $1 * 100)); }
slot_gin() { echo $((3001 + $1 * 100)); }
slot_vite() { echo $((5173 + $1)); }
slot_db() { echo "$(db_name_base)_wt$1"; }

# The DB the primary is pointed at — what every worktree DB is cloned from.
db_name_base() {
  sed -nE 's#^DB_DSN=.*/([^/?]+).*#\1#p' "$(primary_repo)/server/.env" | head -1
}

meta() { echo "$WT_DIR/$1.json"; }

die() {
  echo "worktree: $*" >&2
  exit 1
}

# read_meta NAME FIELD
read_meta() {
  [[ -f "$(meta "$1")" ]] || die "no worktree named '$1' — try: worktree.sh list"
  json_get "$2" <"$(meta "$1")"
}

# The worktree the caller is standing in, if any.
current_name() {
  local here
  here="$(git rev-parse --show-toplevel 2>/dev/null || true)"
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    [[ "$(json_get path <"$f")" == "$here" ]] && basename "$f" .json && return 0
  done
  return 1
}

resolve_name() {
  if [[ -n "${1:-}" ]]; then
    echo "$1"
  else
    current_name || die "not inside a worktree — name one, or run from its checkout"
  fi
}

# claim_slot NAME — reserve the lowest free slot and print it.
#
# The reservation has to be atomic and it has to happen BEFORE the slow work.
# Picking a slot by reading the registry and only recording it at the end
# leaves seconds in which a second `new` reads the same registry, picks the
# same slot, and writes the same ports and database name into a second
# checkout. Both then look correct and one silently serves the other's data.
#
# noclobber makes the create-or-fail one operation, so exactly one caller wins
# each slot.
claim_slot() {
  local n f
  # A registry entry is a claim even when no reservation file backs it — the
  # worktrees that existed before reservations did, and any restored .state.
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    : >>"$WT_DIR/.slot-$(json_get slot <"$f")"
  done
  for ((n = 1; n <= MAX_SLOTS; n++)); do
    if (
      set -o noclobber
      echo "$1" >"$WT_DIR/.slot-$n" 
    ) 2>/dev/null; then
      echo "$n"
      return 0
    fi
  done
  die "all $MAX_SLOTS slots are in use — remove one first"
}

release_slot() { rm -f "$WT_DIR/.slot-$1"; }

# worktree_rollback NAME PATH DB SLOT — undo a partial `new`.
worktree_rollback() {
  echo "worktree: '$1' failed to build — unwinding" >&2
  [[ -d "$2" ]] && git -C "$(primary_repo)" worktree remove "$2" --force 2>/dev/null
  git -C "$(primary_repo)" branch -D "$1" 2>/dev/null
  pg dropdb --if-exists "$3" 2>/dev/null
  rm -f "$(meta "$1")"
  release_slot "$4"
  exit 1
}

port_busy() { ss -ltn "sport = :$1" 2>/dev/null | grep -q LISTEN; }

port_holder() { ss -ltnp "sport = :$1" 2>/dev/null | sed -nE 's/.*pid=([0-9]+).*/\1/p' | head -1; }

orphan_warning() { # orphan_warning PORT PIDFILE WHAT
  if [[ -f "$2" ]]; then
    echo "$3 already up on $1"
  else
    echo "$3 already up on $1 but unmanaged (pid $(port_holder "$1")) — 'down --force' reclaims it" >&2
  fi
}

pg() { PGPASSWORD="${PGPASSWORD:-postgres}" "$@" -h localhost -U postgres; }

# ---------------------------------------------------------------- new --------

cmd_new() {
  local name="${1:-}" base="HEAD"
  shift || true
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --base)
        base="$2"
        shift 2
        ;;
      *) die "unknown flag $1" ;;
    esac
  done
  [[ -n "$name" ]] || die "usage: worktree.sh new NAME [--base REF]"
  [[ "$name" =~ ^[a-z0-9][a-z0-9-]*$ ]] || die "name must be kebab-case"
  [[ -f "$(meta "$name")" ]] && die "'$name' already exists"

  local primary root slot path api gin vite db srcdb
  primary="$(primary_repo)"
  root="$(worktrees_root)"
  slot="$(claim_slot "$name")"
  path="$root/$name"
  api="$(slot_api "$slot")"
  gin="$(slot_gin "$slot")"
  vite="$(slot_vite "$slot")"
  srcdb="$(db_name_base)"
  db="$(slot_db "$slot")"

  for p in "$api" "$gin" "$vite"; do
    port_busy "$p" && {
      release_slot "$slot"
      die "port $p is already listening — slot $slot is not free after all"
    }
  done

  # Anything that dies past this point unwinds. A checkout left behind carries a
  # server/.env naming the slot's ports and database, so the next `new` to win
  # that slot would share them with a directory nothing tracks.
  # shellcheck disable=SC2064
  trap "worktree_rollback '$name' '$path' '$db' '$slot'" ERR

  local carried
  carried="$(git -C "$primary" status --porcelain | wc -l)"
  if [[ "$carried" -gt 0 ]]; then
    echo "note: $carried uncommitted file(s) in the primary will NOT be in this" >&2
    echo "      worktree — it checks out $base. Commit first if you need them." >&2
  fi

  echo "worktree '$name' → slot $slot (api $api · gin $gin · vite $vite · db $db)"

  mkdir -p "$root"
  git -C "$primary" worktree add -b "$name" "$path" "$base" >/dev/null
  echo "  checkout   $path (branch $name off $base)"

  # Shared identity: token, ids.json and the backlog are one set for all
  # sessions. Scratch inside it is per-session, so nothing stomps.
  rm -rf "$path/dev/agent/.state"
  ln -s "$STATE_DIR" "$path/dev/agent/.state"
  echo "  state      → $STATE_DIR (symlink)"

  # Every file below is gitignored, so none of this shows up in the worktree's
  # git status or fights with a tracked file.
  sed -E \
    -e "s#^HTTP_ADDR=.*#HTTP_ADDR=:$api#" \
    -e "s#^DOMAIN=.*#DOMAIN=localhost:$api#" \
    -e "s#^DB_DSN=.*#DB_DSN=postgres://postgres:postgres@localhost:5432/$db#" \
    "$primary/server/.env" >"$path/server/.env"
  echo "  server/.env  HTTP_ADDR=:$api DB=$db"

  cat >"$path/client/web/unify/public/config.js" <<EOF
// Generated by dev/agent/worktree.sh for worktree '$name' (slot $slot).
window.HumanAPI = 'http://localhost:$api/api'
EOF

  if [[ -f "$primary/client/web/unify/.env.e2e" ]]; then
    sed -E "s#^E2E_BASE_URL=.*#E2E_BASE_URL=http://localhost:$vite#" \
      "$primary/client/web/unify/.env.e2e" >"$path/client/web/unify/.env.e2e"
    echo "  .env.e2e     E2E_BASE_URL=http://localhost:$vite"
  fi

  echo "  database     cloning $srcdb → $db"
  pg createdb "$db"
  pg pg_dump --no-owner --no-privileges "$srcdb" | pg psql -q -d "$db" >/dev/null

  python3 -c '
import json, sys, time
name, path, slot, api, gin, vite, db, base, session, out = sys.argv[1:11]
json.dump({
    "name": name, "path": path, "slot": int(slot), "branch": name, "base": base,
    "api": int(api), "gin": int(gin), "vite": int(vite), "db": db,
    "session": session, "created": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
}, open(out, "w"), indent=2)
' "$name" "$path" "$slot" "$api" "$gin" "$vite" "$db" "$base" "$AGENT_SESSION" "$(meta "$name")"

  trap - ERR

  cat <<EOF

ready. next:
  cd $path
  $path/dev/agent/worktree.sh up $name     # installs deps, starts server + webapp
EOF
}

# ------------------------------------------------------------- up / down -----

cmd_up() {
  local name path api gin vite
  name="$(resolve_name "${1:-}")"
  path="$(read_meta "$name" path)"
  api="$(read_meta "$name" api)"
  gin="$(read_meta "$name" gin)"
  vite="$(read_meta "$name" vite)"
  # Both are gitignored, so a fresh checkout has neither: gin writes its binary
  # into build/ and the Makefile tees its log there.
  mkdir -p "$path/.run" "$path/server/build"

  # setsid forks when the caller is already a group leader, so $! names the
  # wrapper and not the leader. Letting the leader write its own $$ is the only
  # reading that survives that.
  start_svc() { # start_svc PIDFILE LOGFILE DIR CMD...
    local pidfile="$1" logfile="$2" dir="$3"
    shift 3
    (cd "$dir" && setsid bash -c 'echo $$ >"$1"; shift; exec "$@"' _ "$pidfile" "$@" \
      >"$logfile" 2>&1 &)
  }

  if [[ ! -d "$path/node_modules" ]]; then
    echo "installing deps (first run in this worktree) …"
    (cd "$path" && pnpm install --silent)
  fi

  if port_busy "$gin"; then
    orphan_warning "$gin" "$path/.run/server.pid" server
  else
    start_svc "$path/.run/server.pid" "$path/.run/server.log" "$path/server" \
      make watch GIN_ARG_PORT="$gin"
    echo "server   gin :$gin → api :$api   (log $path/.run/server.log)"
  fi

  if port_busy "$vite"; then
    orphan_warning "$vite" "$path/.run/webapp.pid" webapp
  else
    start_svc "$path/.run/webapp.pid" "$path/.run/webapp.log" \
      "$path/client/web/unify" pnpm dev --port "$vite" --strictPort
    echo "webapp   http://localhost:$vite   (log $path/.run/webapp.log)"
  fi

  echo
  echo "gin builds on first request to its proxy, so warm it before believing a check:"
  echo "  curl -s localhost:$gin/api/ >/dev/null"
}

cmd_down() {
  local name path force=""
  name="$(resolve_name "${1:-}")"
  [[ "${2:-}" == "--force" || "${1:-}" == "--force" ]] && force=1
  [[ "${1:-}" == "--force" ]] && name="$(resolve_name "")"
  path="$(read_meta "$name" path)"
  for what in server webapp; do
    local pidfile="$path/.run/$what.pid"
    [[ -f "$pidfile" ]] || continue
    local pid
    pid="$(cat "$pidfile")"
    # setsid made it a group leader, so the negative PID reaches the whole
    # tree. Killing the wrapper alone leaves the port bound.
    if kill -0 "$pid" 2>/dev/null; then
      kill -TERM -"$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null || true
      for _ in 1 2 3 4 5 6 7 8 9 10; do
        kill -0 "$pid" 2>/dev/null || break
        sleep 0.3
      done
      kill -KILL -"$pid" 2>/dev/null || true
      echo "stopped $what ($pid)"
    fi
    rm -f "$pidfile"
  done

  # The slot owns these ports, so anything still holding one is this
  # worktree's orphan — but only --force reaches for a pid nothing recorded.
  local stuck=0 holder
  for p in "$(read_meta "$name" gin)" "$(read_meta "$name" vite)" "$(read_meta "$name" api)"; do
    port_busy "$p" || continue
    holder="$(port_holder "$p")"
    if [[ -n "$force" && -n "$holder" ]]; then
      kill -TERM -"$(ps -o pgid= -p "$holder" | tr -d ' ')" 2>/dev/null ||
        kill -TERM "$holder" 2>/dev/null || true
      echo "forced port $p (pid $holder)"
    else
      echo "port $p still bound by pid ${holder:-?} — 'down $name --force' reclaims it" >&2
      stuck=1
    fi
  done
  [[ "$stuck" -eq 1 ]] && return 1
  return 0
}

# ------------------------------------------------------------ list / info ----

cmd_list() {
  printf '%-4s %-18s %-6s %-6s %-6s %-24s %s\n' SLOT NAME API GIN VITE DB STATUS
  printf '%-4s %-18s %-6s %-6s %-6s %-24s %s\n' 0 '(primary)' 1043 3001 5173 \
    "$(db_name_base)" "$(port_busy 3001 && echo up || echo down)"
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    local n s
    n="$(json_get name <"$f")"
    s=""
    port_busy "$(json_get gin <"$f")" && s="server "
    port_busy "$(json_get vite <"$f")" && s="${s}webapp"
    printf '%-4s %-18s %-6s %-6s %-6s %-24s %s\n' \
      "$(json_get slot <"$f")" "$n" "$(json_get api <"$f")" \
      "$(json_get gin <"$f")" "$(json_get vite <"$f")" \
      "$(json_get db <"$f")" "${s:-down}"
  done
  report_orphans
}

# A checkout under the worktrees root with no metadata is unmanaged: `rm` cannot
# see it, and its server/.env names a slot someone else may now own.
report_orphans() {
  local root primary p n found=""
  primary="$(primary_repo)"
  root="$(worktrees_root)"
  [[ -d "$root" ]] || return 0
  for p in "$root"/*; do
    [[ -d "$p" ]] || continue
    n="$(basename "$p")"
    [[ -f "$(meta "$n")" ]] && continue
    [[ -z "$found" ]] && echo && found=1
    echo "orphan: $p has no registry entry — 'git worktree remove $p --force'" >&2
  done
}

cmd_info() {
  local name
  name="$(resolve_name "${1:-}")"
  cat "$(meta "$name")"
}

# --------------------------------------------------------------- land --------

# Put a worktree's commits on main and take the worktree away.
#
# Everything here is a refusal looking for a reason. main lives in the primary's
# working tree, which on this machine usually holds another session's
# uncommitted work — merging into it blind is how someone else's afternoon gets
# a conflict it did not ask for. So: check first, mutate second.
cmd_land() {
  local name="${1:-}" keep=""
  shift || true
  [[ "${1:-}" == "--keep" ]] && keep=1
  [[ -n "$name" ]] || die "usage: worktree.sh land NAME [--keep]"

  local path primary branch_head
  path="$(read_meta "$name" path)"
  primary="$(primary_repo)"

  branch_head="$(git -C "$primary" symbolic-ref --short HEAD 2>/dev/null || echo '')"
  [[ "$branch_head" == "main" ]] ||
    die "the primary is on '$branch_head', not main — land merges into whatever main is checked out as"

  local dirty
  dirty="$(git -C "$path" status --porcelain | grep -v 'dev/agent/\.state$' || true)"
  [[ -n "$dirty" ]] && {
    echo "$dirty" >&2
    die "'$name' has uncommitted work — commit it in the worktree first"
  }

  local ahead
  ahead="$(git -C "$primary" rev-list --count "main..$name")"
  [[ "$ahead" -eq 0 ]] && die "'$name' has no commits main does not already have"

  # The refusal that matters: the branch's files against the primary's dirt.
  local theirs mine overlap
  theirs="$(git -C "$primary" diff --name-only "main...$name")"
  mine="$(git -C "$primary" status --porcelain | awk '{print $2}')"
  overlap="$(comm -12 <(echo "$theirs" | sort -u) <(echo "$mine" | sort -u))"
  [[ -n "$overlap" ]] && {
    echo "the primary has uncommitted changes in files this branch also touches:" >&2
    echo "$overlap" | sed 's/^/  /' >&2
    die "landing would collide with work already in the primary — resolve that first"
  }

  echo "landing '$name' ($ahead commit(s))"
  git -C "$path" rebase main >/dev/null 2>&1 || {
    git -C "$path" rebase --abort 2>/dev/null || true
    die "'$name' does not rebase cleanly onto main — resolve it in $path"
  }
  echo "  rebased onto main"

  git -C "$primary" merge --ff-only "$name" >/dev/null ||
    die "fast-forward refused — main moved again; re-run land"
  echo "  merged (fast-forward, no merge commit)"

  if [[ -n "$keep" ]]; then
    echo "  worktree kept — it is now level with main"
  else
    cmd_rm "$name"
  fi
}

# ----------------------------------------------------------------- rm --------

cmd_rm() {
  local name="${1:-}" keep_branch=""
  shift || true
  [[ "${1:-}" == "--keep-branch" ]] && keep_branch=1
  [[ -n "$name" ]] || die "usage: worktree.sh rm NAME [--keep-branch]"

  local path db slot primary
  path="$(read_meta "$name" path)"
  db="$(read_meta "$name" db)"
  slot="$(read_meta "$name" slot)"
  primary="$(primary_repo)"

  # The .state symlink is this script's own doing, not the human's work.
  local dirty
  dirty="$(git -C "$path" status --porcelain 2>/dev/null |
    grep -v 'dev/agent/\.state$' | head -5 || true)"
  if [[ -n "$dirty" ]]; then
    echo "worktree '$name' still has uncommitted work:" >&2
    echo "$dirty" >&2
    die "commit it or remove the checkout by hand — this script will not discard it"
  fi

  # Unmerged commits are the one thing here that cannot be rebuilt, so say what
  # happened to them. `branch -d` refuses to delete unmerged work, and swallowing
  # that refusal is how a clean-looking "removed" hides a branch still holding
  # the only copy of a commit.
  local ahead
  ahead="$(git -C "$primary" rev-list --count "main..$name" 2>/dev/null || echo 0)"

  cmd_down "$name" || true
  git -C "$primary" worktree remove "$path" --force
  pg dropdb --if-exists "$db"
  rm -f "$(meta "$name")"
  release_slot "$slot"

  if [[ -n "$keep_branch" ]]; then
    echo "removed '$name' (checkout, database $db, slot freed); branch kept"
  elif [[ "$ahead" -gt 0 ]]; then
    echo "removed '$name' (checkout, database $db, slot freed)"
    echo
    echo "branch '$name' KEPT — it holds $ahead commit(s) that are not on main:" >&2
    git -C "$primary" log --oneline "main..$name" >&2
    echo >&2
    echo "  land them:  git merge $name        (from $primary)" >&2
    echo "  or discard: git branch -D $name" >&2
  else
    git -C "$primary" branch -d "$name" >/dev/null 2>&1 || true
    echo "removed '$name' (checkout, database $db, branch, slot freed)"
  fi
}

case "${1:-}" in
  new)
    shift
    cmd_new "$@"
    ;;
  up)
    shift
    cmd_up "$@"
    ;;
  down)
    shift
    cmd_down "$@"
    ;;
  list)
    shift
    cmd_list "$@"
    ;;
  info)
    shift
    cmd_info "$@"
    ;;
  land)
    shift
    cmd_land "$@"
    ;;
  rm)
    shift
    cmd_rm "$@"
    ;;
  *)
    sed -n '2,20p' "${BASH_SOURCE[0]}" | sed 's/^# \?//'
    exit 1
    ;;
esac
