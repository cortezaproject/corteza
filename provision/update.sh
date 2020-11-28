#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BRANCH=${BRANCH:-"develop"}
ZIP="${BRANCH}.zip"
URL=${URL:-"https://github.com/cortezaproject/corteza-ext/archive/${ZIP}"}
DIR="corteza-ext-${BRANCH}"

function download {
  echo -ne "\033[32mDownloading ${URL}\033[39m ... "

  curl -s --location "${URL}" > "${ZIP}"
  unzip -qq "${ZIP}"

  echo "done"
}

function copyExtConfig {
  echo -e "\033[32mCopying ${2} ${1}\033[39m ... "
  mkdir -p "./compose/src/${1}"
  cp "${DIR}/${1}/config/${2}.yaml" "./compose/src/${1}/${2}.yaml"
}

function cleanup {
  echo -ne "\033[32mCleaning up\033[39m ... "
  rm -rf "${ZIP}" "${DIR}"
  echo "done"
}

download

copyExtConfig crm 1000_namespace
copyExtConfig crm 1100_modules
copyExtConfig crm 1200_charts
copyExtConfig crm 1400_pages
copyExtConfig crm 1500_record_settings

copyExtConfig service-solution 1000_namespace
copyExtConfig service-solution 1100_modules
copyExtConfig service-solution 1200_charts
copyExtConfig service-solution 1400_pages
copyExtConfig service-solution 1500_record_settings

cleanup

