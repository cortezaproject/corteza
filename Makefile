.PHONY: build

SPEC=$(GOPATH)/bin/spec

build:
	docker build --rm -t $(shell cat .project) .

$(SPEC):
	go get github.com/titpetric/spec/cmd/spec

spec: $(SPEC)
	cd sam/docs/src && $(SPEC)
	cd sam/ && php _gen.php
