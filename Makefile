.PHONY: dev dev-all test test-lib test-client test-server lint fresh audit codegen tag ftag

WEB_APPS := admin agentic compose one taq home workflow chatbot

codegen:
	@echo "---Running codegen---"
	@(cd $(CURDIR)/lib && make codegen) || (echo "Failed to run codegen"; exit 1)

dev:
	@echo "---Installing dependencies---"
	@pnpm install || (echo "Failed to install dependencies"; exit 1)

dev-all:
	@echo "---Starting server and all web apps---"
	@trap 'kill 0' EXIT; \
	(cd $(CURDIR)/server && $(MAKE) watch) & \
	$(foreach app,$(WEB_APPS),(cd $(CURDIR)/client/web/$(app) && pnpm run dev) & ) \
	(sleep 4 && ss -tlnp | awk '/127\.0\.0\.1:51/{split($$4,a,":");print a[2]}' | sort | while read port; do cmd.exe /c start "http://localhost:$$port" </dev/null; sleep 2; done) & \
	wait

test: test-lib test-client test-server

test-lib:
	@echo "---Testing lib---"
	@(cd $(CURDIR)/lib && make test) || (echo "Failed to test lib"; exit 1)

test-client:
	@echo "---Testing client---"
	@(cd $(CURDIR)/client && make test) || (echo "Failed to test client"; exit 1)

test-server:
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


# Usage: make tag <version>       — tag and push
#        make ftag <version>      — force tag and force push
ifeq (tag,$(firstword $(MAKECMDGOALS)))
  TAG_NAME := $(wordlist 2,2,$(MAKECMDGOALS))
  $(eval $(TAG_NAME):;@:)
endif
ifeq (ftag,$(firstword $(MAKECMDGOALS)))
  TAG_NAME := $(wordlist 2,2,$(MAKECMDGOALS))
  $(eval $(TAG_NAME):;@:)
endif

tag:
ifeq ($(TAG_NAME),)
	$(error Usage: make tag <version>, e.g. make tag 2026.3.1)
endif
	@echo "---Tagging $(TAG_NAME)---"
	git commit -m "$(TAG_NAME)"
	git push
	git tag $(TAG_NAME)
	git push origin $(TAG_NAME)

ftag:
ifeq ($(TAG_NAME),)
	$(error Usage: make ftag <version>, e.g. make ftag 2026.3.1)
endif
	@echo "---Force tagging $(TAG_NAME)---"
	git tag -f $(TAG_NAME)
	git push -f origin $(TAG_NAME)

.DEFAULT_GOAL := dev
