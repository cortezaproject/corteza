build:
	docker build --rm -t $(shell cat .project) .

.PHONY: build