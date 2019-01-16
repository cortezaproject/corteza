.PHONY: nothing docker docker-push realize dep dep.update test test.rbac test.sam test.crm qa critic vet codegen

PKG       = "github.com/$(shell cat .project)"

GO        = go
GOGET     = $(GO) get -u

BASEPKGS = rbac system crm sam
IMAGES   = system crm sam

########################################################################################################################
# Tool bins
DEP      = $(GOPATH)/bin/dep
SPEC      = $(GOPATH)/bin/spec
PROTOC    = $(GOPATH)/bin/protoc-gen-go
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
	@echo - test.sam - individual package unit tests
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

codegen: $(SPEC) dep.codegen
	@PATH=${PATH}:${GOPATH}/bin ./codegen.sh

mailhog.up:
	docker run --rm --publish 8025:8025 --publish 1025:1025 mailhog/mailhog

########################################################################################################################
# QA

test: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./...
	$(GO) tool cover -func=.cover.out

test.sam: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./sam/repository/... ./sam/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.sam.db: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./sam/db/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.pubsub: $(GOTEST)
	$(GOTEST) -run PubSubMemory -covermode count -coverprofile .cover.out -v ./sam/repository/pubsub*.go ./sam/repository/flags*.go ./sam/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/sam\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.events: $(GOTEST)
	$(GOTEST) -run Events -covermode count -coverprofile .cover.out -v ./sam/repository/events*.go ./sam/repository/flags*.go ./sam/repository/error*.go
	perl -pi -e 's/command-line-arguments/.\/sam\/repository/g' .cover.out
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.crm: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./crm/service/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.crm.db: $(GOTEST)
	cd crm/db && $(GO) generate && cd ../..
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./crm/db/...
	$(GO) tool cover -func=.cover.out | grep --color "^\|[^0-9]0.0%"

test.system: $(GOTEST)
	$(GOTEST) -covermode count -coverprofile .cover.out -v ./system/repository/...
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

	$(MOCKGEN) -package service -source sam/service/attachment.go   -destination sam/service/attachment_mock_test.go
	$(MOCKGEN) -package service -source sam/service/channel.go      -destination sam/service/channel_mock_test.go
	$(MOCKGEN) -package service -source sam/service/message.go      -destination sam/service/message_mock_test.go
	$(MOCKGEN) -package service -source system/service/organisation.go -destination system/service/organisation_mock_test.go
	$(MOCKGEN) -package service -source system/service/team.go         -destination system/service/team_mock_test.go
	$(MOCKGEN) -package service -source system/service/user.go         -destination system/service/user_mock_test.go

	$(MOCKGEN) -package mail -source internal/mail/mail.go -destination internal/mail/mail_mock_test.go

	mkdir -p system/repository/mocks
	$(MOCKGEN) -package repository -source system/repository/user.go         -destination system/repository/mocks/user.go
	$(MOCKGEN) -package repository -source system/repository/credentials.go  -destination system/repository/mocks/credentials.go


########################################################################################################################
# Toolset

$(GOTEST):
	$(GOGET) github.com/rakyll/gotest

$(REALIZE):
	$(GOGET) github.com/tockins/realize

$(SPEC):
	$(GOGET) github.com/titpetric/spec/cmd/spec

$(PROTOC):
	$(GOGET) github.com/golang/protobuf/protoc-gen-go

$(GOCRITIC):
	$(GOGET) github.com/go-critic/go-critic/...

$(DEP):
	$(GOGET) github.com/tools/godep

$(MOCKGEN):
	$(GOGET) github.com/golang/mock/gomock
	$(GOGET) github.com/golang/mock/mockgen

clean:
	rm -f $(SPEC) $(PROTOC) $(REALIZE) $(GOCRITIC) $(GOTEST)


