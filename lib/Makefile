.PHONY: build dev

build:
	@echo "---Building js---"
	@(cd $(CURDIR)/js && yarn && yarn build) || (echo "Failed to build js"; exit 1)
	@echo "---Building vue---"
	@(cd $(CURDIR)/vue && yarn && yarn build) || (echo "Failed to build vue"; exit 1)

dev: build
	@echo "---Creating sym link for js---"
	@(cd $(CURDIR)/js && yarn link) || (echo "Failed to link js"; exit 1)
	@echo "---Building vue and creating sym link---"
	@(cd $(CURDIR)/vue && yarn link && yarn link @cortezaproject/corteza-js) || (echo "Failed to process vue lib"; exit 1)

.DEFAULT_GOAL := dev
