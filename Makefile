.PHONY: pack build help qa critic vet codegen provision docs build auth webapp

include Makefile.inc

BUILD_FLAVOUR         ?= corteza
BUILD_TIME            ?= $(shell date +%FT%T%z)
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_ARCH            ?= $(shell go env GOARCH)
BUILD_OS              ?= $(shell go env GOOS)
BUILD_OS_is_windows    = $(filter windows,$(BUILD_OS))
BUILD_DEST_DIR        ?= build
BUILD_NAME             = $(BUILD_FLAVOUR)-server-$(BUILD_VERSION)-$(BUILD_OS)-$(BUILD_ARCH)
BUILD_BIN_NAME         = $(BUILD_NAME)$(if $(BUILD_OS_is_windows),.exe,)

RELEASE_BASEDIR        = $(BUILD_DEST_DIR)/pkg/$(BUILD_FLAVOUR)-server
RELEASE_NAME           = $(BUILD_NAME).tar.gz
RELEASE_EXTRA_FILES   ?= README.md LICENSE CONTRIBUTING.md DCO .env.example
RELEASE_PKEY          ?= .upload-rsa

LDFLAGS_VERSION        = -X github.com/cortezaproject/corteza-server/pkg/version.Version=$(BUILD_VERSION)
LDFLAGS_EXTRA         ?=
LDFLAGS                = -ldflags "$(LDFLAGS_VERSION) $(LDFLAGS_EXTRA)"

# Run go test cmd with flags, eg:
# $> make test.integration TEST_FLAGS="-v"
# $> make test.integration TEST_FLAGS="-v -run SpecialTest"
TEST_FLAGS ?=

COVER_MODE    ?= count
COVER_PROFILE ?= .cover.out
COVER_FLAGS   ?= -covermode=$(COVER_MODE)  -coverprofile=$(COVER_PROFILE)

# Cover package maps for tests tasks
COVER_PKGS_system      = ./system/...
COVER_PKGS_compose     = ./compose/...
COVER_PKGS_federation  = ./federation/...
COVER_PKGS_automation  = ./automation/...
COVER_PKGS_pkg         = ./pkg/...
COVER_PKGS_all         = $(COVER_PKGS_pkg),$(COVER_PKGS_system),$(COVER_PKGS_compose),$(COVER_PKGS_federation),$(COVER_PKGS_automation)
COVER_PKGS_integration = $(COVER_PKGS_all)

TEST_SUITE_app         = ./app/...
TEST_SUITE_pkg         = ./pkg/...
TEST_SUITE_services    = ./compose/... ./system/... ./federation/... ./auth/... ./automation/...
TEST_SUITE_unit        = $(TEST_SUITE_pkg) $(TEST_SUITE_app) $(TEST_SUITE_services)
TEST_SUITE_integration = ./tests/...
TEST_SUITE_store       = ./store/tests/...
TEST_SUITE_all         = $(TEST_SUITE_unit) $(TEST_SUITE_integration) $(TEST_SUITE_store)

# Dev Support apps settings
DEV_MINIO_PORT        ?= 9000
DEV_MAILHOG_SMTP_ADDR ?= 1025
DEV_MAILHOG_HTTP_ADDR ?= 8025

GIN_ARG_LADDR ?= localhost
GIN_ARGS      ?= --laddr $(GIN_ARG_LADDR) --build cmd/corteza --immediate --bin build/gin-bin

GIN_CORTEZA_ENV_FILE ?= .env
GIN_CORTEZA_BIN_ARGS ?= --env-file $(GIN_CORTEZA_ENV_FILE)

DOCKER                ?= docker

help:
	@echo ""
	@echo " Usage: make [target]"
	@echo ""
	@echo " - build             build all apps"
	@echo " - build.<app>       build a specific app"
	@echo " - vet               run go vet on all code"
	@echo " - critic            run go critic on all code"
	@echo " - test.all          run all tests"
	@echo " - test.unit         run all unit tests"
	@echo " - test.integration  run all integration tests"
	@echo ""
	@echo " See tests/README.md for more info on running tests"
	@echo ""

########################################################################################################################
# Building & packing

build: $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME)

$(BUILD_DEST_DIR)/$(BUILD_BIN_NAME):
		GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build $(LDFLAGS) -o $@ cmd/corteza/main.go

