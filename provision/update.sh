#!/bin/bash

set -eu

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

BRANCH=${BRANCH:-"2021.3.x"}
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
  echo -e "\033[32mCopying ${3} ${2}\033[39m ... "
  mkdir -p "./${1}_${2}"
  cp -f "${DIR}/${2}/config/${3}.yaml" "./${1}_${2}/${3}.yaml"
}

function cleanup {
  echo -ne "\033[32mCleaning up\033[39m ... "
  rm -rf "${ZIP}" "${DIR}"
  echo "done"
}

download

copyExtConfig 700 crm 1000_namespace
copyExtConfig 700 crm 1100_modules
copyExtConfig 700 crm 1200_charts
copyExtConfig 700 crm 1400_pages
copyExtConfig 700 crm 1500_record_settings

copyExtConfig 701 service-solution 1000_namespace
copyExtConfig 701 service-solution 1100_modules
copyExtConfig 701 service-solution 1200_charts
copyExtConfig 701 service-solution 1400_pages
copyExtConfig 701 service-solution 1500_record_settings

cleanup

