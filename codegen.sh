#!/bin/bash
set -e

_PWD=$PWD
SPECS=$(find $PWD -name 'spec.json' | xargs -n1 dirname)
for SPEC in $SPECS; do
	echo "=== spec $SPEC ==="
	cd $SPEC && rm -rf spec && /usr/bin/env spec && cd .. && ./README.php && cd $_PWD

	SRC=$(dirname $(dirname $SPEC))
	echo "=== codegen $SRC ==="
	GOPATHS=$(codegen/codegen.php $(basename $SRC) | tee -a /dev/stderr | xargs -n1 dirname | sort | uniq)
	for FOLDER in $GOPATHS; do
		if [[ $FOLDER != "." ]]; then
			echo "== go fmt $FOLDER =="
			cd $FOLDER
			go fmt
			cd $_PWD
		fi
	done
done

