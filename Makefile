.PHONY: pack build help realize qa critic vet codegen integration

GO         = go
GOGET      = $(GO) get -u
GOTEST    ?= go test
GOFLAGS   ?= -mod=vendor

export GOFLAGS

BUILD_FLAVOUR         ?= corteza
BUILD_APPS            ?= system compose messaging monolith
BUILD_TIME            ?= $(shell date +%FT%T%z)
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_ARCH            ?= $(shell go env GOARCH)
BUILD_OS              ?= $(shell go env GOOS)
BUILD_OS_is_windows    = $(filter windows,$(BUILD_OS))
BUILD_DEST_DIR        ?= build
BUILD_NAME             = $(BUILD_FLAVOUR)-server-$*-$(BUILD_VERSION)-$(BUILD_OS)-$(BUILD_ARCH)
BUILD_BIN_NAME         = $(BUILD_NAME)$(if $(BUILD_OS_is_windows),.exe,)

RELEASE_BASEDIR        = $(BUILD_DEST_DIR)/pkg/$(BUILD_FLAVOUR)-server-$*
RELEASE_NAME           = $(BUILD_NAME).tar.gz
RELEASE_EXTRA_FILES   ?= README.md LICENSE CONTRIBUTING.md DCO .env.example
RELEASE_PKEY          ?= .upload-rsa

LDFLAGS_VERSION        = -X github.com/cortezaproject/corteza-server/pkg/version.Version=$(BUILD_VERSION)
LDFLAGS_EXTRA         ?=
LDFLAGS                = -ldflags "$(LDFLAGS_VERSION) $(LDFLAGS_EXTRA)"


# Run watcher with a different event-trigger delay, eg:
# $> WATCH_DELAY=5s make watch.test.integration
WATCH_DELAY ?= 1s

# Run go test cmd with flags, eg:
# $> TEST_FLAGS="-v" make test.integration
# $> TEST_FLAGS="-v -run SpecialTest" make test.integration
TEST_FLAGS ?=

COVER_MODE    ?= count
COVER_PROFILE ?= .cover.out
COVER_FLAGS   ?= -covermode=$(COVER_MODE)  -coverprofile=$(COVER_PROFILE)

# Cover package maps for tests tasks
COVER_PKGS_messaging   = ./messaging/...
COVER_PKGS_system      = ./system/...
COVER_PKGS_compose     = ./compose/...
COVER_PKGS_pkg         = ./pkg/...
COVER_PKGS_all         = $(COVER_PKGS_pkg),$(COVER_PKGS_messaging),$(COVER_PKGS_system),$(COVER_PKGS_compose)
COVER_PKGS_integration = $(COVER_PKGS_all)

TEST_SUITE_pkg         = ./pkg/...
TEST_SUITE_services    = ./compose/... ./messaging/... ./system/...
TEST_SUITE_unit        = $(TEST_SUITE_pkg) $(TEST_SUITE_services)
TEST_SUITE_integration = ./tests/...
TEST_SUITE_all         = $(TEST_SUITE_unit) $(TEST_SUITE_integration)

# Dev Support apps settings
DEV_MINIO_PORT        ?= 9000
DEV_MAILHOG_SMTP_ADDR ?= 1025
DEV_MAILHOG_HTTP_ADDR ?= 8025

DOCKER                ?= docker

########################################################################################################################
# Tool bins
REALIZE     = ${GOPATH}/bin/realize
GOCRITIC    = ${GOPATH}/bin/gocritic
MOCKGEN     = ${GOPATH}/bin/mockgen
STATICCHECK = ${GOPATH}/bin/staticcheck
PROTOGEN    = ${GOPATH}/bin/protoc-gen-go

# Using nodemon in development environment for "watch.*" tasks
# https://nodemon.io
NODEMON      = /usr/local/bin/nodemon
WATCHER      = $(NODEMON) --delay ${WATCH_DELAY} -e go -w . --exec

help:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - build             build all apps
	@echo - build.<app>       build a specific app
	@echo - vet               run go vet on all code
	@echo - critic            run go critic on all code
	@echo - test.all          run all tests
	@echo - test.unit         run all unit tests
	@echo - test.integration  run all integration tests
	@echo
	@echo See tests/README.md for more info on running tests
	@echo

