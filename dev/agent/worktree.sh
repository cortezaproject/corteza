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
#   worktree.sh gc   [--reap]            find abandoned worktrees and databases
#   worktree.sh rm   NAME [--keep-branch]  stop, drop DB, remove checkout
#
# Ports are derived from the slot, off whatever the primary serves, so two
# worktrees can never collide:
#   API <primary api>+slot*100 · vite <primary vite>+slot
#
set -euo pipefail

AGENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$AGENT_DIR/common.sh"

# A worktree's dev/agent/.state is a symlink into the primary, and this script
# runs from the worktree's own copy as often as the primary's. Resolve it to
# the directory it points at: every path under the symlink dies with the
# checkout that `rm` deletes, and `rm -f` on a vanished path succeeds, so the
# ledger entry and the slot outlive a removal that reported freeing them.
STATE_DIR="$(cd "$STATE_DIR" && pwd -P)"

MAX_SLOTS=8
WT_DIR="$STATE_DIR/worktrees"
mkdir -p "$WT_DIR"

# The main checkout — the one holding .git, slot 0, and the real .state.
primary_repo() {
  git -C "$REPO_DIR" worktree list --porcelain | awk '/^worktree /{print $2; exit}'
}

# Where worktrees are created: a sibling dir, never inside the repo.
worktrees_root() { echo "$(dirname "$(primary_repo)")/human-worktrees"; }

# Slot 0 is the primary, and its ports are the developer's rather than this
# script's — so the lanes are numbered off what the primary actually serves.
# Reading it through stack.sh is what keeps the two answers the same.
read -r PRIMARY_API PRIMARY_VITE <<<"$(
  HUMAN_API= HUMAN_WEBAPP= STACK_ROOT="$(primary_repo)" bash "$AGENT_DIR/stack.sh" |
    sed -nE 's#^HUMAN_(API|WEBAPP)=https?://[^:]+:([0-9]+).*#\2#p' | tr '\n' ' '
)"

slot_api() { echo $((PRIMARY_API + $1 * 100)); }
slot_vite() { echo $((PRIMARY_VITE + $1)); }
# Corredor's gRPC port: the slot's API port moved up by 50000 (1343 → 51343).
slot_corredor() { echo $((50000 + $(slot_api "$1"))); }

# corredor_port NAME — from the metadata, or derived for a worktree made
# before the metadata carried it.
corredor_port() {
  local p
  p="$(json_get corredor <"$(meta "$1")" 2>/dev/null || true)"
  [[ -n "$p" ]] || p="$(slot_corredor "$(read_meta "$1" slot)")"
  echo "$p"
}

# append_line FILE LINE — on a line of its own even when the file has no
# trailing newline (the primary's .env ends mid-line).
append_line() {
  [[ ! -s "$1" || "$(tail -c1 "$1" | od -An -c | tr -d ' ')" == '\n' ]] || echo >>"$1"
  printf '%s\n' "$2" >>"$1"
}

# set_env FILE KEY VALUE — the last uncommented KEY= line is the one the
# server reads (see stack.sh), so an existing one is rewritten in place.
set_env() {
  if grep -qE "^[[:space:]]*$2=" "$1"; then
    sed -i -E "s#^[[:space:]]*$2=.*#$2=$3#" "$1"
  else
    append_line "$1" "$2=$3"
  fi
}
slot_db() { echo "$(db_name_base)_wt$1"; }

# The primary's DB_DSN — the one connection string every worktree derives from.
# Only the database name is swapped per slot: user, password, host and port are
# the developer's, and a literal here works on one machine and not the next.
primary_dsn() {
  sed -nE 's/^[[:space:]]*DB_DSN=[[:space:]]*//p' "$(primary_repo)/server/.env" |
    tail -1 | sed -E 's/^["'"'"']//; s/["'"'"']$//'
}

# dsn_part FIELD [DSN] — user, password, host, port, db or scheme.
dsn_part() {
  python3 -c '
import sys, urllib.parse
field, dsn = sys.argv[1], sys.argv[2]
u = urllib.parse.urlsplit(dsn)
print({
    "scheme": u.scheme,
    "user": u.username or "",
    "password": u.password or "",
    "host": u.hostname or "localhost",
    "port": str(u.port or 5432),
    "db": u.path.lstrip("/").split("?")[0],
}[field])
' "$1" "${2:-$(primary_dsn)}"
}

