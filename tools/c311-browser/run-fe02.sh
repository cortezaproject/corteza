#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
compose_port="${C311_COMPOSE_PORT:-18081}"
artifact_dir="${C311_ARTIFACT_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/c311-fe02.XXXXXX")}"

config_file="$repo_root/client/web/compose/public/config.js"
created_config=0
if [[ -L "$config_file" ]]; then
  echo "config.js must not be a symbolic link" >&2
  exit 1
fi
if [[ ! -f "$config_file" ]]; then
  cp "$repo_root/client/web/compose/public/config.example.js" "$config_file"
  created_config=1
fi

if [[ -L "$artifact_dir" ]]; then
  echo "artifact directory must not be a symbolic link" >&2
  exit 1
fi
mkdir -p "$artifact_dir"
chmod 700 "$artifact_dir"
cleanup () {
  kill "${compose_pid:-}" 2>/dev/null || true
  if (( created_config )); then rm -f "$config_file"; fi
}
trap cleanup EXIT INT TERM

corepack yarn --cwd "$repo_root/client/web/compose" serve --port "$compose_port" >"$artifact_dir/compose-fe02.log" 2>&1 & compose_pid=$!
deadline=$((SECONDS + 180))
until curl --fail --silent --max-time 2 "http://127.0.0.1:$compose_port" >/dev/null; do
  if (( SECONDS >= deadline )); then echo "server did not start on $compose_port" >&2; exit 1; fi
  sleep 1
done

C311_COMPOSE_URL="http://127.0.0.1:$compose_port" \
C311_ARTIFACT_DIR="$artifact_dir" \
python3 "$repo_root/tools/c311-browser/fe02_matrix.py"
