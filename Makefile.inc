GO         = go
GOINSTALL  = $(GO) install
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
CODEGEN       = build/codegen

PROTOC      = /usr/local/bin/protoc
FSWATCH     = /usr/local/bin/fswatch

# fswatch is intentionally left out...
BINS = $(GOCRITIC) $(MOCKGEN) $(GOTEST) $(STATICCHECK) $(PROTOGEN) $(GIN) $(STATIK) $(CODEGEN)


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

$(CODEGEN):
	$(GO) build -o $@ cmd/codegen/main.go

clean-tools:
	rm -f $(BINS)
