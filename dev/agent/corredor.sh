#!/usr/bin/env bash
#
# Run Corredor — the automation script runner — for THIS checkout.
#
# The server dials Corredor over gRPC and refuses to boot when it is enabled
# and unreachable, so Corredor comes up first and goes down last.
#
#   corredor.sh up       start it, serving dev/fixtures/corredor
#   corredor.sh down     stop it
#   corredor.sh status   say whether it is listening, and on what
#   corredor.sh env      print the two lines server/.env needs
#
# A worktree gets this from `worktree.sh up`; the primary is where this script
# earns its place. The port is the checkout's API port plus 50000, so two
# checkouts never share one.

set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

RUN_DIR="$REPO_DIR/.run"
PID_FILE="$RUN_DIR/corredor.pid"
LOG_FILE="$RUN_DIR/corredor.log"
BUNDLE_DIR="$RUN_DIR/corredor-bundles"
EXT_PATHS="$REPO_DIR/dev/fixtures/corredor:$REPO_DIR/corredor/usr:$REPO_DIR/corredor/usr/*"

api_port() { sed -nE 's#^https?://[^:]+:([0-9]+).*#\1#p' <<<"$HUMAN_API"; }
grpc_port() { echo $((50000 + $(api_port))); }

listening() { ss -ltn "sport = :$1" 2>/dev/null | grep -q LISTEN; }

cmd_env() {
  cat <<EOF
CORREDOR_ENABLED=true
CORREDOR_ADDR=localhost:$(grpc_port)
EOF
}

cmd_up() {
  local port
  port="$(grpc_port)"

  [[ -f "$REPO_DIR/corredor/package.json" ]] ||
    die "this checkout has no corredor/ — nothing to run"

  if listening "$port"; then
    ok "corredor already listening on :$port"
    return 0
  fi

  [[ -d "$REPO_DIR/corredor/node_modules" ]] || {
    note "installing corredor dependencies (first run) …"
    (cd "$REPO_DIR" && pnpm install --silent)
  }

  mkdir -p "$RUN_DIR" "$BUNDLE_DIR"

  # setsid so the whole tree stops with one signal, and so it outlives the
  # shell that started it. The subshell gets its own descriptors: it inherits
  # the caller's stdout otherwise, and a caller reading this script through a
  # pipe waits on that copy long after corredor is up.
  (
    cd "$REPO_DIR/corredor" && setsid bash -c 'echo $$ >"$1"; shift; exec "$@"' _ "$PID_FILE" \
      env CORREDOR_ENVIRONMENT=dev \
      CORREDOR_ADDR="localhost:$port" \
      CORREDOR_SERVER_CERTIFICATES_ENABLED=false \
      CORREDOR_LOG_PRETTY=true \
      CORREDOR_EXT_SEARCH_PATHS="$EXT_PATHS" \
      CORREDOR_EXT_DEPENDENCIES_AUTO_UPDATE=false \
      CORREDOR_BUNDLER_OUTPUT_PATH="$BUNDLE_DIR" \
      CORREDOR_EXEC_CSERVERS_API_HOST="localhost:$(api_port)" \
      CORREDOR_EXEC_CSERVERS_API_BASEURL_TEMPLATE='http://{host}/api/{service}' \
      CORREDOR_EXEC_CTX_FRONTEND_BASEURL="$HUMAN_WEBAPP" \
      pnpm serve >"$LOG_FILE" 2>&1 &
  ) >/dev/null 2>&1 </dev/null

  for _ in $(seq 1 100); do
    listening "$port" && break
    sleep 0.3
  done

  if listening "$port"; then
    step corredor "grpc :$port   $(paint "$C_DIM" "(log $LOG_FILE)")"
    note "the server needs these in server/.env, then a restart:"
    cmd_env | sed 's/^/      /'
  else
    bad "corredor did not come up on :$port — see $LOG_FILE"
    return 1
  fi
}

cmd_down() {
  local port pid
  port="$(grpc_port)"

  if [[ -f "$PID_FILE" ]]; then
    pid="$(cat "$PID_FILE")"
    kill -TERM -"$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null || true
    rm -f "$PID_FILE"
  fi

  for _ in $(seq 1 40); do
    listening "$port" || break
    sleep 0.3
  done

  if listening "$port"; then
    bad "port $port is still bound"
    return 1
  fi

  ok "corredor stopped"
}

cmd_status() {
  local port
  port="$(grpc_port)"

  if listening "$port"; then
    ok "corredor listening on :$port $(paint "$C_DIM" "(log $LOG_FILE)")"
  else
    echo "  $(paint "$C_DIM" "corredor is not running (would listen on :$port)")"
  fi

  if grep -qE '^CORREDOR_ENABLED=true' "$SERVER_DIR/.env" 2>/dev/null; then
    ok "server/.env has CORREDOR_ENABLED=true"
  else
    warn "server/.env does not enable corredor — 'corredor.sh env' prints what it needs"
  fi
}

case "${1:-}" in
  up) cmd_up ;;
  down) cmd_down ;;
  status) cmd_status ;;
  env) cmd_env ;;
  *)
    sed -n '3,16p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
    exit 1
    ;;
esac
