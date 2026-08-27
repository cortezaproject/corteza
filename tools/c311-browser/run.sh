#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
compose_port="${C311_COMPOSE_PORT:-18081}"
admin_port="${C311_ADMIN_PORT:-18082}"
artifact_dir="${C311_ARTIFACT_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/c311-fe01.XXXXXX")}"

mkdir -p "$artifact_dir"
cleanup () {
  kill "${compose_pid:-}" "${admin_pid:-}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

corepack yarn --cwd "$repo_root/client/web/compose" serve --port "$compose_port" >"$artifact_dir/compose.log" 2>&1 & compose_pid=$!
corepack yarn --cwd "$repo_root/client/web/admin" serve --port "$admin_port" >"$artifact_dir/admin.log" 2>&1 & admin_pid=$!

for port in "$compose_port" "$admin_port"; do
  deadline=$((SECONDS + 180))
  until curl --fail --silent --max-time 2 "http://127.0.0.1:$port" >/dev/null; do
    if (( SECONDS >= deadline )); then
      echo "server did not start on $port" >&2
      exit 1
    fi
    sleep 1
  done
done

C311_COMPOSE_URL="http://127.0.0.1:$compose_port" \
C311_ADMIN_URL="http://127.0.0.1:$admin_port" \
C311_ARTIFACT_DIR="$artifact_dir" \
python3 "$repo_root/tools/c311-browser/fe01_matrix.py"
