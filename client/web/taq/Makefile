.PHONY: dev build preview test lint format release release-clean upload

BUILD_FLAVOUR         ?= human
BUILD_VERSION         ?= $(shell git describe --tags --abbrev=0)
BUILD_DEST_DIR        ?= build
APP_NAME              ?= $(notdir $(CURDIR))
RELEASE_NAME           = $(BUILD_FLAVOUR)-webapp-$(APP_NAME)-$(BUILD_VERSION).tar.gz
RELEASE_PKEY          ?= .upload-rsa

dev:
	@pnpm install && pnpm dev

build:
	@pnpm build

preview:
	@pnpm preview

test:
	@pnpm test:unit

lint:
	@pnpm lint

format:
	@pnpm format

release-clean:
	rm -rf $(BUILD_DEST_DIR)

release:
	@ mkdir -p $(BUILD_DEST_DIR)
	tar -C dist -czf $(BUILD_DEST_DIR)/$(RELEASE_NAME) .

upload: $(RELEASE_PKEY)
	@ echo "put $(BUILD_DEST_DIR)/*.tar.gz" | sftp -q -o "StrictHostKeyChecking no" -i $(RELEASE_PKEY) $(RELEASE_SFTP_URI)
	@ rm -f $(RELEASE_PKEY)

$(RELEASE_PKEY):
	@ echo $(RELEASE_SFTP_KEY) | base64 -d > $@
	@ chmod 0400 $@

.DEFAULT_GOAL := dev
