#!/usr/bin/env bash
#
# From a fresh clone to a working dev stack, in two commands.
#
#   make setup         everything that does not need a running server
#   make setup-agent   the agent identities, once the stack is up
#   make doctor        the same checks, writing nothing
#
# The split is not cosmetic: dev/agent/bootstrap.sh provisions users over the
# API, so it cannot run before the server does. Everything either phase writes
# is gitignored, and both are idempotent — a second run reports rather than
# repeats.
#
# What it will NOT do is edit a file you already have. server/.env in a working
# checkout is hand-tuned; this creates it when it is missing and otherwise only
# says which of the values the guide depends on are not set.
set -euo pipefail

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
UNIFY_DIR="$REPO_DIR/client/web/unify"
SERVER_ENV="$REPO_DIR/server/.env"
CONFIG_JS="$UNIFY_DIR/public/config.js"
E2E_ENV="$UNIFY_DIR/.env.e2e"
UI_PASSWORD="$REPO_DIR/dev/agent/.state/ui-password"

MODE=setup
[[ "${1:-}" == "--agent" ]] && MODE=agent
[[ "${1:-}" == "--check" ]] && MODE=check

# Where this checkout's stack listens. Re-read after anything writes server/.env
# or .env.e2e, since those files are what it reads.
reload_stack() {
  eval "$(HUMAN_API= HUMAN_BASE= HUMAN_AUTH= HUMAN_WEBAPP= \
    bash "$REPO_DIR/dev/agent/stack.sh")"
}
reload_stack

PROBLEMS=0

say() { printf '  %-12s %s\n' "$1" "$2"; }
bad() {
  printf '  %-12s ! %s\n' "$1" "$2"
  PROBLEMS=$((PROBLEMS + 1))
}
hint() { printf '  %-12s   %s\n' '' "$1"; }
head_() { printf '\n%s\n' "$1"; }

writing() { [[ "$MODE" != check ]]; }

# env_get FILE KEY — the last uncommented assignment, which is the one godotenv
# hands the server. Same rule as dev/agent/stack.sh.
env_get() {
  [[ -f "$1" ]] || return 0
  sed -nE "s/^[[:space:]]*$2=[[:space:]]*(.*)$/\1/p" "$1" | tail -1 |
    sed -E 's/^["'\'']//; s/["'\'']$//'
}

port_of() { # port_of ADDR — ":1043" / "127.0.0.1:1043" / "localhost:1043"
  [[ "$1" == *:* ]] || return 0
  echo "${1##*:}"
}

# ---------------------------------------------------------- prerequisites ----

check_tools() {
  head_ "Prerequisites"

  local missing=0
  for t in go node pnpm python3 curl; do
    if command -v "$t" >/dev/null 2>&1; then
      continue
    fi
    bad "$t" "not on PATH"
    missing=1
  done

  if [[ "$missing" == 0 ]]; then
    local nodemajor pnpmmajor
    nodemajor="$(node -v | sed -E 's/^v([0-9]+).*/\1/')"
    pnpmmajor="$(pnpm -v | cut -d. -f1)"
    say tools "go $(go env GOVERSION | sed 's/^go//') · node $(node -v) · pnpm $(pnpm -v)"
    [[ "$nodemajor" -ge 22 ]] || bad node "22+ required (package.json engines), have $(node -v)"
    [[ "$pnpmmajor" -ge 10 ]] || bad pnpm "10+ required, have $(pnpm -v)"
  fi

  if command -v jq >/dev/null 2>&1; then
    say jq "present"
  elif writing; then
    say jq "installing"
    sudo -n apt-get install -y jq >/dev/null 2>&1 || brew install jq >/dev/null 2>&1 ||
      bad jq "could not auto-install — sudo apt-get install jq"
  else
    bad jq "missing — sudo apt-get install jq"
  fi
}

# ------------------------------------------------------------ node deps ------

install_deps() {
  head_ "Dependencies"

  if writing; then
    (cd "$REPO_DIR" && pnpm install --silent)
    say pnpm "installed"
    if (cd "$UNIFY_DIR" && npx playwright install chromium >/dev/null 2>&1); then
      say playwright "chromium installed"
    else
      bad playwright "chromium install failed — cd client/web/unify && npx playwright install chromium"
    fi
  else
    [[ -d "$REPO_DIR/node_modules" ]] && say pnpm "node_modules present" || bad pnpm "no node_modules — make setup"
  fi
}

