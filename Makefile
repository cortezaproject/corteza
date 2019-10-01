.PHONY: help docker docker-push realize dep dep.update qa critic vet codegen integration

GO        = go
GOGET     = $(GO) get -u
GOTEST    ?= go test

BASEPKGS = system compose messaging
IMAGES   = corteza-server-system corteza-server-compose corteza-server-messaging corteza-server
TESTABLE = messaging system compose pkg internal

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


########################################################################################################################
# Tool bins
DEP         = $(GOPATH)/bin/dep
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
	@echo - docker-images: builds docker images locally
	@echo - docker-push:   push built images
	@echo
	@echo - vet - run go vet on all code
	@echo - critic - run go critic on all code
	@echo - test.all - run all tests
	@echo - test.unit - run all unit tests
	@echo - test.integration - run all integration tests
	@echo
	@echo See tests/README.md for more info
	@echo


docker-images: $(IMAGES:%=docker-image.%)

docker-image.%: Dockerfile.%
	@ docker build --no-cache --rm -f Dockerfile.$* -t cortezaproject/$*:latest .

docker-push: $(IMAGES:%=docker-push.%)

docker-push.%: Dockerfile.%
	@ docker push cortezaproject/$*:latest


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.update: $(DEP)
	$(DEP) ensure -update -v

dep: $(DEP)
	$(DEP) ensure -v

codegen: $(PROTOGEN)
	./codegen.sh

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

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

# Runs ALL tests
test.unit:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_unit)

# Testing pkg
test.pkg:
	$(GOTEST) $(TEST_FLAGS) $(TEST_SUITE_pkg)


# Testing messaging, system, compose
test.%:
	$(GOTEST) $(TEST_FLAGS) ./$*/...

test: test.unit

# Outputs cross-package imports that should not be there.
test.cross-dep:
	@ grep -rE "github.com/cortezaproject/corteza-server/(compose|messaging)/" system || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|messaging)/" compose || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|compose)/" messaging || exit 0
	@ grep -rE "github.com/cortezaproject/corteza-server/(system|compose|messaging)/" pkg || exit 0

# Drone tasks
# Run drone's integration pipeline
drone.integration:
	rm -f build/gen*
	drone exec --pipeline integration


vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./pkg/... ./system/... ./messaging/... ./compose/...

qa: vet critic test

mocks: $(GOMOCK)
	# Cleanup all pre-generated
	find . -name '*_mock_test.go' -delete
	rm -rf system/repository/mocks && mkdir -p system/repository/mocks
	rm -rf compose/service/mocks && mkdir -p compose/service/mocks


	$(MOCKGEN) -package repository -source system/repository/user.go         -destination system/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/repository/credentials.go  -destination system/repository/mocks/credentials.go

	$(MOCKGEN) -package service_mocks -source compose/service/automation_runner.go -destination compose/service/mocks/automation_runner.go

	$(MOCKGEN) -package mail  -source pkg/mail/mail.go                           -destination pkg/mail/mail_mock_test.go



########################################################################################################################
# Toolset

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(DEP):
	$(GOGET) github.com/tools/godep

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

$(STATICCHECK):
	$(GOGET) honnef.co/go/tools/cmd/staticcheck

$(PROTOGEN):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

$(NODEMON):
	npm install -g nodemon

clean:
	rm -f $(REALIZE) $(GOCRITIC)
