.PHONY: nothing docker docker-push realize dep dep.update protobuf test test.rbac test.sam test.crm qa critic vet codegen

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

nothing:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - docker: builds docker images locally
	@echo - docker-push: push built images
	@echo
	@echo - vet - run go vet on all code
	@echo - critic - run go critic on all code
	@echo - test.crm - individual package unit tests
	@echo - test.sam - individual package unit tests
	@echo - test.rbac - individual package unit tests
	@echo - test - run all available unit tests
	@echo - qa - run vet, critic and test on code
	@echo

docker:
	docker build --no-cache --rm --build-arg APP=sam -f docker/Dockerfile -t crusttech/sam .
	docker build --no-cache --rm --build-arg APP=crm -f docker/Dockerfile -t crusttech/crm .

docker-push:
	docker push crusttech/sam
	docker push crusttech/crm

########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.update: $(DEP)
	$(DEP) ensure -update -v

dep: $(DEP)
	$(DEP) ensure -v

codegen: $(SPEC)
	@PATH=${PATH}:${GOPATH}/bin ./codegen.sh

protobuf: $(PROTOC)
	# @todo this needs work (it hangs and outputs nothing)
	$(PROTOC) --go_out=plugins=grpc:. -I. sam/chat/*.proto


########################################################################################################################
# QA

test: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./...
	$(GO) tool cover -func=.cover.out

test.sam: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./sam/repository/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.crm: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./crm/repository/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.rbac: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./rbac/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.store: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./store/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

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