# ---------------------------------------------------------- server/.env ------

write_server_env() {
  head_ "server/.env"

  if [[ ! -f "$SERVER_ENV" ]]; then
    if ! writing; then
      bad "server/.env" "missing — make setup creates it from .env.min.example"
      return
    fi
    cp "$REPO_DIR/server/.env.min.example" "$SERVER_ENV"
    cat >>"$SERVER_ENV" <<EOF

# Pinned, not derived. worktree.sh clones this file with a new DB_DSN, and an
# unset secret is derived from the DSN — so each worktree would mint tokens the
# others reject, through the one dev/agent/.state/token they all share.
AUTH_JWT_SECRET=$(head -c 36 /dev/urandom | base64 | tr -d '\n/+=')
EOF
    say "server/.env" "created from .env.min.example"
    reload_stack
  else
    say "server/.env" "exists — checked, not modified"
  fi

  local addr domain base env_ rate dsn
  addr="$(env_get "$SERVER_ENV" HTTP_ADDR)"
  domain="$(env_get "$SERVER_ENV" DOMAIN)"
  base="$(env_get "$SERVER_ENV" HTTP_API_BASE_URL)"
  env_="$(env_get "$SERVER_ENV" ENVIRONMENT)"
  rate="$(env_get "$SERVER_ENV" AUTH_REQUEST_RATE_LIMIT)"
  dsn="$(env_get "$SERVER_ENV" DB_DSN)"

  [[ -n "$dsn" ]] || bad DB_DSN "unset — the server runs on an in-memory database that vanishes on restart"
  [[ -n "$(port_of "$addr")" ]] || bad HTTP_ADDR "no port in '$addr'"

  if [[ "$(port_of "$domain")" != "$(port_of "$addr")" ]]; then
    bad DOMAIN "'$domain' disagrees with HTTP_ADDR '$addr' — logins appear to work, then bounce"
  fi

  [[ "$base" == "/api" ]] ||
    bad HTTP_API_BASE_URL "'$base' — dev/agent/* and the webapp both assume /api"

  [[ "$env_" == "dev" ]] ||
    bad ENVIRONMENT "'$env_' is treated as production — the server then refuses to provision super users"

  [[ "$rate" == "0" ]] ||
    bad AUTH_REQUEST_RATE_LIMIT "'$rate' — a playwright run trips the per-IP auth limit and fails as blank 429 pages"

  [[ -n "$(env_get "$SERVER_ENV" AUTH_JWT_SECRET)" ]] ||
    hint "AUTH_JWT_SECRET unset — derived from DB_DSN, so worktrees will churn the shared token"

  # `make claude` exports the token once and cannot refresh it in flight, so a
  # short lifetime shows up as the local MCP server going quiet mid-session.
  case "$(env_get "$SERVER_ENV" AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME)" in
    [2-9][0-9]h | [0-9][0-9][0-9]h | [0-9][0-9][0-9][0-9]h) ;;
    *) hint "AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME is under a day — the human-local MCP token will expire mid-session (720h suits dev)" ;;
  esac

  say api "$HUMAN_API"
  say webapp "$HUMAN_WEBAPP"
}

# ------------------------------------------------------- public/config.js ----

write_config_js() {
  head_ "client/web/unify/public/config.js"

  if [[ -f "$CONFIG_JS" ]]; then
    # A commented-out alternative is not the value: the shipped example carries
    # one, and reading it reports a healthy checkout as pointed at a remote host.
    local have
    have="$(sed -nE "\#^[[:space:]]*//#d; s/.*HumanAPI[[:space:]]*=[[:space:]]*['\"]([^'\"]+)['\"].*/\1/p" \
      "$CONFIG_JS" | tail -1)"
    if [[ "$have" == "$HUMAN_API" ]]; then
      say config.js "→ $have"
    else
      bad config.js "points at '$have', this checkout's server is $HUMAN_API"
      hint "delete it and re-run, or edit window.HumanAPI"
    fi
    return
  fi

  if ! writing; then
    bad config.js "missing — the webapp cannot reach the API without it"
    return
  fi

  cat >"$CONFIG_JS" <<EOF
