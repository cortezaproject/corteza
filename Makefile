.PHONY: help docker docker-push realize dep dep.update test test.messaging test.compose qa critic vet codegen integration

GO        = go
GOGET     = $(GO) get -u
GOTEST    ?= go test

BASEPKGS = system compose messaging
IMAGES   = corteza-server-system corteza-server-compose corteza-server-messaging corteza-server

# Run watcher with a different event-trigger delay, eg:
# $> WATCH_DELAY=5s make watch.test.integration
WATCH_DELAY ?= 1s

# Run go test cmd with flags, eg:
# $> TEST_FLAGS="-v" make test.integration
# $> TEST_FLAGS="-v -run SpecialTest" make test.integration
TEST_FLAGS ?=

TEST_INTEGRATION_COVER_PROFILE_OUT ?= .integration.cover.out

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
	@echo - test.compose - individual package unit tests
	@echo - test.messaging - individual package unit tests
	@echo - test - run all available unit tests
	@echo - qa - run vet, critic and test on code
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

watch.test.integration: $(NODEMON)
	# Development helper - watches for file
	# changes & reruns integration tests
	$(WATCHER) "make test.integration || exit 0"

watch.test.integration.coverage: $(NODEMON)
	# Development helper - watches for file
	# changes & reruns integration tests
	$(WATCHER) "make test.integration.coverage || exit 0"

########################################################################################################################
# QA

## refactored

test.integration:
	$(GOTEST) $(TEST_FLAGS) ./tests/...

test.integration.coverage:
	$(GOTEST) $(TEST_FLAGS) -covermode=count -coverprofile=$(TEST_INTEGRATION_COVER_PROFILE_OUT) -coverpkg=./... ./tests/...

## old:

test:
	# Run basic unit tests
	$(GOTEST) ./pkg/... ./internal/... ./compose/... ./messaging/... ./system/...

test-coverage:
	overalls -project=github.com/cortezaproject/corteza-server -covermode=count -- -coverpkg=./... --tags=integration -p 1
	mv overalls.coverprofile coverage.txt

test.internal:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/...
	$(GO) tool cover -func=.cover.out

test.messaging:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./messaging/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.pubsub:
	$(GOTEST) -run PubSubMemory -covermode count -coverprofile .cover.out -v ./messaging/repository/pubsub*.go ./messaging/repository/flags*.go ./messaging/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/internal\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.events:
	$(GOTEST) -run Events -covermode count -coverprofile .cover.out -v ./messaging/repository/events*.go ./messaging/repository/flags*.go ./messaging/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/internal\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.compose:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./compose/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.system:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./system/repository/... ./system/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.mail:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/mail/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.store:
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/store/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.cross-dep:
	# Outputs cross-package imports that should not be there.
	grep -rE "github.com/cortezaproject/corteza-server/(compose|messaging)/" system || exit 0
	grep -rE "github.com/cortezaproject/corteza-server/(system|messaging)/" compose || exit 0
	grep -rE "github.com/cortezaproject/corteza-server/(system|compose)/" messaging || exit 0
	grep -rE "github.com/cortezaproject/corteza-server/(system|compose|messaging)/" pkg || exit 0
	grep -rE "github.com/cortezaproject/corteza-server/(system|compose|messaging)/" internal || exit 0

integration:
	# Run drone's integration pipeline
	rm -f build/gen*
	drone exec --pipeline integration


vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./pkg/... ./internal/... ./system/... ./messaging/... ./compose/...

qa: vet critic test

mocks: $(GOMOCK)
	# Cleanup all pre-generated
	find . -name '*_mock_test.go' -delete
	rm -rf system/internal/repository/mocks && mkdir -p system/internal/repository/mocks
	rm -rf compose/internal/service/mocks && mkdir -p compose/internal/service/mocks

	$(MOCKGEN) -package repository -source system/internal/repository/user.go         -destination system/internal/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/internal/repository/credentials.go  -destination system/internal/repository/mocks/credentials.go

	$(MOCKGEN) -package service_mocks -source compose/internal/service/automation_runner.go -destination compose/internal/service/mocks/automation_runner.go

	$(MOCKGEN) -package mail  -source internal/mail/mail.go                           -destination internal/mail/mail_mock_test.go



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
