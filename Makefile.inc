GO         = go
GOGET      = $(GO) get -u
GOFLAGS   ?= -trimpath
GOPATH    ?= $(HOME)/go

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
CODEGEN       = build/codegen

PROTOC      = /usr/local/bin/protoc
FSWATCH     = /usr/local/bin/fswatch

# fswatch is intentionally left out...
BINS = $(GOCRITIC) $(MOCKGEN) $(GOTEST) $(STATICCHECK) $(PROTOGEN) $(GIN) $(STATIK) $(CODEGEN)


########################################################################################################################
# Go Toolset
$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/mockgen

$(STATICCHECK):
	$(GOGET) honnef.co/go/tools/cmd/staticcheck

$(PROTOGEN):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

$(PROTOGEN_GRPC):
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc

$(GIN):
	$(GOGET) github.com/codegangsta/gin

$(GOTEST):
	$(GOGET) github.com/rakyll/gotest

$(STATIK):
	$(GOGET) github.com/goware/statik

$(MODOUTDATED):
	$(GOGET) github.com/psampaz/go-mod-outdated

$(CODEGEN):
	$(GO) build -o $@ cmd/codegen/main.go


clean-tools:
	rm -f $(BINS)
