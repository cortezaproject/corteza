LIB_SUBDIRS := $(filter %/, $(wildcard */))

.PHONY: build dev

build:
	@echo "---Building js---"
	@(cd $(CURDIR)/js && yarn && yarn build) || (echo "Failed to build js"; exit 1)
	@echo "---Building vue---"
	@(cd $(CURDIR)/vue && yarn && yarn build) || (echo "Failed to build vue"; exit 1)

dev: build

.DEFAULT_GOAL := dev