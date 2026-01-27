.PHONY: test test-js test-vue lint

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

.DEFAULT_GOAL := dev
