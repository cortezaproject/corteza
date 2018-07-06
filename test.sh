#!/bin/bash

set -e
set -u

PKG=github.com/crusttech/crust

go test `cd ${GOPATH}/src/; find ${PKG} -type f -name '*_test.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`