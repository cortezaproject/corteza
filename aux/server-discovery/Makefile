.PHONY: pack build help qa critic vet codegen provision docs build auth webapp

include Makefile.inc

BUILD_FLAVOUR         ?= corteza
BUILD_TIME            ?= $(shell date +%FT%T%z)
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_ARCH            ?= $(shell go env GOARCH)
BUILD_OS              ?= $(shell go env GOOS)
BUILD_OS_is_windows    = $(filter windows,$(BUILD_OS))
BUILD_DEST_DIR        ?= build
BUILD_NAME             = $(BUILD_FLAVOUR)-server-discovery-$(BUILD_VERSION)-$(BUILD_OS)-$(BUILD_ARCH)
BUILD_BIN_NAME         = $(BUILD_NAME)$(if $(BUILD_OS_is_windows),.exe,)

RELEASE_BASEDIR        = $(BUILD_DEST_DIR)/pkg/$(BUILD_FLAVOUR)-server
RELEASE_NAME           = $(BUILD_NAME).tar.gz
RELEASE_EXTRA_FILES   ?= README.md LICENSE CONTRIBUTING.md DCO .env.example
RELEASE_PKEY          ?= .upload-rsa

LDFLAGS_VERSION        = -X github.com/cortezaproject/corteza-server/pkg/version.Version=$(BUILD_VERSION)
LDFLAGS_EXTRA         ?=
LDFLAGS                = -ldflags "$(LDFLAGS_VERSION) $(LDFLAGS_EXTRA)"

########################################################################################################################
# Building & packing

build: $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME)

$(BUILD_DEST_DIR)/$(BUILD_BIN_NAME):
		GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build $(LDFLAGS) -o $@ *.go

release: build $(BUILD_DEST_DIR)/$(RELEASE_NAME)

$(BUILD_DEST_DIR)/$(RELEASE_NAME):
	@ mkdir -p $(RELEASE_BASEDIR) $(RELEASE_BASEDIR)/bin
	@ cp $(RELEASE_EXTRA_FILES) $(RELEASE_BASEDIR)/
	@ cp $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME) $(RELEASE_BASEDIR)/bin/$(BUILD_FLAVOUR)-server
	tar -C $(dir $(RELEASE_BASEDIR)) -czf $(BUILD_DEST_DIR)/$(RELEASE_NAME) $(notdir $(RELEASE_BASEDIR))

release-clean:
	rm -rf $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME)
	rm -rf $(BUILD_DEST_DIR)/$(RELEASE_NAME)

upload: $(RELEASE_PKEY)
	@ echo "put $(BUILD_DEST_DIR)/*.tar.gz" | sftp -q -o "StrictHostKeyChecking no" -i $(RELEASE_PKEY) $(RELEASE_SFTP_URI)
	@ rm -f $(RELEASE_PKEY)

$(RELEASE_PKEY):
	@ echo $(RELEASE_SFTP_KEY) | base64 -d > $@
	@ chmod 0400 $@

########################################################################################################################
# Development

watch: $(GIN)
	$(GIN) $(GIN_ARGS) run -- serve
