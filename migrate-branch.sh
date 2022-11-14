#!/bin/bash

# Feel free to remove this script when migration is complete

set -eu

function help {
  echo "This script migrates your branch, commit-by-commit from"
  echo "your local git repository to a monorepo repository."
  echo ""
  echo "In a nutshell it executes combination of format-patch and am"
  echo "commands with a couple of additional checks and cleanups."
  echo ""
  echo "./migrate-branch <legacy-repo-path> branch[, ...]"
  echo ""
}

LBLUE=$(tput setaf 051)
LRED=$(tput setaf 196)
LGREEN=$(tput setaf 046)
BLACK=$(tput setaf 000)
BG_GREEN=$(tput setab 2)
BOLD=$(tput bold)
CLR=$(tput sgr0)

function _header {
  echo -e "${BG_GREEN}${BLACK}> ${1}${CLR}"
}

function _info { echo -ne "${LBLUE}${1}${CLR}";   }
function _err  { echo -ne "${LRED}${1}${CLR}";  }

function _infoln { _info "${1}\n"; }
function _errln { _err "${1}\n"; }
function _fatal {
  _errln "${1}";
  exit 1;
}

LEGACYRP=${1:-""}

# show help if no legacy repo path is specified or if value ends with "help" or "h"
if [ -z "${LEGACYRP}" ] || [[ "${LEGACYRP}" =~ (help|h)$ ]]; then
  help
  exit 0
fi


shift
BRANCHES=${@:-""}

if [ -z "${LEGACYRP}" ]; then
  _fatal "Use: migrate-branch <legacy-repo-path> [branch, ...]"
fi

if [ ! -d "${LEGACYRP}" ]; then
  _fatal "Error: legacy repo path does not exist"
fi

if [ -z "${BRANCHES}" ]; then
  # fetch current branch from legacy repo
  BRANCHES=$(git -C "${LEGACYRP}" branch --quiet --all --list)

  # if no branches ar found, exit
  if [ -z "${BRANCHES}" ]; then
    _fatal "Error: no branches found in ${LEGACYRP}"
  fi

  _infoln "Found branches in ${LEGACYRP}:"
  echo "${BRANCHES}"
  exit 0
fi

# extract repository name from the remote
REPO=$(git -C "${LEGACYRP}" remote get-url origin | sed -E 's/.*\/(.*)\.git/\1/')

# guess new path from legacy repo
case $REPO in
  corteza-server)            NEW_PATH="server" ;;
  corteza-webapp)            NEW_PATH="client/web" ;;
  corteza-webapp-admin)      NEW_PATH="client/web/admin" ;;
  corteza-webapp-reporter)   NEW_PATH="client/web/reporter" ;;
  corteza-webapp-compose)    NEW_PATH="client/web/compose" ;;
  corteza-webapp-privacy)    NEW_PATH="client/web/privacy" ;;
  corteza-webapp-workflow)   NEW_PATH="client/web/workflow" ;;
  corteza-webapp-discovery)  NEW_PATH="client/web/discovery" ;;
  corteza-webapp-one)        NEW_PATH="client/web/one" ;;
  corteza-vue)               NEW_PATH="lib/vue" ;;
  corteza-js)                NEW_PATH="lib/js" ;;

  corteza-server-discovery) PATH="server-discovery" ;;

  *) _fatal "Unsupported legacy repo: ${REPO}" ;;
esac

# get current branch
BASE_BRANCH=$(git branch --quiet --show-current)

# If $BASE_BRANCH does not end with .x, ask for confirmation to proceed
if [[ ! "${BASE_BRANCH}" =~ \.x$ ]]; then
  _infoln "Current branch (on mono repo) is ${BOLD}${BASE_BRANCH}${LBLUE} and does not end with .x (e.g. a version branch 2021.9.x)"
  echo "If you know what you are doing proceed other abort and switch to a version branch"
  read -n 1 -p "Are you sure you want to proceed? [y/N] "; echo
  if [ "${REPLY}" != "y" ]; then
    _fatal "Aborting"
  fi
fi


# loop through branches
for BRANCH in ${BRANCHES}; do
  NEW_BRANCH="legacy-${REPO}-${BRANCH}"

  # Check if $BASE_BRANCH is a suffix of $BRANCH?
  read -n 1 -p "Migrating branch ${BOLD}${BRANCH}${CLR} as ${NEW_BRANCH} using ${BASE_BRANCH} as base.\n Proceed? [y/N] "; echo
  if [ "${REPLY}" != "y" ]; then
    echo "Skipping"
    continue
  fi

  # get list of all commits since HEAD of base-branch branch
  COMMITS=$(git -C "${LEGACYRP}" log --reverse --pretty=format:"%H" "${BASE_BRANCH}..${BRANCH}")

  _infoln "Found $(echo ${COMMITS} | wc -l) commits to migrate"

  # switch to base branch to ensure new one will be created from this
  git switch "${BASE_BRANCH}"

  # create feature branch from base branch
  git switch -C "${NEW_BRANCH}"

  # take all commits between base and feature branch
  # and apply them to the mono-repo in the appropriate path
  git -C "${LEGACYRP}" format-patch \
      --subject-prefix="" \
      --stdout \
      "${BASE_BRANCH}..${BRANCH}" \
  | git am \
    --ignore-whitespace \
    --no-scissors \
    --ignore-space-change \
    --directory="${NEW_PATH}"

  git switch "${BASE_BRANCH}"
done