// Generated by dev/setup.sh for this checkout.
window.HumanAPI = '$HUMAN_API'
window.i18nPseudoModeEnabled = false
EOF
  say config.js "created → $HUMAN_API"
}

# ------------------------------------------------------------- .env.e2e ------

write_e2e_env() {
  head_ "client/web/unify/.env.e2e"

  local pass='<run make setup-agent>'
  [[ -f "$UI_PASSWORD" ]] && pass="$(cat "$UI_PASSWORD")"

  if [[ ! -f "$E2E_ENV" ]]; then
    if ! writing; then
      bad .env.e2e "missing — make e2e cannot log in"
      return
    fi
    cat >"$E2E_ENV" <<EOF
# Generated by dev/setup.sh. E2E attaches to your running dev stack — it never
# boots or seeds anything. The password is the one bootstrap.sh generated for
# agent@local.dev; make setup-agent writes it here.
E2E_BASE_URL=$HUMAN_WEBAPP
E2E_USER=agent@local.dev
E2E_PASS=$pass
EOF
    say .env.e2e "created → $HUMAN_WEBAPP"
    reload_stack
    return
  fi

  say .env.e2e "exists — checked, not modified"

  local url user
  url="$(env_get "$E2E_ENV" E2E_BASE_URL)"
  user="$(env_get "$E2E_ENV" E2E_USER)"
  if [[ "$url" != "$HUMAN_WEBAPP" ]]; then
    bad E2E_BASE_URL "'$url' — this checkout's webapp is $HUMAN_WEBAPP"
    hint "vite binds VITE_PORT (client/web/unify/.env.local, else .env); this file has to name the same port"
  fi
  case "$(env_get "$E2E_ENV" E2E_PASS)" in
    '' | change-me | '<run make setup-agent>')
      bad E2E_PASS "not filled in — make setup-agent writes it once $user exists"
      ;;
  esac
}

# ------------------------------------------------------------- database ------

