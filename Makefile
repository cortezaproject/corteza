.PHONY: e2e e2e-ui setup setup-agent doctor claude claude-yolo intent-check intent-hooks dev dev-all test test-lib test-client test-server lint fresh audit codegen mcp-apps tag ftag

WEB_APPS := unify

# Server codegen first (cue models/stores + rest.yaml request/handlers), then
# the lib/js API clients, so the clients are generated from the final rest.yaml.
codegen:
	@echo "---Running server codegen---"
	@(cd $(CURDIR)/server && make codegen codegen-legacy) || (echo "Failed to run server codegen"; exit 1)
	@echo "---Running lib codegen---"
	@(cd $(CURDIR)/lib && make codegen) || (echo "Failed to run lib codegen"; exit 1)

# MCP Apps views: one self-contained HTML per view, committed where the server
# embeds it (server/compose/agentic/mcpui).
mcp-apps:
	@(cd $(CURDIR)/client/web/mcp-apps && pnpm build)

dev:
	@echo "---Installing dependencies---"
	@pnpm install || (echo "Failed to install dependencies"; exit 1)

dev-all:
	@echo "---Starting server and all web apps---"
	@trap 'kill 0' EXIT; \
	(cd $(CURDIR)/server && $(MAKE) watch) & \
	$(foreach app,$(WEB_APPS),(cd $(CURDIR)/client/web/$(app) && pnpm run dev) & ) \
	(sleep 4 && $(CURDIR)/dev/agent/open-webapp.sh) & \
	wait

test: test-lib test-client test-compile test-server

test-lib:
	@echo "---Testing lib---"
	@(cd $(CURDIR)/lib && make test) || (echo "Failed to test lib"; exit 1)

test-client:
	@echo "---Testing client---"
	@(cd $(CURDIR)/client && make test) || (echo "Failed to test client"; exit 1)

test-server:
	@echo "---Testing server---"
	@(cd $(CURDIR)/server && make test) || (echo "Failed to test server"; exit 1)

# Fails when any Go package cannot BUILD its test binary.
#
# `go test ./...` prints "[build failed]" for such a package and carries on, so
# its tests silently stop running and nothing says so. federation/service sat
# like that for three weeks with ten tests dormant, nine of which passed once it
# compiled again; running this found three more packages in the same state.
#
# -run='^$' matches no test, so every test binary is compiled and none execute:
# about 35 seconds for the whole tree.
test-compile:
	@echo "---Checking every package compiles its tests---"
	@(cd $(CURDIR)/server && go test -run='^$$' ./... >/dev/null) || \
		(echo "A package cannot compile its tests — its coverage is silently off. Run: cd server && go test -run='^\$$' ./..."; exit 1)

lint:
	@echo "---Linting libs---"
	@(cd $(CURDIR)/lib && make lint) || (echo "Failed to lint libs"; exit 1)
	@echo "---Linting clients---"
	@(cd $(CURDIR)/client && make lint) || (echo "Failed to lint clients"; exit 1)

fresh:
	@echo "---Fresh install---"
	@rm -rf node_modules lib/*/node_modules client/web/*/node_modules corredor/node_modules
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

# --- Intent system (see .intent/SPEC.md) ---
intent-check:
	@node .intent/intent.mjs check && node .intent/intent.mjs coverage

intent-hooks:
	@git config core.hooksPath .intent/hooks
	@echo "git hooks path set to .intent/hooks (pre-commit intent check active)"

# Onboarding for a fresh clone. The whole procedure is dev/setup.sh; both
# targets are idempotent and everything they write is gitignored.
#
# Two phases because bootstrap.sh provisions users over the API, so it cannot
# run before the server does. `make setup` prints what to do between them.
setup:
	@dev/setup.sh

# The agent toolkit's own identities. Needs the stack up.
setup-agent:
	@dev/setup.sh --agent

# Every check `make setup` makes, writing nothing. Non-zero if something is off.
doctor:
	@dev/setup.sh --check

# Launch Claude Code with a token for this checkout's Human MCP server.
#
# Flags go after `--`: make parses a bare `--flag` as its own option and dies
# before the recipe runs, but after `--` it lands in MAKECMDGOALS instead. The
# stub below is what stops make then trying to BUILD each flag as a target —
# the same trick `tag` uses for its version argument.
#
#   make claude
#   make claude MCP=1                    # also attach this checkout's human-local MCP
#   make claude -- --dangerously-skip-permissions
#   make claude ARGS="--model opus"      # for a flag containing '=', which
#                                        # make would read as an assignment
ifeq (claude,$(firstword $(MAKECMDGOALS)))
  CLAUDE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(foreach a,$(CLAUDE_ARGS),$(eval $(a):;@:))
endif

claude:
	@HUMAN_MCP=$(MCP) dev/claude.sh $(CLAUDE_ARGS) $(ARGS)

# The one flag common enough to be worth not typing `--` for.
claude-yolo:
	@HUMAN_MCP=$(MCP) dev/claude.sh --dangerously-skip-permissions $(ARGS)

# E2E (Playwright) — needs the dev stack running and client/web/unify/.env.e2e filled.
e2e:
	@(cd $(CURDIR)/client/web/unify && npx playwright test) || (echo "e2e failed (is the dev stack running and .env.e2e filled?)"; exit 1)

e2e-ui:
	@cd $(CURDIR)/client/web/unify && npx playwright test --ui