# dsn_with_db NAME — the primary's DSN pointed at a different database.
dsn_with_db() {
  python3 -c '
import sys, urllib.parse
dsn, db = sys.argv[1], sys.argv[2]
u = urllib.parse.urlsplit(dsn)
print(urllib.parse.urlunsplit((u.scheme, u.netloc, "/" + db, u.query, u.fragment)))
' "$(primary_dsn)" "$1"
}

# The DB the primary is pointed at — what every worktree DB is cloned from.
db_name_base() { dsn_part db; }

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

# Signal whatever is listening on a port, by its own process group — a server
# started under a watcher is a group leader in its own right, and reaching it
# through the group that started the watcher does not work.
stop_port_group() { # stop_port_group PORT SIGNAL
  local holder pgid
  holder="$(port_holder "$1")"
  [[ -n "$holder" ]] || return 0
  pgid="$(ps -o pgid= -p "$holder" 2>/dev/null | tr -d ' ')"
  [[ -n "$pgid" ]] && kill -"$2" -"$pgid" 2>/dev/null ||
    kill -"$2" "$holder" 2>/dev/null || true
}

orphan_warning() { # orphan_warning PORT PIDFILE WHAT
  if [[ -f "$2" ]]; then
    echo "  $(paint "$C_DIM" "$G_OK $3 already up on $1")"
  else
    warn "$3 already up on $1 but unmanaged (pid $(port_holder "$1")) — 'down --force' reclaims it"
  fi
}

# Creating and cloning a database is not something the app's own role can
# always do, so the identity is overridable — but it defaults to the primary's
# rather than to a literal `postgres`, which is a superuser on one machine and
# absent on the next.
pg() {
  PGPASSWORD="${PGSUPERPASS:-$(dsn_part password)}" \
    "$@" \
    -h "${PGSUPERHOST:-$(dsn_part host)}" \
    -p "${PGSUPERPORT:-$(dsn_part port)}" \
    -U "${PGSUPERUSER:-$(dsn_part user)}"
}

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
  case "$(dsn_part scheme)" in
    postgres | postgresql) ;;
    *) die "the primary's DB_DSN is '$(dsn_part scheme)'; a worktree is a cloned postgres database, so there is nothing to clone from" ;;
  esac
  [[ "$name" =~ ^[a-z0-9][a-z0-9-]*$ ]] || die "name must be kebab-case"
  [[ -f "$(meta "$name")" ]] && die "'$name' already exists"

  local primary root slot path api vite db srcdb
  primary="$(primary_repo)"
  root="$(worktrees_root)"
  slot="$(claim_slot "$name")"
  path="$root/$name"
  api="$(slot_api "$slot")"
  vite="$(slot_vite "$slot")"
  srcdb="$(db_name_base)"
  db="$(slot_db "$slot")"

  for p in "$api" "$vite"; do
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

  echo
  echo "$(paint "$C_BOLD" "$G_SECTION worktree '$name'") $(paint "$C_DIM" "slot $slot · api $api · vite $vite · db $db")"

  mkdir -p "$root"
  git -C "$primary" worktree add -b "$name" "$path" "$base" >/dev/null
  step checkout "$path $(paint "$C_DIM" "(branch $name off $base)")"

  # Shared identity: token, ids.json and the backlog are one set for all
  # sessions. Scratch inside it is per-session, so nothing stomps.
  rm -rf "$path/dev/agent/.state"
  ln -s "$STATE_DIR" "$path/dev/agent/.state"
  step state "→ $STATE_DIR $(paint "$C_DIM" "(symlink)")"

  # Every file below is gitignored, so none of this shows up in the worktree's
  # git status or fights with a tracked file.
  sed -E \
    -e "s#^HTTP_ADDR=.*#HTTP_ADDR=:$api#" \
    -e "s#^DOMAIN=.*#DOMAIN=localhost:$api#" \
    -e "s#^DOMAIN_WEBAPP=.*#DOMAIN_WEBAPP=localhost:$vite#" \
    -e "s#^DB_DSN=.*#DB_DSN=$(dsn_with_db "$db")#" \
    "$primary/server/.env" >"$path/server/.env"

  # The webapp is on vite, not on the API port. DOMAIN_WEBAPP is the only thing
  # that says so, and server/pkg/weburl builds every link a tool result hands
  # back from it; without it they all point at the API and 404.
  grep -q '^DOMAIN_WEBAPP=' "$path/server/.env" ||
    append_line "$path/server/.env" "DOMAIN_WEBAPP=localhost:$vite"
  step server/.env "HTTP_ADDR=:$api DOMAIN_WEBAPP=localhost:$vite DB=$db"

  cat >"$path/client/web/unify/public/config.js" <<EOF