# Reads a DSN into DSN_* without shelling out: setup runs before the toolkit is
# usable, and a missing python3 is one of the things it reports.
parse_dsn() {
  local rest creds hostport
  DSN_SCHEME="${1%%://*}"
  rest="${1#*://}"
  rest="${rest%%\?*}"
  DSN_USER=""
  DSN_PASS=""
  if [[ "$rest" == *@* ]]; then
    creds="${rest%%@*}"
    rest="${rest#*@}"
    DSN_USER="${creds%%:*}"
    [[ "$creds" == *:* ]] && DSN_PASS="${creds#*:}"
  fi
  hostport="${rest%%/*}"
  DSN_DB=""
  [[ "$rest" == */* ]] && DSN_DB="${rest#*/}"
  DSN_HOST="${hostport%%:*}"
  DSN_PORT=5432
  [[ "$hostport" == *:* ]] && DSN_PORT="${hostport##*:}"
}

check_database() {
  head_ "Database"

  local dsn
  dsn="$(env_get "$SERVER_ENV" DB_DSN)"

  if [[ -z "$dsn" ]]; then
    bad database "no DB_DSN — in-memory, and everything vanishes on restart"
    return
  fi

  parse_dsn "$dsn"

  if [[ "$DSN_SCHEME" == sqlite3 ]]; then
    say database "sqlite — the server creates the file on boot"
    return
  fi

  if ! command -v psql >/dev/null 2>&1; then
    bad psql "not on PATH — cannot check whether '$DSN_DB' exists"
    hint "sudo -u postgres createuser -P $DSN_USER"
    hint "sudo -u postgres createdb -O $DSN_USER $DSN_DB"
    return
  fi

  if PGPASSWORD="$DSN_PASS" psql -h "$DSN_HOST" -p "$DSN_PORT" -U "$DSN_USER" \
    -d "$DSN_DB" -tAc 'select 1' >/dev/null 2>&1; then
    say database "$DSN_USER@$DSN_HOST:$DSN_PORT/$DSN_DB reachable"
    return
  fi

  if ! writing; then
    bad database "cannot reach $DSN_USER@$DSN_HOST:$DSN_PORT/$DSN_DB"
    hint "sudo -u postgres createuser -P $DSN_USER"
    hint "sudo -u postgres createdb -O $DSN_USER $DSN_DB"
    return
  fi

  # Only ever tried non-interactively: a make target that stops on a hidden
  # sudo password prompt looks like a hang.
  say database "creating role '$DSN_USER' and database '$DSN_DB'"
  if sudo -n -u postgres psql -tAc \
    "select 1 from pg_roles where rolname='$DSN_USER'" 2>/dev/null | grep -q 1 ||
    sudo -n -u postgres createuser "$DSN_USER" 2>/dev/null; then
    sudo -n -u postgres psql -q -c \
      "alter role \"$DSN_USER\" login password '$DSN_PASS'" >/dev/null 2>&1 || true
    sudo -n -u postgres createdb -O "$DSN_USER" "$DSN_DB" 2>/dev/null || true
  fi

  if PGPASSWORD="$DSN_PASS" psql -h "$DSN_HOST" -p "$DSN_PORT" -U "$DSN_USER" \
    -d "$DSN_DB" -tAc 'select 1' >/dev/null 2>&1; then
    say database "$DSN_USER@$DSN_HOST:$DSN_PORT/$DSN_DB created"
  else
    bad database "could not create it without a password-less sudo — run these yourself:"
    hint "sudo -u postgres createuser -P $DSN_USER"
    hint "sudo -u postgres createdb -O $DSN_USER $DSN_DB"
  fi
}

# ------------------------------------------------------- agent identities ----

setup_agent() {
  head_ "Agent identities"

  if ! curl -sf -m 5 "$HUMAN_BASE/version" >/dev/null 2>&1; then
    bad server "not answering on $HUMAN_BASE — start it first:"
    hint "cd server && make watch"
    hint "cd client/web/unify && pnpm dev"
    exit 1
  fi
  say server "up on $HUMAN_BASE"

  cat <<EOF

  Sign up your own login BEFORE this runs, if you have not: the first
  non-system user in the database is auto-promoted to super-admin, and
  agent@local.dev is about to take that slot. For a login created later:

    cd server && ./build/dev-bin --env-file .env roles useradd super-admin you@example.tld

EOF

  "$REPO_DIR/dev/agent/bootstrap.sh"

  if [[ -f "$UI_PASSWORD" ]]; then
    local pass
    pass="$(cat "$UI_PASSWORD")"
    if [[ -f "$E2E_ENV" ]]; then
      sed -i -E "s#^E2E_PASS=.*#E2E_PASS=$pass#" "$E2E_ENV"
      grep -q '^E2E_PASS=' "$E2E_ENV" || echo "E2E_PASS=$pass" >>"$E2E_ENV"
      say .env.e2e "E2E_PASS filled from .state/ui-password"
    fi
  else
    bad .env.e2e "bootstrap left no .state/ui-password — E2E_PASS not filled"
  fi
}

# ------------------------------------------------------------------ main -----

case "$MODE" in
  agent)
    setup_agent
    cat <<EOF

Ready. Launch Claude Code with 'make claude' — it mints the token the
human-local MCP server needs, which .mcp.json cannot do for itself:

  make claude
  make claude-yolo    # --dangerously-skip-permissions
EOF
    ;;

  check)
    check_tools
    install_deps
    write_server_env
    write_config_js
    write_e2e_env
    check_database
    head_ "$([[ "$PROBLEMS" == 0 ]] && echo 'Nothing to fix.' || echo "$PROBLEMS problem(s) above.")"
    [[ "$PROBLEMS" == 0 ]]
    ;;

  *)
    check_tools
    install_deps
    write_server_env
    write_config_js
    write_e2e_env
    check_database
    cat <<EOF

Setup done$([[ "$PROBLEMS" == 0 ]] || echo " — $PROBLEMS thing(s) marked ! above need you"). Next:

  1. cd server && make watch                 API on $HUMAN_BASE
  2. cd client/web/unify && pnpm dev         webapp on $HUMAN_WEBAPP
  3. open $HUMAN_WEBAPP and sign up — the FIRST user becomes super-admin
  4. make setup-agent                        the agent toolkit's own identities

The server rebuilds and restarts itself on any .go write — give it ~15s.
EOF
    ;;
esac
