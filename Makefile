.PHONY: build realize dep spec protobuf critic test test.rbac qa qa.test qa.vet codegen

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

########################################################################################################################
# Tool bins
SPEC      = $(GOPATH)/bin/spec
PROTOC    = $(GOPATH)/bin/protoc-gen-go
REALIZE   = ${GOPATH}/bin/realize
GOTEST    = ${GOPATH}/bin/gotest
GOCRITIC  = ${GOPATH}/bin/gocritic

build:
	docker build --no-cache --rm -t $(shell cat .project) .


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

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

test: $(GOTEST)
	$(GOTEST) -cover -v ./...

test.rbac: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./rbac
	go tool cover -func=.cover.out 

qa: qa.vet qa.test

qa.vet:
	go vet `cd ${GOPATH}/src/; find $(PKG) -type f -name '*.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`

qa.test:
	go test `cd ${GOPATH}/src/; find $(PKG) -type f -name '*_test.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`



########################################################################################################################
# Toolset

$(GOTEST):
	$(GOGET) github.com/rakyll/gotest

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(SPEC):
	$(GOGET) github.com/titpetric/spec/cmd/spec

$(PROTOC):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

clean.tools:
	rm -f $(SPEC) $(PROTOC) $(REALIZE)