release: build $(BUILD_DEST_DIR)/$(RELEASE_NAME)

$(BUILD_DEST_DIR)/$(RELEASE_NAME):
	@ mkdir -p $(RELEASE_BASEDIR) $(RELEASE_BASEDIR)/bin
	@ cp $(RELEASE_EXTRA_FILES) $(RELEASE_BASEDIR)/
	@ cp -r provision $(RELEASE_BASEDIR)
	@ rm -f $(RELEASE_BASEDIR)/provision/README.adoc $(RELEASE_BASEDIR)/provision/update.sh
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
	$(GIN) $(GIN_ARGS) -- $(GIN_CORTEZA_BIN_ARGS) serve

realize: watch # BC

mailhog.up:
	$(DOCKER) run --rm --publish $(DEV_MAILHOG_HTTP_ADDR):8025 --publish $(DEV_MAILHOG_SMTP_ADDR):1025 mailhog/mailhog

minio.up:
	# Runs temp minio server
	# No volume mounts because we do not want the data to persist
	$(DOCKER) run --rm --publish $(DEV_MINIO_PORT):9000 --env-file .env minio/minio server /data

# Development helper - reruns test when files change
#
# make watch.test.unit
# make watch.test.pkg
# make watch.test.all
# make watch.test.pkg TEST_FLAGS="-v"
watch.test.%: $(FSWATCH)
	( make test.$* || exit 0 ) && ( $(FSWATCH) -o . | xargs -n1 -I{} make test.$* )

watch.test: watch.test.unit

# See codegen/README.md for details
codegen:
	@ make -C codegen all

cue.fmt: $(CUE)
	$(CUE) fmt -v codegen/*.cue
	$(CUE) fmt -v codegen/schema/*.cue

	$(CUE) fmt -v app/*.cue
	$(CUE) fmt -v app/options/*.cue

	$(CUE) fmt system/*.cue
	$(CUE) fmt compose/*.cue
	$(CUE) fmt automation/*.cue
	$(CUE) fmt federation/*.cue

codegen-legacy: $(CODEGEN)
	@ $(CODEGEN) -v

watch.codegen-legacy: $(CODEGEN)
	@ $(CODEGEN) -w -v

clean.codegen-legacy:
	rm -f $(CODEGEN)

provision:
	$(MAKE) --directory=provision clean all

webapp:
	@ $(MAKE) --directory=webapp

locale.update:
	$(GO) get github.com/cortezaproject/corteza-locale@2022.9.x
	$(GO) mod vendor
	git add --all vendor/github.com/cortezaproject
	git commit -m 'Update corteza-locale dep' vendor/github.com/cortezaproject vendor/modules.txt go.mod go.sum


vendors.check: $(MODOUTDATED)
	$(GO) list -mod=mod -u -m -json all | $(MODOUTDATED) -update -direct

#######################################################################################################################
# Quality Assurance

# Adds -coverprofile flag to test flags
# and executes test.cover... task
test.coverprofile.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverprofile=$(COVER_PROFILE)" make test.cover.$*

# Adds -coverpkg flag
test.cover.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverpkg=$(COVER_PKGS_$*)" make test.$*

# Runs integration tests
test.integration: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_integration)

# Runs one suite from integration tests
test.integration.%: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) ./tests/$*/...

# Runs store tests
test.store: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_store)

# Runs ALL tests
test.all:
	@echo "pass"
#test.all: $(GOTEST)
#	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_all)

# Unit testing testing, system or compose
test.unit.%: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) ./$*/...

# Runs ALL tests
test.unit: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_unit)

# Testing pkg
test.pkg: $(GOTEST)
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_pkg)

# Test defaults to test.unit
test: test.unit


vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./pkg/... ./system/... ./compose/... ./automation/...

qa: vet critic test

mocks: $(MOCKGEN)
	$(MOCKGEN) -package mail -source pkg/mail/mail.go -destination pkg/mail/mail_mock_test.go


########################################################################################################################
# Toolset

# @todo this will most likely need some special care for other platforms
$(FSWATCH):
	ifeq ($(UNAME_S),Darwin)
		brew install fswatch
	endif

# https://grpc.io/docs/protoc-installation/
# @todo $ apt install -y protobuf-compiler
$(PROTOC):
	ifeq ($(UNAME_S),Darwin)
		brew install protobuf
	endif

#
