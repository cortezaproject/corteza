WEB_SUBDIRS := $(filter %/, $(wildcard web/*/))

.PHONY: yarn

yarn:
	@for dir in $(WEB_SUBDIRS); do \
		echo "Entering $$dir"; \
		cd $$dir && yarn; \
		cd ../..; \
	done

.DEFAULT_GOAL := yarn