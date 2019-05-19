.PHONY: help docker docker-push realize dep dep.update test test.messaging test.compose qa critic vet codegen integration

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

BASEPKGS = system compose messaging
IMAGES   = system compose messaging

########################################################################################################################
# Tool bins
DEP      = $(GOPATH)/bin/dep
REALIZE   = ${GOPATH}/bin/realize
GOTEST    = ${GOPATH}/bin/gotest
GOCRITIC  = ${GOPATH}/bin/gocritic
MOCKGEN   = ${GOPATH}/bin/mockgen

help:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - docker: builds docker images locally
	@echo - docker-push: push built images
	@echo
	@echo - vet - run go vet on all code
	@echo - critic - run go critic on all code
	@echo - test.compose - individual package unit tests
	@echo - test.messaging - individual package unit tests
	@echo - test - run all available unit tests
	@echo - qa - run vet, critic and test on code
	@echo


docker: $(IMAGES:%=docker-image.%)

docker-image.%: Dockerfile.%
	@ docker build --no-cache --rm -f Dockerfile.$* -t crusttech/$*:latest .

docker-push: $(IMAGES:%=docker-push.%)

docker-push.%: Dockerfile.%
	@ docker push crusttech/$*:latest


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.update: $(DEP)
	$(DEP) ensure -update -v

dep: $(DEP)
	$(DEP) ensure -v

codegen:
	./codegen.sh

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

########################################################################################################################
# QA

test:
	# Run basic unit tests
	$(GO) test ./cmd/... ./internal/... ./compose/... ./messaging/... ./system/...

test-coverage:
	overalls -project=github.com/crusttech/crust -covermode=count -debug -- -coverpkg=./... --tags=integration
	mv overalls.coverprofile coverage.txt

test.internal: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/...
	$(GO) tool cover -func=.cover.out

test.messaging: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./messaging/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.pubsub: $(GOTEST)
	$(GOTEST) -run PubSubMemory -covermode count -coverprofile .cover.out -v ./messaging/internal/repository/pubsub*.go ./messaging/internal/repository/flags*.go ./messaging/internal/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/internal\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.events: $(GOTEST)
	$(GOTEST) -run Events -covermode count -coverprofile .cover.out -v ./messaging/internal/repository/events*.go ./messaging/internal/repository/flags*.go ./messaging/internal/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/internal\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.compose: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./compose/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.system: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./system/internal/repository/... ./system/internal/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.mail: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/mail/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.store: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/store/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.cross-dep:
	# Outputs cross-package imports that should not be there.
	grep -rE "crust/(compose|messaging)/" system || exit 0
	grep -rE "crust/(system|messaging)/" compose || exit 0
	grep -rE "crust/(system|compose)/" messaging || exit 0

integration:
	# Run drone's integration pipeline
	rm -f build/gen*
	drone exec --pipeline integration


vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

qa: vet critic test

mocks: $(GOMOCK)
	# Cleanup all pre-generated
	find . -name '*_mock_test.go' -delete
	rm -rf system/internal/repository/mocks && mkdir -p system/internal/repository/mocks

	$(MOCKGEN) -package repository -source system/internal/repository/user.go         -destination system/internal/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/internal/repository/credentials.go  -destination system/internal/repository/mocks/credentials.go

	$(MOCKGEN) -package mail  -source internal/mail/mail.go                           -destination internal/mail/mail_mock_test.go



########################################################################################################################
# Toolset

$(GOTEST):
	$(GOGET) github.com/rakyll/gotest

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(DEP):
	$(GOGET) github.com/tools/godep

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

clean:
	rm -f $(REALIZE) $(GOCRITIC) $(GOTEST)
