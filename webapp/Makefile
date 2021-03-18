WGET         ?= wget
WGET_FLAGS   ?= -q
TAR          ?= tar
TAR_FLAGS    ?= -xzmok
RELEASE_PAGE ?= https://releases.cortezaproject.org/files
APPS         ?= admin compose workflow
ALL_APPS     ?= $(APPS) one
VERSION      ?= 2021.3.0-rc.1
PACKAGES      = $(addprefix corteza-webapp-,$(addsuffix -$(VERSION).tar.gz,$(ALL_APPS)))

all: install

install: $(PACKAGES)
	rm -rf $(addprefix public/,$(APPS))
	mkdir -p $(addprefix public/,$(APPS))
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-one-$(VERSION).tar.gz -C public
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-admin-$(VERSION).tar.gz -C public/admin
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-compose-$(VERSION).tar.gz -C public/compose
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-workflow-$(VERSION).tar.gz -C public/workflow

download: $(PACKAGES)

$(PACKAGES):
	$(WGET) $(WGET_FLAGS) $(RELEASE_PAGE)/$(@)

clean:
	rm -f $(PACKAGES)
