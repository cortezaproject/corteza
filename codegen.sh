#!/bin/bash
set -e
_PWD=$PWD
SPECS=$(find -name 'spec.json' | xargs -n1 dirname)
for SPEC in $SPECS; do
	echo "=== spec $SPEC ==="
	cd $SPEC && spec && cd $_PWD

	SRC=$(dirname $(dirname $SPEC))
	echo "=== codegen $SRC ==="
	cd $SRC && ../codegen/codegen.php && go fmt && cd $_PWD
done