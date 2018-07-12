.PHONY: build

SPEC=$(GOPATH)/bin/spec
PROTOC=$(GOPATH)/bin/protoc-gen-go

build:
	docker build --rm -t $(shell cat .project) .

$(SPEC):
	go get -u github.com/titpetric/spec/cmd/spec

$(PROTOC):
	go get -u github.com/golang/protobuf/protoc-gen-go

spec: $(SPEC)
	cd sam/docs/src && $(SPEC)
	cd sam/ && ./_gen.sh

protobuf: $(PROTOC)
	# @todo this needs work (it hangs and outputs nothing)
	$(PROTOC) --go_out=plugins=grpc:. -I. sam/chat/*.proto

dep:
	dep ensure -v

dep.update:
	dep ensure -update -v
