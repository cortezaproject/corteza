#!/bin/bash
set -e

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

	./build/gen-type-set --types Attachment --output crm/types/attachment.gen.go
	./build/gen-type-set --types Module     --output crm/types/module.gen.go
	./build/gen-type-set --types Page       --output crm/types/page.gen.go
	./build/gen-type-set --types Chart      --output crm/types/chart.gen.go
	./build/gen-type-set --types Trigger    --output crm/types/trigger.gen.go
	./build/gen-type-set --types Record     --output crm/types/record.gen.go

	./build/gen-type-set --with-primary-key=false --types ModuleField --output crm/types/module_field.gen.go
	./build/gen-type-set --with-primary-key=false --types RecordValue --output crm/types/record_value.gen.go

	./build/gen-type-set --types MessageAttachment --output messaging/types/attachment.gen.go
	./build/gen-type-set --types Mention           --output messaging/types/mention.gen.go
	./build/gen-type-set --types MessageFlag       --output messaging/types/message_flag.gen.go
	./build/gen-type-set --types Message           --output messaging/types/message.gen.go
	./build/gen-type-set --types Channel           --output messaging/types/channel.gen.go
	./build/gen-type-set --types Webhook           --output messaging/types/webhook.gen.go

	./build/gen-type-set --with-primary-key=false --types ChannelMember --output messaging/types/channel_member.gen.go
	./build/gen-type-set --with-primary-key=false --types Command       --output messaging/types/command.gen.go
	./build/gen-type-set --with-primary-key=false --types CommandParam  --output messaging/types/command_param.gen.go
	./build/gen-type-set --with-primary-key=false --types Unread        --output messaging/types/unread.gen.go

	./build/gen-type-set --types User         --output system/types/user.gen.go
	./build/gen-type-set --types Application  --output system/types/application.gen.go
	./build/gen-type-set --types Role         --output system/types/role.gen.go
	./build/gen-type-set --types Organisation --output system/types/organisation.gen.go
	./build/gen-type-set --types Credentials  --output system/types/credentials.gen.go

	green "OK"
}

types

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

database

function files {
	yellow "> files"
	FOLDERS=$(find . -type d -wholename '*/data')
	for FOLDER in $FOLDERS; do
		FOLDER=$(dirname $FOLDER)
		FOLDER=${FOLDER:2}
		echo $FOLDER
		cd $FOLDER && $GOPATH/bin/statik -p files -m -Z -f -src=data && cd $_PWD
	done
	green "OK"
}

files

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
			green "OK"
		fi
	done
}

specs

gofmt
