GOPKG ?=	moul.io/zapfilter
DOCKER_IMAGE ?=	moul/zapfilter
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate:
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	go doc -all > .tmp/godoc.txt
	embedmd -w README.md
	rm -rf .tmp
