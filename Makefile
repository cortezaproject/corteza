.PHONY: build realize dep spec protobuf critic test test.rbac qa qa.test qa.vet codegen

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

BASEPKGS = rbac auth crm sam

########################################################################################################################
# Tool bins
DEP      = $(GOPATH)/bin/dep
SPEC      = $(GOPATH)/bin/spec
PROTOC    = $(GOPATH)/bin/protoc-gen-go
REALIZE   = ${GOPATH}/bin/realize
GOTEST    = ${GOPATH}/bin/gotest
GOCRITIC  = ${GOPATH}/bin/gocritic
MOCKGEN   = ${GOPATH}/bin/mockgen

build:
	docker build --no-cache --rm -t $(shell cat .project) .


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.update: $(DEP)
	$(DEP) ensure -update -v

dep: $(DEP)
	$(DEP) ensure -v

codegen: $(SPEC)
	./codegen.sh

protobuf: $(PROTOC)
	# @todo this needs work (it hangs and outputs nothing)
	$(PROTOC) --go_out=plugins=grpc:. -I. sam/chat/*.proto


########################################################################################################################
# QA

critic: $(GOCRITIC)
	PATH=${PATH}:${GOPATH}/bin $(GOCRITIC) check-project .

test: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./...
	$(GO) tool cover -func=.cover.out

test.rbac: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./rbac/...
	$(GO) tool cover -func=.cover.out

vet:
	$(GO) vet `cd ${GOPATH}/src/; find $(PKG) -type f -name '*.go' -and -not -path '*vendor*'|xargs -n1 dirname|uniq`

qa: vet critic test

mocks: $(GOMOCK)
	# See https://github.com/golang/mock for details
	$(MOCKGEN) -source sam/service/service_test.go -destination sam/service/service_mock_test.go -package service


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

$(DEP):
	$(GOGET) github.com/tools/godep

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

clean:
	rm -f $(SPEC) $(PROTOC) $(REALIZE) $(GOCRITIC) $(GOTEST)


