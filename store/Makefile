
GOTOOL    ?= go tool
GOTEST    ?= go test
GOFLAGS   ?= -mod=vendor

COVER_MODE    ?= count
COVER_PROFILE ?= .cover.out
COVER_FLAGS   ?= -covermode=$(COVER_MODE)  -coverprofile=$(COVER_PROFILE)


cover:
	$(GOTEST) -coverprofile=$(COVER_PROFILE) -coverpkg=./... ./tests/...
	$(GOTOOL) cover -html=$(COVER_PROFILE)
