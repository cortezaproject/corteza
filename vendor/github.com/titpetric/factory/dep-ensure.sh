#!/bin/bash
set -e
if [ -d "vendor" ]; then
	docker run --rm -v $(pwd):/go/src/go.casinogrounds.com/youtube -w /go/src/go.casinogrounds.com/youtube titpetric/golang dep ensure -update -v
else
	docker run --rm -v $(pwd):/go/src/go.casinogrounds.com/youtube -w /go/src/go.casinogrounds.com/youtube titpetric/golang dep init -v
fi
