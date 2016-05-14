all : clean deps test build
.PHONY: all

LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S %Z')"
LDFLAGS += -X "main.build=$(CI_BUILD_NUMBER)"

EXECUTABLE ?= $(shell basename '$(shell pwd)')
COMMIT ?= $(or $(CI_COMMIT), $(shell git rev-parse --short HEAD))

LDFLAGS += -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

watch:
	go get github.com/onsi/ginkgo/ginkgo
	ginkgo watch -r -cover

savedeps:
	rm -rf vendor Godeps
	godep save -t ./...

clean:
	rm -rf $(EXECUTABLE)
	go clean -v -i ./...

deps:
	go get -t -v ./...

test:
	go test `go list ./... | grep -v /vendor/` -cover -ginkgo.failFast

# test:
# 	@for PKG in $(PACKAGES); do go test -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

$(EXECUTABLE): $(wildcard *.go)
	go build -o $(EXECUTABLE) -ldflags '-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

