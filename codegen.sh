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

function permissions {
	yellow "> permissions"
	if [ ! -f "build/gen-permissions" ]; then
		CGO_ENABLED=0 go build -o ./build/gen-permissions codegen/v2/permissions.go 
	fi
	
	./build/gen-permissions -package types -function "func (c *Organisation) Permissions() []rules.OperationGroup" -input sam/types/permissions/1-organisation.json -output sam/types/organisation.perms.go
	./build/gen-permissions -package types -function "func (c *Team) Permissions() []rules.OperationGroup" -input sam/types/permissions/2-team.json -output sam/types/team.perms.go
	./build/gen-permissions -package types -function "func (c *Channel) Permissions() []rules.OperationGroup" -input sam/types/permissions/3-channel.json -output sam/types/channel.perms.go

	green "OK"
}

permissions

function types {
	yellow "> types"
	if [ ! -f "build/gen-type-set" ]; then
		CGO_ENABLED=0 go build -o ./build/gen-type-set codegen/v2/type-set.go 
	fi

	./build/gen-type-set --types Module,Page,Chart,Trigger,Record \
	                      --no-pk-types ModuleField,RecordValue \
	                      --output crm/types/type.gen.go

	./build/gen-type-set --types MessageAttachment --output sam/types/attachment.gen.go
	./build/gen-type-set --types Channel --output sam/types/channel.gen.go
	./build/gen-type-set --no-pk-types ChannelMember --output sam/types/channel_member.gen.go
	./build/gen-type-set --no-pk-types Command,CommandParam --output sam/types/command.gen.go
	./build/gen-type-set --types Mention --output sam/types/mention.gen.go
	./build/gen-type-set --types MessageFlag --output sam/types/message_flag.gen.go
	./build/gen-type-set --types Message --output sam/types/message.gen.go
	./build/gen-type-set --no-pk-types Unread --output sam/types/unread.gen.go

	./build/gen-type-set --types User --output system/types/user.gen.go
	./build/gen-type-set --types Credentials --output system/types/credentials.gen.go
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

_PWD=$PWD
SPECS=$(find $PWD -name 'spec.json' | xargs -n1 dirname)
for SPEC in $SPECS; do
	yellow "> spec $SPEC"
	cd $SPEC && rm -rf spec && /usr/bin/env spec && cd $_PWD
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

gofmt