// Generated by dev/agent/worktree.sh for worktree '$name' (slot $slot).
window.HumanAPI = 'http://localhost:$api/api'
EOF

  # Where vite binds. vite.config.js reads VITE_PORT with strictPort on, so a
  # lane whose port is taken fails at startup instead of drifting onto the next
  # one while .env.e2e and the toolkit keep naming this one.
  cat >"$path/client/web/unify/.env.local" <<EOF
# Generated by dev/agent/worktree.sh for worktree '$name' (slot $slot).
VITE_PORT=$vite
EOF
  step .env.local "VITE_PORT=$vite"

  if [[ -f "$primary/client/web/unify/.env.e2e" ]]; then
    sed -E "s#^E2E_BASE_URL=.*#E2E_BASE_URL=http://localhost:$vite#" \
      "$primary/client/web/unify/.env.e2e" >"$path/client/web/unify/.env.e2e"
    step .env.e2e "E2E_BASE_URL=http://localhost:$vite"
  fi

  step database "cloning $srcdb → $db"
  pg createdb "$db"
  pg pg_dump --no-owner --no-privileges "$srcdb" | pg psql -q -d "$db" >/dev/null

  python3 -c '
import json, sys, time
name, path, slot, api, vite, corredor, db, base, session, out = sys.argv[1:11]
json.dump({
    "name": name, "path": path, "slot": int(slot), "branch": name, "base": base,
    "api": int(api), "vite": int(vite), "corredor": int(corredor), "db": db,
    "session": session, "created": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
}, open(out, "w"), indent=2)
' "$name" "$path" "$slot" "$api" "$vite" "$(slot_corredor "$slot")" "$db" "$base" "$AGENT_SESSION" "$(meta "$name")"

  trap - ERR

  cat <<EOF

ready. next:
  cd $path
  $path/dev/agent/worktree.sh up $name     # installs deps, starts server + webapp
EOF
}

# ------------------------------------------------------------- up / down -----

cmd_up() {
  local name path api vite
  name="$(resolve_name "${1:-}")"
  path="$(read_meta "$name" path)"
  api="$(read_meta "$name" api)"
  vite="$(read_meta "$name" vite)"
  # Both are gitignored, so a fresh checkout has neither: the watcher writes its
  # binary into build/ and the Makefile tees its log there.
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
    note "installing deps (first run in this worktree) …"
    (cd "$path" && pnpm install --silent)
  fi

  # A first build is minutes. Doing it here reports that, where leaving it to
  # the watcher would have `up` announce a server and return while the port is
  # still dead.
  if [[ ! -x "$path/server/build/dev-bin" ]]; then
    note "building the server (first run in this worktree) …"
    (cd "$path/server" && go build -o build/dev-bin ./cmd/human)
  fi

  # Corredor listens before the server boots: with CORREDOR_ENABLED=true the
  # server's Connect() is a boot failure when nothing answers. A checkout
  # without corredor/ leaves the server's Corredor settings as they are.
  if [[ -f "$path/corredor/package.json" ]]; then
    local cport
    cport="$(corredor_port "$name")"
    set_env "$path/server/.env" CORREDOR_ENABLED true
    set_env "$path/server/.env" CORREDOR_ADDR "localhost:$cport"
    mkdir -p "$path/.run/corredor-bundles"
    if port_busy "$cport"; then
      orphan_warning "$cport" "$path/.run/corredor.pid" corredor
    else
      start_svc "$path/.run/corredor.pid" "$path/.run/corredor.log" "$path/corredor" \
        env CORREDOR_ENVIRONMENT=dev \
        CORREDOR_ADDR="localhost:$cport" \
        CORREDOR_SERVER_CERTIFICATES_ENABLED=false \
        CORREDOR_LOG_PRETTY=true \
        CORREDOR_EXT_SEARCH_PATHS="$path/dev/fixtures/corredor:$path/corredor/usr:$path/corredor/usr/*" \
        CORREDOR_EXT_DEPENDENCIES_AUTO_UPDATE=false \
        CORREDOR_BUNDLER_OUTPUT_PATH="$path/.run/corredor-bundles" \
        CORREDOR_EXEC_CSERVERS_API_HOST="localhost:$api" \
        CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE='http://{host}/api/{service}' \
        CORREDOR_EXEC_CTX_FRONTEND_BASEURL="http://localhost:$vite" \
        pnpm serve
      for ((i = 0; i < 100; i++)); do
        port_busy "$cport" && break
        sleep 0.3
      done
      if port_busy "$cport"; then
        printf '  %s %-8s %s\n' "$(paint "$C_GREEN" "$G_OK")" corredor \
          "grpc :$cport   $(paint "$C_DIM" "(log $path/.run/corredor.log)")"
      else
        warn "corredor did not come up on :$cport — the server will fail to boot; see $path/.run/corredor.log"
      fi
    fi
  fi

  if port_busy "$api"; then
    orphan_warning "$api" "$path/.run/server.pid" server
  else
    start_svc "$path/.run/server.pid" "$path/.run/server.log" "$path/server" \
      make watch
    printf '  %s %-8s %s\n' "$(paint "$C_GREEN" "$G_OK")" server \
      "api :$api   $(paint "$C_DIM" "(log $path/.run/server.log)")"
  fi

  if port_busy "$vite"; then
    orphan_warning "$vite" "$path/.run/webapp.pid" webapp
  else
    start_svc "$path/.run/webapp.pid" "$path/.run/webapp.log" \
      "$path/client/web/unify" pnpm dev
    printf '  %s %-8s %s\n' "$(paint "$C_GREEN" "$G_OK")" webapp \
      "$(paint "$C_CYAN" "http://localhost:$vite")   $(paint "$C_DIM" "(log $path/.run/webapp.log)")"
  fi

  echo
  note "The server rebuilds and restarts itself on any .go write — give it ~15s."
}

