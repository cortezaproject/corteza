.PHONY: all install-fresh install-packages build clean cleaner

WGET         ?= wget
WGET_FLAGS   ?= -q
TAR          ?= tar
TAR_FLAGS    ?= -xzmok
RELEASE_PAGE ?= https://releases.cortezaproject.org/files
APPS         ?= admin compose workflow reporter discovery
ALL_APPS     ?= $(APPS) one
VERSION      ?= $(shell git describe --tags --abbrev=0)
PACKAGES      = $(addprefix corteza-webapp-,$(addsuffix -$(VERSION).tar.gz,$(ALL_APPS)))
SOURCE_LOC   ?= "../../webapp-"
SOURCES       = $(addprefix $(SOURCE_LOC),$(ALL_APPS))

all: install-packages

install-packages: $(PACKAGES)
	$(info installing packages to public/)
	@ rm -rf $(addprefix public/,$(APPS))
	@ mkdir -p $(addprefix public/,$(APPS))
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-one-$(VERSION).tar.gz -C public
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-admin-$(VERSION).tar.gz -C public/admin
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-compose-$(VERSION).tar.gz -C public/compose
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-workflow-$(VERSION).tar.gz -C public/workflow
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-reporter-$(VERSION).tar.gz -C public/reporter
	@ $(TAR) $(TAR_FLAGS) -f corteza-webapp-discovery-$(VERSION).tar.gz -C public/discovery

download: $(PACKAGES)

$(PACKAGES):
	$(info downloading $(RELEASE_PAGE)/$(@))
	@ $(WGET) $(WGET_FLAGS) $(RELEASE_PAGE)/$(@)

install-fresh: build
	@ rm -rf public
	@ cp -r $(SOURCE_LOC)one/dist       public/
	@ cp -r $(SOURCE_LOC)admin/dist     public/admin
	@ cp -r $(SOURCE_LOC)compose/dist   public/compose
	@ cp -r $(SOURCE_LOC)workflow/dist  public/workflow
	@ cp -r $(SOURCE_LOC)reporter/dist  public/reporter
	@ cp -r $(SOURCE_LOC)discovery/dist public/discovery

build: $(SOURCES)

$(SOURCES):
	$(info building $(@))
	@ cd $(@) && yarn build

clean:
	@ rm -f $(PACKAGES)

cleaner:
	@ rm -rf public corteza-webapp-*.tar.gz
