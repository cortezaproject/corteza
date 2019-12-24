#!/bin/bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

function download {
  SRC="https://raw.githubusercontent.com/cortezaproject/corteza-configs/${1}"
  DST=${2}
  echo -ne "\033[32m${DST}\033[39m (${SRC}) ..."
  curl -s $SRC > ${DST}
  echo "done"
}

# ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------

function getCrmConfig {
  NAMES="1000_namespace 1100_modules 1200_charts 1400_pages 1500_record_settings"

  for NAME in $NAMES; do
    download "master/crm/${NAME}.yaml" "./compose/src/${NAME}_crm.yaml"
  done
}

getCrmConfig

# ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------ ------

function getServiceCloudConfig {
  NAMES="1000_namespace 1100_modules 1200_charts 1400_pages 1500_record_settings"

  for NAME in $NAMES; do
    download "master/service-cloud/${NAME}.yaml" "./compose/src/${NAME}_service_cloud.yaml"
  done
}

getServiceCloudConfig

