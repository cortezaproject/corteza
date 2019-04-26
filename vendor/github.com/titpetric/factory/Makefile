.PHONY: all test

all:
	drone exec

test:
	gotest ./prof* -v -cover
	gotest ./sonyflake* -v -cover
	gotest ./database* ./profiler.go -v -cover
	gotest ./resputil -v -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
