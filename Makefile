.PHONY: build realize dep spec protobuf qa qa.test qa.vet codegen

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

########################################################################################################################
# Tool bins
SPEC      = $(GOPATH)/bin/spec
PROTOC    = $(GOPATH)/bin/protoc-gen-go
REALIZE   = ${GOPATH}/bin/realize


build:
	docker build --rm -t $(shell cat .project) .


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.update:
	dep ensure -update -v

dep:
	dep ensure -v

codegen: $(SPEC)
	./codegen.sh

protobuf: $(PROTOC)
	# @todo this needs work (it hangs and outputs nothing)
	$(PROTOC) --go_out=plugins=grpc:. -I. sam/chat/*.proto


########################################################################################################################
# QA

qa: qa.vet qa.test

qa.vet:
	go vet `cd ${GOPATH}/src/; find $(PKG) -type f -name '*.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`

qa.test:
	go test `cd ${GOPATH}/src/; find $(PKG) -type f -name '*_test.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`



########################################################################################################################
# Toolset

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(SPEC):
	$(GOGET) github.com/titpetric/spec/cmd/spec

$(PROTOC):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

