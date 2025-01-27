WEB_SUBDIRS := $(filter %/, $(wildcard $(CURDIR)/web/*/))

.PHONY: dev test lint fresh

yarn:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Yarning $$dir---"; \
		(cd $$dir && yarn) || (echo "Failed to yarn $$dir"; exit 1); \
	done

dev:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Installing and linking clients $$dir---"; \
		(cd $$dir && yarn && yarn cdeps) || (echo "Failed to run yarn or link libs $$dir"; exit 1); \
	done

test:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Testing $$dir---"; \
		(cd $$dir && yarn test) || (echo "Failed to test $$dir"; exit 1); \
	done

lint:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Linting $$dir---"; \
		(cd $$dir && yarn lint) || (echo "Failed to lint $$dir"; exit 1); \
	done

fresh:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Fresh $$dir---"; \
		(cd $$dir && rm -rf node_modules yarn.lock && yarn && yarn cdeps) || (echo "Failed to fresh $$dir"; exit 1); \
	done

.DEFAULT_GOAL := dev