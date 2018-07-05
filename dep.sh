#!/bin/bash
set -e
PROJECT=$(basename $(dirname $(readlink -f $0)))

if [ -d "vendor" ]; then
	docker run --rm -v $(pwd):/go/src/github.com/crusttech/$PROJECT -w /go/src/github.com/crusttech/$PROJECT -e GOOS=${OS} -e GOARCH=${ARCH} -e CGO_ENABLED=0 -e GOARM=7 titpetric/golang dep ensure -update -v
else
	docker run --rm -v $(pwd):/go/src/github.com/crusttech/$PROJECT -w /go/src/github.com/crusttech/$PROJECT -e GOOS=${OS} -e GOARCH=${ARCH} -e CGO_ENABLED=0 -e GOARM=7 titpetric/golang dep init
fi
