.PHONY: dep test build release upload

YARN_FLAGS            ?= --non-interactive --no-progress --silent --emoji false
YARN                   = yarn $(YARN_FLAGS)

REPO_NAME             ?= corteza-webapp-privacy

BUILD_FLAVOUR         ?= corteza
BUILD_FLAGS           ?= --production
BUILD_DEST_DIR         = dist
BUILD_TIME            ?= $(shell date +%FT%T%z)
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_NAME             = $(REPO_NAME)-$(BUILD_VERSION)

RELEASE_NAME           = $(BUILD_NAME).tar.gz
RELEASE_EXTRA_FILES   ?= README.md LICENSE CONTRIBUTING.md DCO
RELEASE_PKEY          ?= .upload-rsa

dep:
	$(YARN) install

test:
	$(YARN) lint
	$(YARN) test:unit

build:
	export BUILD_VERSION=${BUILD_VERSION} && $(YARN) build $(BUILD_FLAGS)

release:
	@ cp $(RELEASE_EXTRA_FILES) $(BUILD_DEST_DIR)
	@ tar -C $(BUILD_DEST_DIR) -czf $(RELEASE_NAME) $(dir $(BUILD_DEST_DIR))

upload: $(RELEASE_PKEY)
	@ echo "put *.tar.gz" | sftp -q -o "StrictHostKeyChecking no" -i $(RELEASE_PKEY) $(RELEASE_SFTP_URI)
	@ rm -f $(RELEASE_PKEY)

$(RELEASE_PKEY):
	@ echo $(RELEASE_SFTP_KEY) | base64 -d > $(RELEASE_PKEY)
	@ chmod 0400 $@
