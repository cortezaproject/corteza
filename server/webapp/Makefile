.PHONY: all install-fresh install-packages build download clean cleaner

WGET         ?= wget
WGET_FLAGS   ?= -q
TAR          ?= tar
TAR_FLAGS    ?= -xzmok
APP          ?= unify
VERSION      ?= $(shell git describe --tags --abbrev=0)
PACKAGE       = human-webapp-$(APP)-$(VERSION).tar.gz
SOURCE_LOC   ?= ../../client/web/$(APP)

all: install-packages

install-packages: $(PACKAGE)
	$(info installing package to public/)
	@ rm -rf public
	@ mkdir -p public
	@ $(TAR) $(TAR_FLAGS) -f $(PACKAGE) -C public

download: $(PACKAGE)

$(PACKAGE):
	$(info downloading $(RELEASE_PAGE)/$(@))
	@ $(WGET) $(WGET_FLAGS) $(RELEASE_PAGE)/$(@)

install-fresh: build
	@ rm -rf public
	@ cp -r $(SOURCE_LOC)/dist public

build:
	$(info building $(APP))
	@ cd $(SOURCE_LOC) && pnpm build

clean:
	@ rm -f $(PACKAGE)

cleaner:
	@ rm -rf public human-webapp-*.tar.gz
