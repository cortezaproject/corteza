.PHONY: dev

dev:
	@echo "---Building libs---"
	@(cd $(CURDIR)/lib && make dev) || (echo "Failed to build libs"; exit 1)
	@echo "---Yarning clients---"
	@(cd $(CURDIR)/client && make yarn) || (echo "Failed to yarn clients"; exit 1)

.DEFAULT_GOAL := dev