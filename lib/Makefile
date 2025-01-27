.PHONY: build dev test test-js test-vue lint fresh

build:
	@echo "---Building js---"
	@(cd $(CURDIR)/js && yarn && yarn build) || (echo "Failed to build js"; exit 1)
	@echo "---Building vue---"
	@(cd $(CURDIR)/vue && yarn && yarn build) || (echo "Failed to build vue"; exit 1)

dev: build
	@echo "---Creating sym link for js---"
	@(cd $(CURDIR)/js && yarn && yarn link && yarn build) || (echo "Failed to link js"; exit 1)
	@echo "---Building vue and creating sym link---"
	@(cd $(CURDIR)/vue && yarn && yarn cdeps && yarn build && yarn link) || (echo "Failed to process vue lib"; exit 1)

test: test-js test-vue

test-js:
	@echo "---Testing js---"
	@(cd $(CURDIR)/js && yarn test) || (echo "Failed to test js"; exit 1)

test-vue:
	@echo "---Testing vue---"
	@(cd $(CURDIR)/vue && yarn test) || (echo "Failed to test vue"; exit 1)

lint:
	@echo "---Linting js---"
	@(cd $(CURDIR)/js && yarn lint) || (echo "Failed to lint js"; exit 1)
	@echo "---Linting vue---"
	@(cd $(CURDIR)/vue && yarn lint) || (echo "Failed to lint vue"; exit 1)

fresh:
	@echo "---Fresh js---"
	@(cd $(CURDIR)/js && rm -rf node_modules yarn.lock && yarn && yarn link && yarn build) || (echo "Failed to fresh js"; exit 1)
	@echo "---Fresh vue---"
	@(cd $(CURDIR)/vue && rm -rf node_modules yarn.lock && yarn && yarn cdeps && yarn build && yarn link) || (echo "Failed to fresh vue"; exit 1)

.DEFAULT_GOAL := dev
