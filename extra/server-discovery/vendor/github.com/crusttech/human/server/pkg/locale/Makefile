
PATH_ROOT     := ../../..
PATH_INTERNAL := $(PATH_ROOT)/locale
VERSION       := 2022.9.x

EXTERNAL_PKG  := corteza-locale-$(VERSION).zip
EXTERNAL_DIR  := corteza-locale-$(VERSION)

# Default task
copy-english-from-root-locale-dir: src/en

# Copy internal language files (english) to pkg/locale/src
#
# This way languages will be embedded into the build
src/en:
	rm -rf $*
	cp -r $(PATH_INTERNAL)/en ./src

# Catch-all for external locales
src/%: unpack-external
	cp -r $(EXTERNAL_DIR)/src/$* ./src

unpack-external: $(EXTERNAL_DIR)

$(EXTERNAL_DIR): $(EXTERNAL_PKG)
	unzip $@

$(EXTERNAL_PKG):
	curl -L https://github.com/cortezaproject/corteza-locale/archive/refs/heads/$(VERSION).zip -o $@

clean:
	rm -rf $(EXTERNAL_PKG) $(EXTERNAL_DIR)