cmd_down() {
  local name path force=""
  name="$(resolve_name "${1:-}")"
  [[ "${2:-}" == "--force" || "${1:-}" == "--force" ]] && force=1
  [[ "${1:-}" == "--force" ]] && name="$(resolve_name "")"
  path="$(read_meta "$name" path)"
  local api vite cport
  api="$(read_meta "$name" api)"
  vite="$(read_meta "$name" vite)"
  cport="$(corredor_port "$name")"
  for what in server webapp corredor; do
    local pidfile="$path/.run/$what.pid" port
    [[ -f "$pidfile" ]] || continue
    case "$what" in
      server) port="$api" ;;
      webapp) port="$vite" ;;
      corredor) port="$cport" ;;
    esac
    local pid
    pid="$(cat "$pidfile")"
    # setsid made it a group leader, so the negative PID reaches the whole
    # tree. Killing the wrapper alone leaves the port bound.
    if kill -0 "$pid" 2>/dev/null; then
      # The watcher runs the server in a process group of its own, so the
      # wrapper's group never contains it. Nor does the watcher get to stop it
      # on the way out: it is piped into `tee`, which the same group signal
      # kills, and the next line it logs takes it down before its shutdown
      # runs. So ask the listener to close first, by its own group, while it
      # still has a parent to be reaped by.
      stop_port_group "$port" TERM
      kill -TERM -"$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null || true
      # The port is what `up` checks, and it outlives the group leader —
      # `make` exits at once, so waiting on that alone puts the SIGKILL below
      # in the middle of the server's shutdown.
      for ((i = 0; i < 40; i++)); do
        kill -0 "$pid" 2>/dev/null || port_busy "$port" || break
        sleep 0.3
      done
      kill -KILL -"$pid" 2>/dev/null || true
      # Last resort for a service this slot started: the wait above already
      # gave it a full shutdown window.
      if port_busy "$port"; then stop_port_group "$port" KILL; fi
      echo "  $(paint "$C_GREEN" "$G_OK") stopped $what $(paint "$C_DIM" "($pid)")"
    fi
    rm -f "$pidfile"
  done

  # The slot owns these ports, so anything still holding one is this
  # worktree's orphan — but only --force reaches for a pid nothing recorded.
  local stuck=0 holder
  for p in "$(read_meta "$name" vite)" "$(read_meta "$name" api)" "$cport"; do
    port_busy "$p" || continue
    holder="$(port_holder "$p")"
    if [[ -n "$force" && -n "$holder" ]]; then
      kill -TERM -"$(ps -o pgid= -p "$holder" | tr -d ' ')" 2>/dev/null ||
        kill -TERM "$holder" 2>/dev/null || true
      ok "forced port $p $(paint "$C_DIM" "(pid $holder)")"
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
  printf '%s\n' "$(paint "$C_BOLD" \
    "$(printf '%-4s %-18s %-6s %-6s %-24s %s' SLOT NAME API VITE DB STATUS)")"
  printf '%-4s %-18s %-6s %-6s %-24s %s\n' 0 '(primary)' "$PRIMARY_API" "$PRIMARY_VITE" \
    "$(db_name_base)" "$(port_busy "$PRIMARY_API" &&
      paint "$C_GREEN" up || paint "$C_DIM" down)"
  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    local n s
    n="$(json_get name <"$f")"
    s=""
    port_busy "$(json_get api <"$f")" && s="server "
    port_busy "$(json_get vite <"$f")" && s="${s}webapp "
    port_busy "$(corredor_port "$n")" && s="${s}corredor"
    s="${s% }"
    printf '%-4s %-18s %-6s %-6s %-24s %s\n' \
      "$(json_get slot <"$f")" "$n" "$(json_get api <"$f")" \
      "$(json_get vite <"$f")" \
      "$(json_get db <"$f")" \
      "$([[ -n "$s" ]] && paint "$C_GREEN" "$s" || paint "$C_DIM" down)"
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

  section "landing '$name'" && echo "  $(paint "$C_DIM" "$ahead commit(s)")"
  git -C "$path" rebase main >/dev/null 2>&1 || {
    git -C "$path" rebase --abort 2>/dev/null || true
    die "'$name' does not rebase cleanly onto main — resolve it in $path"
  }
  step rebased "onto main"

  git -C "$primary" merge --ff-only "$name" >/dev/null ||
    die "fast-forward refused — main moved again; re-run land"
  step merged "$(paint "$C_DIM" "fast-forward, no merge commit")"

  if [[ -n "$keep" ]]; then
    note "worktree kept — it is now level with main"
  else
    cmd_rm "$name"
  fi
}

