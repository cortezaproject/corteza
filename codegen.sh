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

function gofmt {
	yellow "> go fmt ./..."
	go fmt ./...
	green "OK"
}

function types {
	yellow "> types"
	if [ ! -f "build/gen-type-set" ]; then
		CGO_ENABLED=0 go build -o ./build/gen-type-set codegen/v2/type-set.go
	fi
	if [ ! -f "build/gen-type-set-test" ]; then
		CGO_ENABLED=0 go build -o ./build/gen-type-set-test codegen/v2/type-set-test.go
	fi

	./build/gen-type-set --types Namespace   --output compose/types/namespace.gen.go
	./build/gen-type-set --types Attachment  --output compose/types/attachment.gen.go
	./build/gen-type-set --types Module      --output compose/types/module.gen.go
	./build/gen-type-set --types Page        --output compose/types/page.gen.go
	./build/gen-type-set --types Chart       --output compose/types/chart.gen.go
	./build/gen-type-set --types Record      --output compose/types/record.gen.go
	./build/gen-type-set --types ModuleField --output compose/types/module_field.gen.go

	./build/gen-type-set-test --types Namespace   --output compose/types/namespace.gen_test.go
	./build/gen-type-set-test --types Attachment  --output compose/types/attachment.gen_test.go
	./build/gen-type-set-test --types Module      --output compose/types/module.gen_test.go
	./build/gen-type-set-test --types Page        --output compose/types/page.gen_test.go
	./build/gen-type-set-test --types Chart       --output compose/types/chart.gen_test.go
	./build/gen-type-set-test --types Record      --output compose/types/record.gen_test.go
	./build/gen-type-set-test --types ModuleField --output compose/types/module_field.gen_test.go

	./build/gen-type-set --with-primary-key=false --types RecordValue --output compose/types/record_value.gen.go
	./build/gen-type-set-test --with-primary-key=false --types RecordValue --output compose/types/record_value.gen_test.go

	./build/gen-type-set --types MessageAttachment --output messaging/types/attachment.gen.go
	./build/gen-type-set --types Mention           --output messaging/types/mention.gen.go
	./build/gen-type-set --types MessageFlag       --output messaging/types/message_flag.gen.go
	./build/gen-type-set --types Message           --output messaging/types/message.gen.go
	./build/gen-type-set --types Channel           --output messaging/types/channel.gen.go
	./build/gen-type-set --types Webhook           --output messaging/types/webhook.gen.go

	./build/gen-type-set-test --types MessageAttachment --output messaging/types/attachment.gen_test.go
	./build/gen-type-set-test --types Mention           --output messaging/types/mention.gen_test.go
	./build/gen-type-set-test --types MessageFlag       --output messaging/types/message_flag.gen_test.go
	./build/gen-type-set-test --types Message           --output messaging/types/message.gen_test.go
	./build/gen-type-set-test --types Channel           --output messaging/types/channel.gen_test.go
	./build/gen-type-set-test --types Webhook           --output messaging/types/webhook.gen_test.go

	./build/gen-type-set --with-primary-key=false --types ChannelMember --output messaging/types/channel_member.gen.go
	./build/gen-type-set --with-primary-key=false --types Command       --output messaging/types/command.gen.go
	./build/gen-type-set --with-primary-key=false --types CommandParam  --output messaging/types/command_param.gen.go
	./build/gen-type-set --with-primary-key=false --types Unread        --output messaging/types/unread.gen.go

	./build/gen-type-set-test --with-primary-key=false --types ChannelMember --output messaging/types/channel_member.gen_test.go
	./build/gen-type-set-test --with-primary-key=false --types Command       --output messaging/types/command.gen_test.go
	./build/gen-type-set-test --with-primary-key=false --types CommandParam  --output messaging/types/command_param.gen_test.go
	./build/gen-type-set-test --with-primary-key=false --types Unread        --output messaging/types/unread.gen_test.go

	./build/gen-type-set --types User         --output system/types/user.gen.go
	./build/gen-type-set --types Application  --output system/types/application.gen.go
	./build/gen-type-set --types Role         --output system/types/role.gen.go
	./build/gen-type-set --types Organisation --output system/types/organisation.gen.go
	./build/gen-type-set --types Credentials  --output system/types/credentials.gen.go
	./build/gen-type-set --types Reminder     --output system/types/reminder.gen.go
	./build/gen-type-set --types Attachment   --output system/types/attachment.gen.go

	./build/gen-type-set-test --types User         --output system/types/user.gen_test.go
	./build/gen-type-set-test --types Application  --output system/types/application.gen_test.go
	./build/gen-type-set-test --types Role         --output system/types/role.gen_test.go
	./build/gen-type-set-test --types Organisation --output system/types/organisation.gen_test.go
	./build/gen-type-set-test --types Credentials  --output system/types/credentials.gen_test.go
	./build/gen-type-set-test --types Reminder     --output system/types/reminder.gen_test.go
	./build/gen-type-set-test --types Attachment   --output system/types/attachment.gen_test.go

	./build/gen-type-set --types Value --output pkg/settings/types.gen.go --with-primary-key=false --package settings
	./build/gen-type-set-test --types Value --output pkg/settings/types.gen_test.go --with-primary-key=false --package settings

	./build/gen-type-set --types Rule      --output pkg/permissions/rule.gen.go     --with-primary-key=false --package permissions
	./build/gen-type-set --types Resource  --output pkg/permissions/resource.gen.go --with-primary-key=false --package permissions

	./build/gen-type-set-test --types Rule      --output pkg/permissions/rule.gen_test.go     --with-primary-key=false --package permissions
	./build/gen-type-set-test --types Resource  --output pkg/permissions/resource.gen_test.go --with-primary-key=false --package permissions

	./build/gen-type-set --types Script --output pkg/corredor/types.gen.go --with-primary-key=false --package corredor
	./build/gen-type-set-test --types Script --output pkg/corredor/types.gen_test.go --with-primary-key=false --package corredor


	green "OK"
}


