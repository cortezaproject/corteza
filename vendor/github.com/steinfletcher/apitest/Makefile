.PHONY: test test-examples docs fmt vet

test:
	go test ./... -v -covermode=atomic -coverprofile=coverage.out

test-examples:
	cd examples && go test -v ./... && \
	cd sequence-diagrams-with-sqlite-database && make test && cd ..

fmt:
	bash -c 'diff -u <(echo -n) <(gofmt -s -d ./)'

vet:
	bash -c 'diff -u <(echo -n) <(go vet ./...)'

test-all: fmt vet test test-examples

docs:
	cd docs && hugo server -w && cd -