# ------------------------------------------------------------------ gc -------

# Find what nothing is going to clean up on its own.
#
# Closing a tab reaps none of this: the checkout, the branch, the database and
# the slot all survive, and the servers are their own process group so they
# outlive the terminal too. Only rm and land remove anything, and a session
# that ends mid-task calls neither.
#
# Reports by default. --reap removes only what cannot lose anything: an entry
# whose checkout is gone, a database no entry claims, and a worktree that is
# clean, fully merged and not serving. Uncommitted work is never touched — it
# is the only thing here that exists nowhere else.
cmd_gc() {
  local reap=""
  [[ "${1:-}" == "--reap" ]] && reap=1

  local primary base found=0 held=0
  primary="$(primary_repo)"
  base="$(db_name_base)"

  for f in "$WT_DIR"/*.json; do
    [[ -f "$f" ]] || continue
    local name path db slot dirty ahead running
    name="$(json_get name <"$f")"
    path="$(json_get path <"$f")"
    db="$(json_get db <"$f")"
    slot="$(json_get slot <"$f")"

    # Checkout gone: nothing left to lose, only residue to drop.
    if [[ ! -d "$path" ]]; then
      found=$((found + 1))
      if [[ -n "$reap" ]]; then
        pg dropdb --if-exists "$db" 2>/dev/null || true
        git -C "$primary" worktree prune 2>/dev/null || true
        rm -f "$f"
        release_slot "$slot"
        ok "reaped  $name $(paint "$C_DIM" "— checkout was already gone (database $db, slot $slot)")"
      else
        warn "stale   $name — checkout gone, database $db and slot $slot still held"
      fi
      continue
    fi

    running=""
    port_busy "$(json_get api <"$f")" && running="serving"
    port_busy "$(json_get vite <"$f")" && running="serving"
    port_busy "$(corredor_port "$name")" && running="serving"

    dirty="$(git -C "$path" status --porcelain 2>/dev/null |
      grep -v 'dev/agent/\.state$' || true)"
    ahead="$(git -C "$primary" rev-list --count "main..$name" 2>/dev/null || echo 0)"

    if [[ -n "$dirty" || "$ahead" -gt 0 || -n "$running" ]]; then
      held=$((held + 1))
      echo "HOLD    $name — $(
        [[ -n "$running" ]] && printf 'servers up; '
        [[ "$ahead" -gt 0 ]] && printf '%s unmerged commit(s); ' "$ahead"
        [[ -n "$dirty" ]] && printf '%s uncommitted file(s); ' "$(echo "$dirty" | wc -l)"
        true
      )not touched"
      [[ -n "$dirty" ]] && echo "$dirty" | sed 's/^/          /'
      [[ "$ahead" -gt 0 ]] && git -C "$primary" log --oneline "main..$name" | sed 's/^/          /'
      echo "          $path  ·  http://localhost:$(json_get vite <"$f")"
      continue
    fi

    found=$((found + 1))
    if [[ -n "$reap" ]]; then
      cmd_rm "$name" >/dev/null && echo "reaped  $name — clean and fully merged"
    else
      ok "reapable $name $(paint "$C_DIM" "— clean, fully merged, not serving")"
    fi
  done

  # A database whose worktree nobody records. Costs disk and a name forever.
  local d
  while read -r d; do
    [[ -n "$d" ]] || continue
    local claimed=""
    for f in "$WT_DIR"/*.json; do
      [[ -f "$f" ]] || continue
      [[ "$(json_get db <"$f")" == "$d" ]] && claimed=1 && break
    done
    [[ -n "$claimed" ]] && continue
    found=$((found + 1))
    if [[ -n "$reap" ]]; then
      pg dropdb --if-exists "$d" && echo "reaped  database $d — no worktree claims it"
    else
      warn "orphan  database $d — no worktree claims it"
    fi
  done < <(pg psql -tAc "select datname from pg_database where datname like '${base}_wt%'" 2>/dev/null || true)

  report_orphans

  if [[ "$found" -eq 0 && "$held" -eq 0 ]]; then
    ok "nothing to collect"
  elif [[ -z "$reap" && "$found" -gt 0 ]]; then
    echo
    echo "$found item(s) safe to remove — 'worktree.sh gc --reap'"
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

  # The checkout owns its todo, so what is still open on it goes to the global
  # pool rather than out with the directory.
  local released
  released="$("$primary/dev/agent/backlog.sh" release "wt:$name" 2>/dev/null || true)"
  [[ "$released" =~ ^[0-9]+$ ]] || released=0
  if [[ "$released" -gt 0 ]]; then
    step backlog "$released open item(s) $(paint "$C_DIM" "→ global pool")"
  fi

  cmd_down "$name" || true
  git -C "$primary" worktree remove "$path" --force
  pg dropdb --if-exists "$db"
  rm -f "$(meta "$name")"
  release_slot "$slot"

  # The ledger entry and the slot are the two things a removal can silently not
  # do. Check them at the primary's real path rather than through $WT_DIR: if
  # that resolved to the deleted checkout, a check written against it asks the
  # same broken question and answers "gone" about a file that is still there.
  # $primary was resolved before the checkout was removed; resolving it again
  # here would run git inside the directory this function just deleted.
  local real_wt residue=""
  real_wt="$primary/dev/agent/.state/worktrees"
  if [[ -e "$real_wt/$name.json" ]]; then residue="ledger entry $real_wt/$name.json"; fi
  if [[ -e "$real_wt/.slot-$slot" ]]; then residue="${residue:+$residue and }slot $slot"; fi
  if [[ -n "$residue" ]]; then
    die "checkout and database $db are gone, but $residue survived — the run cannot report this removed"
  fi

  if [[ -n "$keep_branch" ]]; then
    ok "removed '$name' $(paint "$C_DIM" "(checkout, database $db, slot freed); branch kept")"
  elif [[ "$ahead" -gt 0 ]]; then
    ok "removed '$name' $(paint "$C_DIM" "(checkout, database $db, slot freed)")"
    echo
    echo "branch '$name' KEPT — it holds $ahead commit(s) that are not on main:" >&2
    git -C "$primary" log --oneline "main..$name" >&2
    echo >&2
    echo "  land them:  git merge $name        (from $primary)" >&2
    echo "  or discard: git branch -D $name" >&2
  else
    git -C "$primary" branch -d "$name" >/dev/null 2>&1 || true
    ok "removed '$name' $(paint "$C_DIM" "(checkout, database $db, branch, slot freed)")"
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
  gc)
    shift
    cmd_gc "$@"
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
