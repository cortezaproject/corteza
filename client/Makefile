WEB_SUBDIRS := $(filter %/, $(wildcard $(CURDIR)/web/*/))

.PHONY: yarn

yarn:
	@for dir in $(WEB_SUBDIRS); do \
		echo "---Yarning $$dir---"; \
		(cd $$dir && yarn) || (echo "Failed to yarn $$dir"; exit 1); \
	done

.DEFAULT_GOAL := yarn