.PHONY: dev test lint fresh audit codegen

codegen:
	@echo "---Running codegen---"
	@(cd $(CURDIR)/lib && make codegen) || (echo "Failed to run codegen"; exit 1)

dev:
	@echo "---Installing dependencies---"
	@pnpm install || (echo "Failed to install dependencies"; exit 1)

test:
	@echo "---Testing libs---"
	@(cd $(CURDIR)/lib && make test) || (echo "Failed to test libs"; exit 1)
	@echo "---Testing clients---"
	@(cd $(CURDIR)/client && make test) || (echo "Failed to test clients"; exit 1)
	@echo "---Testing server---"
	@(cd $(CURDIR)/server && make test) || (echo "Failed to test server"; exit 1)

lint:
	@echo "---Linting libs---"
	@(cd $(CURDIR)/lib && make lint) || (echo "Failed to lint libs"; exit 1)
	@echo "---Linting clients---"
	@(cd $(CURDIR)/client && make lint) || (echo "Failed to lint clients"; exit 1)

fresh:
	@echo "---Fresh install---"
	@rm -rf node_modules lib/*/node_modules client/web/*/node_modules
	@pnpm install

audit:
	@echo "---Audit dependencies---"
	@pnpm audit

.DEFAULT_GOAL := dev