########################################################################################################################
# Building & packing

build: $(addprefix build., $(BUILD_APPS))

build.%: cmd/%
	GOOS=$(BUILD_OS) GOARCH=$(BUILD_ARCH) go build $(LDFLAGS) -o $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME) cmd/$*/main.go

release.%: $(addprefix build., %)
	@ mkdir -p $(RELEASE_BASEDIR) $(RELEASE_BASEDIR)/bin
	@ cp $(RELEASE_EXTRA_FILES) $(RELEASE_BASEDIR)/
	@ cp $(BUILD_DEST_DIR)/$(BUILD_BIN_NAME) $(RELEASE_BASEDIR)/bin/$(BUILD_FLAVOUR)-server-$*
	@ tar -C $(dir $(RELEASE_BASEDIR)) -czf $(BUILD_DEST_DIR)/$(RELEASE_NAME) $(notdir $(RELEASE_BASEDIR))

release: $(addprefix release.,$(BUILD_APPS))

release-clean:
	@ rm -rf $(RELEASE_BASEDIR)

upload: $(RELEASE_PKEY)
	@ echo "put $(BUILD_DEST_DIR)/*.tar.gz" | sftp -q -i $(RELEASE_PKEY) $(RELEASE_SFTP_URI)
	@ rm -f $(RELEASE_PKEY)

$(RELEASE_PKEY):
	@ echo $(RELEASE_SFTP_KEY) | base64 -d > $(RELEASE_PKEY)
	@ chmod 0400 $@

########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

codegen: $(PROTOGEN)
	./codegen.sh

mailhog.up:
	$(DOCKER) run --rm --publish $(DEV_MAILHOG_HTTP_ADDR):8025 --publish $(DEV_MAILHOG_SMTP_ADDR):1025 mailhog/mailhog

minio.up:
	# Runs temp minio server
	# No volume mounts because we do not want the data to persist
	$(DOCKER) run --rm --publish $(DEV_MINIO_PORT):9000 --env-file .env minio/minio server /data

watch.test.%: $(NODEMON)
	# Development helper - watches for file
	# changes & reruns  tests
	$(WATCHER) "make test.$* || exit 0"

########################################################################################################################
# Quality Assurance

# Adds -coverprofile flag to test flags
# and executes test.cover... task
test.coverprofile.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverprofile=$(COVER_PROFILE)" make test.cover.$*

# Adds -coverpkg flag
test.cover.%:
	@ TEST_FLAGS="$(TEST_FLAGS) -coverpkg=$(COVER_PKGS_$*)" make test.$*

# Runs integration tests
test.integration:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_integration)

# Runs one suite from integration tests
test.integration.%:
	$(GOTEST) $(TEST_FLAGS) ./tests/$*/...

# Runs ALL tests
test.all:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_all)

# Unit testing testing messaging, system or compose
test.unit.%:
	$(GOTEST) $(TEST_FLAGS) ./$*/...

# Runs ALL tests
test.unit:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_unit)

# Testing pkg
test.pkg:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_pkg)

test: test.unit

# Outputs cross-package imports that should not be there.
test.cross-dep:
	@ grep -rE "github.com/cortezaproject/corteza-server/(compose|messaging)/" system || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|messaging)/" compose || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|compose)/" messaging || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|compose|messaging)/" pkg || exit 0

vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./pkg/... ./system/... ./messaging/... ./compose/...

qa: vet critic test

mocks: $(MOCKGEN)
	# Cleanup all pre-generated
	rm -rf system/repository/mocks && mkdir -p system/repository/mocks
	rm -rf compose/service/mocks && mkdir -p compose/service/mocks

	$(MOCKGEN) -package repository -source system/repository/user.go         -destination system/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/repository/credentials.go  -destination system/repository/mocks/credentials.go

	$(MOCKGEN) -package mail  -source pkg/mail/mail.go                           -destination pkg/mail/mail_mock_test.go



########################################################################################################################
# Toolset

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/mockgen

$(STATICCHECK):
	$(GOGET) honnef.co/go/tools/cmd/staticcheck

$(PROTOGEN):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

$(NODEMON):
	npm install -g nodemon

clean:
	rm -f $(REALIZE) $(GOCRITIC)
