#!/usr/bin/env bash
#
# Colour for the scripts a human reads.
#
# Colour is switched off whenever stdout is not a terminal, so a piped or
# redirected run stays greppable and a log file holds no escape sequences.
# NO_COLOR turns it off outright (no-color.org). FORCE_COLOR turns it back on
# for a pager or a CI that renders escape sequences, and outranks NO_COLOR —
# a standing preference loses to the flag typed on the command in front of it.
#
# Glyphs are NOT gated the same way: ✓ and ✗ are ordinary characters that a log
# file and a grep both survive, and gating them would change what the tests
# reading this output see.
#
# dev/setup.sh sources this before it has checked a single prerequisite, so
# this file may depend on nothing but bash.

if [[ -n "${FORCE_COLOR:-}" ]] || { [[ -t 1 ]] && [[ -z "${NO_COLOR:-}" ]]; }; then
  C_RESET=$'\033[0m'
  C_BOLD=$'\033[1m'
  C_DIM=$'\033[2m'
  C_RED=$'\033[31m'
  C_GREEN=$'\033[32m'
  C_YELLOW=$'\033[33m'
  C_CYAN=$'\033[36m'
else
  C_RESET=''
  C_BOLD=''
  C_DIM=''
  C_RED=''
  C_GREEN=''
  C_YELLOW=''
  C_CYAN=''
fi

G_OK='✓'
G_BAD='✗'
G_HINT='↳'
G_SECTION='▸'

# paint COLOUR TEXT — TEXT in COLOUR, and TEXT alone where colour is off.
paint() { printf '%s%s%s' "$1" "$2" "$C_RESET"; }

# One vocabulary for every script that reports to a human, so a result reads
# the same whichever one produced it. bad and warn go to stderr, which is where
# the messages they replaced already went — a caller redirecting stdout still
# sees what went wrong.
#
# A script needing a counter or an exit code of its own defines its own; these
# report and nothing more.
ok() { echo "  $(paint "$C_GREEN" "$G_OK") $*"; }
bad() { echo "  $(paint "$C_RED" "$G_BAD") $(paint "$C_RED" "$*")" >&2; }
warn() { echo "  $(paint "$C_YELLOW" '!') $(paint "$C_YELLOW" "$*")" >&2; }
note() { echo "    $(paint "$C_DIM" "$G_HINT $*")"; }
section() { printf '\n%s\n' "$(paint "$C_BOLD" "$G_SECTION $*")"; }

# step LABEL DETAIL — one thing a command did, in an aligned column.
step() { printf '  %s %-12s %s\n' "$(paint "$C_GREEN" "$G_OK")" "$1" "$2"; }

# Colour for the grep and python that some scripts render through. Exported so
# a child process can reach them; empty when colour is off, so the child needs
# no switch of its own.
export C_RESET C_BOLD C_DIM C_RED C_GREEN C_YELLOW C_CYAN

# Spelled as an if rather than `[[ ... ]] && VAR=x`: that form returns 1 when
# the test is false, and as the last line of a sourced file it makes `source`
# itself fail — which under `set -e` kills the caller with no message at all.
if [[ -n "$C_RESET" ]]; then
  GREP_COLOUR_WHEN=always
else
  GREP_COLOUR_WHEN=never
fi
export GREP_COLOUR_WHEN