function database {
	yellow "> database"
	FOLDERS=$(find . -type d -wholename '*/schema/mysql')
	for FOLDER in $FOLDERS; do
		FOLDER=$(dirname $(dirname $FOLDER))
		FOLDER=${FOLDER:2}
		echo $FOLDER
		cd $FOLDER && $GOPATH/bin/statik -p mysql -m -Z -f -src=schema/mysql && cd $_PWD
	done
	green "OK"
}


function provision {
	yellow "> provision files"
	for FOLDER in system compose messaging; do
   	$GOPATH/bin/statik -p $FOLDER -m -Z -f -src="./provision/$FOLDER/src" -dest "./provision"
	done
	green "OK"
}


function events {
  if [ ! -f "build/event-gen" ]; then
		CGO_ENABLED=0 go build -o ./build/event-gen ./codegen/v2/events
	fi

	for SERVICE in system compose messaging; do
	  yellow "> event files for ${SERVICE}"
	  ./build/event-gen --service ${SERVICE}
	done
	green "OK"
}


function specs {
	yellow "> specs"
	if [ ! -f "build/gen-spec" ]; then
		CGO_ENABLED=0 go build -o ./build/gen-spec codegen/v2/spec.go
	fi
	_PWD=$PWD
	SPECS=$(find $PWD -name 'spec.json' | xargs -n1 dirname)
	for SPEC in $SPECS; do
		yellow "> spec $SPEC"
		cd $SPEC && rm -rf spec && ../../build/gen-spec && cd $_PWD
		green "OK"
	done

	for SPEC in $SPECS; do
		SRC=$(basename $SPEC)
		if [ -d "codegen/$SRC" ]; then
			yellow "> README $SRC"
			codegen/codegen.php $SRC
			rsync -a codegen/common/ $SRC/
			green "OK"
		fi
	done
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
  types)
    types
    ;;
  database)
    database
    ;;
  provision)
    provision
    ;;
  specs)
    specs
    ;;
  proto)
    proto
    ;;
  events)
    events
    ;;
  all)
    types
    database
    provision
    specs
    proto
    events
esac

# Always finish with fmt
gofmt
