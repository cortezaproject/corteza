.PHONY: generate

GOPATH   ?= $(go env GOPATH)
BRANCH   ?= develop
ZIP      ?= $(BRANCH).zip
DIR      ?= corteza-protobuf-$(BRANCH)
PROTO    ?= service-corredor.proto
URL      ?= https://github.com/cortezaproject/corteza-protobuf/archive/$(ZIP)

PARTS    ?= compose system
OUTPUTS  ?= $(addsuffix /static.go, $(PARTS))
PROTOC   ?= /usr/local/bin/protoc
CURL     ?= /usr/bin/curl
UNZIP    ?= /usr/bin/unzip

PROTOGEN      ?= $(GOPATH)/bin/protoc-gen-go
PROTOGEN_GRPC ?= $(GOPATH)/bin/protoc-gen-go-grpc


generate:
	PATH=$(PATH):$(GOPATH)/bin protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO)

$(PROTO): $(ZIP)
	$(UNZIP) -qq -o $(ZIP)
	mv $(DIR)/$(PROTO) $(PROTO)
	rm -rf $(DIR)

$(ZIP):
	curl -o $(ZIP) -s --location "$(URL)"

clear:
	rm -f $(ZIP)

deps: $(PROTOC) $(PROTOGEN) $(PROTOGEN_GRPC)

$(PROTOC):
	$(MAKE) --directory=../../ $(PROTOC)

$(PROTOGEN):
	$(MAKE) --directory=../../ $(PROTOGEN)

$(PROTOGEN_GRPC):
	$(MAKE) --directory=../../ $(PROTOGEN_GRPC)
