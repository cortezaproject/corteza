GO         = go
GOINSTALL  = $(GO) install
GOFLAGS   ?= -trimpath
GOPATH    ?= $(HOME)/go
ROOT_DIR   = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))


CGO_ENABLED = 1

export GOFLAGS
export CGO_ENABLED

########################################################################################################################
# Tool bins
GOCRITIC      = $(GOPATH)/bin/gocritic
MOCKGEN       = $(GOPATH)/bin/mockgen
GOTEST        = $(GOPATH)/bin/gotest
STATICCHECK   = $(GOPATH)/bin/staticcheck
PROTOGEN      = $(GOPATH)/bin/protoc-gen-go
PROTOGEN_GRPC = $(GOPATH)/bin/protoc-gen-go-grpc
GIN           = $(GOPATH)/bin/gin
STATIK        = $(GOPATH)/bin/statik
MODOUTDATED   = $(GOPATH)/bin/go-mod-outdated
CUE           = $(GOPATH)/bin/cue
CODEGEN       = $(GOPATH)/bin/corteza-codegen
JSONTPLEXEC   = $(GOPATH)/bin/corteza-json-tpl-exec

PROTOC      = /usr/local/bin/protoc
FSWATCH     = /usr/local/bin/fswatch

# fswatch is intentionally left out...
BINS         = $(GOCRITIC) $(MOCKGEN) $(GOTEST) $(STATICCHECK) $(PROTOGEN) $(GIN) $(STATIK)

# Corteza's tools
CORTEZA_BINS = $(CODEGEN) $(JSONTPLEXEC)


########################################################################################################################
# Go Toolset
$(GOCRITIC):
	$(GOINSTALL) github.com/go-critic/go-critic@latest

$(MOCKGEN):
	$(GOINSTALL) github.com/golang/mock/mockgen@latest

$(STATICCHECK):
	$(GOINSTALL) honnef.co/go/tools/cmd/staticcheck@latest

$(PROTOGEN):
	$(GOINSTALL) github.com/golang/protobuf/protoc-gen-go@latest

$(PROTOGEN_GRPC):
	$(GOINSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$(GIN):
	$(GOINSTALL) github.com/codegangsta/gin@latest

$(GOTEST):
	$(GOINSTALL) github.com/rakyll/gotest@latest

$(STATIK):
	$(GOINSTALL) github.com/goware/statik@latest

$(MODOUTDATED):
	$(GOINSTALL) github.com/psampaz/go-mod-outdated@latest

$(CUE):
	$(GOINSTALL) cuelang.org/go/cmd/cue@latest

$(CODEGEN):
	$(GO) build -o $@ $(ROOT_DIR)/cmd/codegen/main.go

$(JSONTPLEXEC):
	$(GO) build -o $@ $(ROOT_DIR)/codegen/tool

clean-tools:
	rm -f $(BINS)

clean-corteza-tools:
	rm -f $(CORTEZA_BINS)

########################################################################################################################

COLOUR_GREEN=\033[0;32m
COLOUR_RED=\033[0;31m
COLOUR_BLUE=\033[0;34m
COLOUR_END=\033[0m


#
