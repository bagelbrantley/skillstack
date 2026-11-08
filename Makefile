BINARY_NAME=skillstack
MODULE=github.com/bagelbrantley/skillstack
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
GIT_COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE?=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS=-ldflags "-X $(MODULE)/pkg/version.Version=$(VERSION) -X $(MODULE)/pkg/version.GitCommit=$(GIT_COMMIT) -X $(MODULE)/pkg/version.BuildDate=$(BUILD_DATE)"

.PHONY: all build test clean run

all: test build

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go
	@echo "Built bin/$(BINARY_NAME) successfully!"

test:
	go test -v ./...

clean:
	rm -rf bin/

run: build
	./bin/$(BINARY_NAME)
