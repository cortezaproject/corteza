#!/usr/bin/env bash
# Print the shape of a REST endpoint before calling it: method, path, params, traps.
#
# Usage: shape.sh [SERVICE[/GROUP [CALL]]] | --grep WORD
#   shape.sh                        services and their endpoint groups
#   shape.sh compose/record         every call of the record group
#   shape.sh compose/record update  method, path, parameters and known traps
#   shape.sh --grep undelete        every call whose name or path has the word
#
# Reads server/<service>/rest.yaml, the codegen input the handlers are made
# from, and traps.yaml next to this script. Needs python3 with PyYAML.
set -euo pipefail
source "$(dirname "${BASH_SOURCE[0]}")/common.sh"
python3 -c 'import yaml' 2>/dev/null || {
  bad "python3 cannot import yaml (PyYAML); install it: pip install pyyaml / apt install python3-yaml"
  exit 1
}
exec python3 "$AGENT_DIR/shape.py" "$@"
