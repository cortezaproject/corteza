#!/bin/bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

function download {
  SRC="https://raw.githubusercontent.com/cortezaproject/corteza-configs/master/${1}"
  DST=${2}
  echo -ne "\033[32m${DST}\033[39m (${SRC}) ..."
  curl -s $SRC > ${DST}
  echo "done"
}

function getCrmConfig {
  NAMES="1000_namespace 1100_modules 1200_charts 1300_scripts 1400_pages 1500_record_settings"

  for NAME in $NAMES; do
    download "crm/${NAME}.yaml" "./compose/src/${NAME}.yaml"
  done
}

getCrmConfig

