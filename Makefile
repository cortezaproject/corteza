.PHONY: nothing docker docker-push realize dep dep.update test test.rbac test.messaging test.crm qa critic vet codegen

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

BASEPKGS = rbac system crm messaging
IMAGES   = system crm messaging

########################################################################################################################
# Tool bins
DEP      = $(GOPATH)/bin/dep
REALIZE   = ${GOPATH}/bin/realize
GOTEST    = ${GOPATH}/bin/gotest
GOCRITIC  = ${GOPATH}/bin/gocritic
MOCKGEN   = ${GOPATH}/bin/mockgen

nothing:
	@echo
	@echo Usage: make [target]
	@echo
	@echo - docker: builds docker images locally
	@echo - docker-push: push built images
	@echo
	@echo - vet - run go vet on all code
	@echo - critic - run go critic on all code
	@echo - test.crm - individual package unit tests
	@echo - test.messaging - individual package unit tests
	@echo - test.rbac - individual package unit tests
	@echo - test - run all available unit tests
	@echo - qa - run vet, critic and test on code
	@echo


docker: $(IMAGES:%=docker-image.%)

docker-image.%: %
	@ docker build --no-cache --rm -f Dockerfile.$^ -t crusttech/api-$^:latest .

docker-push: $(IMAGES:%=docker-push.%)

docker-push.%: %
	@ docker push crusttech/api-$^:latest


########################################################################################################################
# Development

realize: $(REALIZE)
	$(REALIZE) start

dep.codegen:
	go install github.com/titpetric/statik

dep.update: $(DEP)
	$(DEP) ensure -update -v

dep: $(DEP)
	$(DEP) ensure -v

codegen: dep.codegen
	@PATH=${PATH}:${GOPATH}/bin ./codegen.sh

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

########################################################################################################################
# QA

test: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./...
	$(GO) tool cover -func=.cover.out

test.messaging: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./messaging/repository/... ./messaging/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.messaging.db: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./messaging/db/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.pubsub: $(GOTEST)
	$(GOTEST) -run PubSubMemory -covermode count -coverprofile .cover.out -v ./messaging/repository/pubsub*.go ./messaging/repository/flags*.go ./messaging/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.events: $(GOTEST)
	$(GOTEST) -run Events -covermode count -coverprofile .cover.out -v ./messaging/repository/events*.go ./messaging/repository/flags*.go ./messaging/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/messaging\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.crm: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./crm/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.crm.db: $(GOTEST)
	cd crm/db && $(GO) generate && cd ../..
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./crm/db/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.system: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./system/repository/... ./system/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.system.db: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./system/db/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"


test.rbac: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/rbac/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.rules: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/rules/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.mail: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/mail/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.rbac.resources: $(GOTEST)
	$(GOTEST) -run Resources -covermode count -coverprofile .cover.out -v ./internal/rbac/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.rbac.sessions: $(GOTEST)
	$(GOTEST) -run Sessions -covermode count -coverprofile .cover.out -v ./internal/rbac/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.store: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./internal/store/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

vet:
	$(GO) vet ./...

critic: $(GOCRITIC)
	$(GOCRITIC) check-project .

qa: vet critic test

mocks: $(GOMOCK)
	# Cleanup all pre-generated
	rm -f */*/*_mock_test.go */*/mocks/*

	# See https://github.com/golang/mock for details
	$(MOCKGEN) -package service -source crm/service/notification.go -destination crm/service/notification_mock_test.go

	$(MOCKGEN) -package service -source messaging/service/attachment.go   -destination messaging/service/attachment_mock_test.go
	$(MOCKGEN) -package service -source messaging/service/channel.go      -destination messaging/service/channel_mock_test.go
	$(MOCKGEN) -package service -source messaging/service/message.go      -destination messaging/service/message_mock_test.go
	$(MOCKGEN) -package service -source system/service/organisation.go -destination system/service/organisation_mock_test.go
	$(MOCKGEN) -package service -source system/service/role.go         -destination system/service/role_mock_test.go
	$(MOCKGEN) -package service -source system/service/user.go         -destination system/service/user_mock_test.go

	$(MOCKGEN) -package mail  -source internal/mail/mail.go        -destination internal/mail/mail_mock_test.go
	$(MOCKGEN) -package rules -source internal/rules/interfaces.go -destination internal/rules/resources_mock_test.go

	mkdir -p system/repository/mocks
	$(MOCKGEN) -package repository -source system/repository/user.go         -destination system/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/repository/credentials.go  -destination system/repository/mocks/credentials.go


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


