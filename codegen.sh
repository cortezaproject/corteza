#!/bin/bash
set -e

if [ -f ./build ]; then
	find ./build -name gen-* -delete
fi

if [ -f ./.env ]; then
  source .env
fi;

_PWD=$PWD

function yellow {
	echo -e "\033[33m$@\033[39m"
}
function green {
	echo -e "\033[32m$@\033[39m"
}

function proto {
	yellow "> proto"

	# Where should we look for the files
	PROTOBUF_PATH="codegen/corteza-protobuf"
	CORTEZA_PROTOBUF_PATH=${CORTEZA_PROTOBUF_PATH:-"${PROTOBUF_PATH}"}

  # Download protobufs to the primary location
  BRANCH=${BRANCH:-"develop"}
  ZIP="${BRANCH}.zip"
  URL=${URL:-"https://github.com/cortezaproject/corteza-protobuf/archive/${ZIP}"}
  rm -rf "${PROTOBUF_PATH}"
  curl -s --location "${URL}" > "codegen/${ZIP}"
  unzip -qq -o -d "codegen/" "codegen/${ZIP}"
  mv -f "codegen/corteza-protobuf-${BRANCH}" "${PROTOBUF_PATH}"

  DIR=./pkg/corredor
  mkdir -p ${DIR}
	yellow "  ${CORTEZA_PROTOBUF_PATH} >> ${DIR}"
	PATH=$PATH:$GOPATH/bin protoc \
		--proto_path ${CORTEZA_PROTOBUF_PATH} \
		--go_out="plugins=grpc:./${DIR}" \
		service-corredor.proto

	yellow "  ${CORTEZA_PROTOBUF_PATH} >> system/proto"
	PATH=$PATH:$GOPATH/bin protoc \
		--proto_path ${CORTEZA_PROTOBUF_PATH}/system \
		--go_out=plugins=grpc:system/proto \
		user.proto role.proto
  green "OK"
}

case ${1:-"all"} in
  provision)
    provision
    ;;
  proto)
    proto
    ;;
  all)
    provision
    proto
esac
