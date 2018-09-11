#!/bin/bash
set -e

function gofmt {
	echo "=== fmt all folders ==="
	GOPATHS=$(find -name '*.go' | grep -v vendor/ | xargs -n1 dirname | sort | uniq)
	for FOLDER in $GOPATHS; do
		#echo "== go fmt $FOLDER =="
		cd $FOLDER
		go fmt | xargs -n1 -I {} echo "$FOLDER/{}"
		cd $_PWD
	done
}

_PWD=$PWD
SPECS=$(find $PWD -name 'spec.json' | xargs -n1 dirname)
for SPEC in $SPECS; do
	echo "=== spec $SPEC ==="
	if [ -x "$(dirname $SPEC)/README.php" ]; then
		cd $SPEC && rm -rf spec && /usr/bin/env spec && cd .. && ./README.php && cd $_PWD
	fi

	SRC=$(dirname $(dirname $SPEC))
	echo "=== codegen $SRC ==="

	codegen/codegen.php $(basename $SRC) | tee -a /dev/stderr
done

gofmt