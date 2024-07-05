WEB_SUBDIRS := $(filter %/, $(wildcard $(CURDIR)/web/*/))

.PHONY: dev

yarn:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Yarning $$dir---"; \
		(cd $$dir && yarn) || (echo "Failed to yarn $$dir"; exit 1); \
	done

dev:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Installing and linking clients $$dir---"; \
		(cd $$dir && yarn && yarn link @cortezaproject/corteza-js && yarn link @cortezaproject/corteza-vue) || (echo "Failed to run yarn or link libs $$dir"; exit 1); \
	done

.DEFAULT_GOAL := dev