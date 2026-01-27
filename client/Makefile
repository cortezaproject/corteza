WEB_SUBDIRS := $(filter %/, $(wildcard $(CURDIR)/web/*/))

.PHONY: test lint build

test:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Testing $$dir---"; \
		(cd $$dir && yarn test:unit) || (echo "Failed to test $$dir"; exit 1); \
	done

lint:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Linting $$dir---"; \
		(cd $$dir && yarn lint) || (echo "Failed to lint $$dir"; exit 1); \
	done

build:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Building $$dir---"; \
		(cd $$dir && yarn build) || (echo "Failed to build $$dir"; exit 1); \
	done

.DEFAULT_GOAL := dev