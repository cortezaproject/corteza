
include ../Makefile.inc

GOTOOL    ?= go tool
GOTEST    ?= go test
GOFLAGS   ?= -mod=vendor

COVER_MODE    ?= count
COVER_PROFILE ?= .cover.out
COVER_FLAGS   ?= -covermode=$(COVER_MODE)  -coverprofile=$(COVER_PROFILE)

# Run go test cmd with flags, eg:
# $> TEST_FLAGS="-v" make test
# $> TEST_FLAGS="-v -run 'Test_Store/.+/ComposeCharts'" make test
TEST_FLAGS ?=

test.store:
	$(GOTEST) $(TEST_FLAGS) ./tests/...

test.cover.store:
	@ $(GOTEST) $(TEST_FLAGS) -coverprofile=$(COVER_PROFILE) -coverpkg=./... ./tests/...
	@ $(GOTOOL) cover -html=$(COVER_PROFILE)
	@ $(GOTOOL) cover
	@ rm $(COVER_PROFILE)

watch.test.%: $(FSWATCH)
	( make test.$* || exit 0 ) && ( $(FSWATCH) -o . | xargs -n1 -I{} make test.$* )

watch.codegen:
	@ make -C ../ $@

codegen:
	@ make -C ../ $@
