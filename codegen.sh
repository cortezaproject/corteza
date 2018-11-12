#!/bin/bash
set -e

function gofmt {
	echo "=== fmt all folders ==="
	GOPATHS=$(find . -name '*.go' | grep -v vendor/ | xargs -n1 dirname | sort | uniq)
	for FOLDER in $GOPATHS; do
		#echo "== go fmt $FOLDER =="
		cd $FOLDER
		go fmt | xargs -n1 -I {} echo "$FOLDER/{}"
		cd $_PWD
	done
}

echo "=== go generate ==="
cd system/db && go generate && cd ../..
cd crm/db && go generate && cd ../..
cd crm/repository && go generate && cd ../..
cd sam/db && go generate && cd ../..

_PWD=$PWD
SPECS=$(find $PWD -name 'spec.json' | xargs -n1 dirname)
for SPEC in $SPECS; do
	echo "=== spec $SPEC ==="
	cd $SPEC && rm -rf spec && /usr/bin/env spec && cd $_PWD

	SRC=$(dirname $(dirname $SPEC))
	if [ -d "codegen/$(basename $SRC)" ]; then
		echo "=== codegen $SRC ==="
		codegen/codegen.php $(basename $SRC)
	fi
done

echo "=== codegen permissions ==="

go run sam/types/permissions/main.go -package types -function "func (c *Organisation) Permissions() []rbac.OperationGroup" -input sam/types/permissions/1-organisation.json -output sam/types/organisation_perms.go
go run sam/types/permissions/main.go -package types -function "func (c *Team) Permissions() []rbac.OperationGroup" -input sam/types/permissions/2-team.json -output sam/types/team_perms.go
go run sam/types/permissions/main.go -package types -function "func (c *Channel) Permissions() []rbac.OperationGroup" -input sam/types/permissions/3-channel.json -output sam/types/channel_perms.go

gofmt
