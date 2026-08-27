#!/usr/bin/env bash
#
# Colour for the scripts a human reads.
#
# Colour is switched off whenever stdout is not a terminal, so a piped or
# redirected run stays greppable and a log file holds no escape sequences.
# NO_COLOR turns it off outright (no-color.org) and FORCE_COLOR turns it back
# on, for a pager or a CI that renders them.
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
