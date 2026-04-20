.PHONY: test test-js test-vue lint codegen

codegen:
	@echo "---Running codegen---"
	@(cd $(CURDIR)/js && pnpm codegen) || (echo "Failed to run codegen"; exit 1)

test: test-js test-vue

test-js:
	@echo "---Testing js---"
	@(cd $(CURDIR)/js && pnpm test) || (echo "Failed to test js"; exit 1)

test-vue:
	@echo "---Testing vue---"
	@(cd $(CURDIR)/vue && pnpm test:unit) || (echo "Failed to test vue"; exit 1)

lint:
	@echo "---Linting js---"
	@(cd $(CURDIR)/js && pnpm lint) || (echo "Failed to lint js"; exit 1)
	@echo "---Linting vue---"
	@(cd $(CURDIR)/vue && pnpm lint) || (echo "Failed to lint vue"; exit 1)

.DEFAULT_GOAL := dev
