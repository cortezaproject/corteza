WGET         ?= wget
WGET_FLAGS   ?= -q
TAR          ?= tar
TAR_FLAGS    ?= -xzmok
RELEASE_PAGE ?= https://releases.cortezaproject.org/files
APPS         ?= auth admin compose
ALL_APPS     ?= $(APPS) one
VERSION      ?= unstable
PACKAGES      = $(addprefix corteza-webapp-,$(addsuffix -$(VERSION).tar.gz,$(ALL_APPS)))

all: install

install: $(PACKAGES)
	rm -rf $(addprefix public/,$(APPS))
	mkdir -p $(addprefix public/,$(APPS))
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-one-$(VERSION).tar.gz -C public
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-auth-$(VERSION).tar.gz -C public/auth
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-admin-$(VERSION).tar.gz -C public/admin
	$(TAR) $(TAR_FLAGS) -f corteza-webapp-compose-$(VERSION).tar.gz -C public/compose

download: $(PACKAGES)

$(PACKAGES):
	$(WGET) $(WGET_FLAGS) $(RELEASE_PAGE)/$(@)

clean:
	rm -f $(PACKAGES)
