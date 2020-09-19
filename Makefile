PROJECT_NAME := aton
BINARY_NAME := aton
VERSION := $(shell git describe --tags)
GO_LINT_CI_VERSION := v1.31.0
TIME := $(shell date +%Y-%m-%dT%T%z)
BUILD := $(shell git rev-parse --short HEAD)
DIST_FOLDER := ./dist
BINARY_OUTPUT := $(DIST_FOLDER)/$(BINARY_NAME)
LD_FLAGS=-ldflags "-s -w \
		-X=main.Name=$(PROJECT_NAME) \
		-X=main.Version=$(VERSION) \
		-X=main.Build=$(BUILD) \
		-X=main.BuildTime=$(TIME)"
FLAGS=-trimpath -tags timetzdata

.DEFAULT_GOAL := build

lint:
	golangci-lint run -v
lint-fix:
	golangci-lint run -v --fix

all: lint test build

test: test-unit test-integration test-racy test-bench

test-unit:
	go test -v --tags="unit" ./...
test-integration:
	go test -v --tags="integration" ./...
test-e2e:
	go test -v --tags="e2e" ./...
test-racy:
	go test -race -v --tags="racy" ./...
test-bench:
	go test -v -bench=. ./...
build: clean
	mkdir -p $(DIST_FOLDER)
	go build $(FLAGS) $(LD_FLAGS) -o $(BINARY_OUTPUT)
	@echo "Binary output at $(BINARY_OUTPUT)"
docker:
	docker build -t ghcr.io/eloylp/aton:$(BUILD) .
clean:
	rm -rf $(DIST_FOLDER)