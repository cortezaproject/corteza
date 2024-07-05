.PHONY: dev

dev:
	@echo "---Processing libs---"
	@(cd $(CURDIR)/lib && make dev) || (echo "Failed to build libs"; exit 1)
	@echo "---Processing clients---"
	@(cd $(CURDIR)/client && make dev) || (echo "Failed to yarn clients"; exit 1)

.DEFAULT_GOAL := dev