.PHONY: setup-hooks

setup-hooks:
	@ echo "Setting up precommit hooks..."

	@ cp ./.githooks/pre-commit .git/hooks/pre-commit
	@ chmod +x .git/hooks/pre-commit

	@ echo "Hooks setup complete. You're all set!"
