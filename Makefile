all : clean test build
# .PHONY: all

LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S %Z')"
LDFLAGS += -X "main.build=$(CI_BUILD_NUMBER)"

EXECUTABLE ?= $(shell basename '$(shell pwd)')
COMMIT ?= $(or $(CI_COMMIT), $(shell git rev-parse --short HEAD))

LDFLAGS += -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

PLATFORMS=linux_amd64 linux_386 linux_arm darwin_amd64 darwin_386 freebsd_amd64 freebsd_386 windows_386 windows_amd64

FLAGS_all = GOROOT=$(GOROOT) GOPATH=$(GOPATH)
FLAGS_linux_amd64   = $(FLAGS_all) GOOS=linux   GOARCH=amd64
FLAGS_linux_386     = $(FLAGS_all) GOOS=linux   GOARCH=386
FLAGS_linux_arm     = $(FLAGS_all) GOOS=linux   GOARCH=arm   GOARM=5 # ARM5 support for Raspberry Pi
FLAGS_darwin_amd64  = $(FLAGS_all) GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0
FLAGS_darwin_386    = $(FLAGS_all) GOOS=darwin  GOARCH=386   CGO_ENABLED=0
FLAGS_freebsd_amd64  = $(FLAGS_all) GOOS=freebsd  GOARCH=amd64 CGO_ENABLED=0
FLAGS_freebsd_386    = $(FLAGS_all) GOOS=freebsd  GOARCH=386   CGO_ENABLED=0
FLAGS_windows_386   = $(FLAGS_all) GOOS=windows GOARCH=386   CGO_ENABLED=0
FLAGS_windows_amd64 = $(FLAGS_all) GOOS=windows GOARCH=amd64 CGO_ENABLED=0

EXTENSION_windows_386=.exe
EXTENSION_windows_amd64=.exe

print-%: ; @echo $*=$($*)

build-local: clean $(wildcard ../*.go)
	go build -ldflags '-s -w $(LDFLAGS)' -o $(EXECUTABLE) $(wildcard ../*.go)

out/%/.built: $(wildcard ../*.go)
	@echo -n 'Building $(EXECUTABLE)-$(subst _,-,$*)$(EXTENSION_$*) ... '
	@$(FLAGS_$*) go build -ldflags '-s -w $(LDFLAGS)' -o out/$(EXECUTABLE)-$(subst _,-,$*)$(EXTENSION_$*) $(wildcard ../*.go)
	@echo 'done'

build: clean $(foreach PLATFORM,$(PLATFORMS),out/$(PLATFORM)/.built)
.PHONY: build

watch:
	@go get github.com/onsi/ginkgo/ginkgo
	@ginkgo watch -r -cover

savedeps:
	rm -rf vendor Godeps
	godep save -t ./...

clean:
	@echo "Cleaning up..."
	@rm -rf ./out $(EXECUTABLE)
	@go clean -i ./...

deps:
	@echo "Fetching dependencies..."
	@go get -t ./...

test: deps
	@echo "Executing tests..."
	@go test $(PACKAGES) -cover

