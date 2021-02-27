PROJECT_NAME := aton
VERSION := $(shell git describe --tags)
GO_LINT_CI_VERSION := v1.31.0
TIME := $(shell date +%Y-%m-%dT%T%z)
BUILD := $(shell git rev-parse --short HEAD)
DIST_FOLDER := ./dist
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
	go test -v --tags="unit" $(shell go list ./... | grep -v node)
test-integration:
	go test -v --tags="integration" $(shell go list ./... | grep -v node)
test-node-docker:
	docker run -u $(shell id -u) --rm \
	-v $(shell pwd):/home/$(shell id -u -n)/app \
	-v $(shell go env GOCACHE):$(shell go env GOCACHE) \
	-v $(shell go env GOMODCACHE):$(shell go env GOMODCACHE) \
	ghcr.io/eloylp/aton-test \
	go test -v --tags="node" ./...
# The CI build image must be prepared with deps apt-get install -y libdlib-dev libblas-dev liblapack-dev libjpeg62-turbo-dev
test-node:
	go test -v --tags="node" ./...
test-e2e:
	go test -v --tags="e2e" $(shell go list ./... | grep -v node)
proto:
	protoc -I components/proto components/proto/node.proto components/proto/system.proto --go_out=plugins=grpc:components/
	find ./components/github.com/ -type f -name "*pb.go" -exec mv {} ./components/proto \;
	rm -rf ./components/github.com
test-racy:
	go test -race -v --tags="racy" $(shell go list ./... | grep -v node)
test-bench:
	go test -v -bench=. $(shell go list ./... | grep -v node)
build: clean
	mkdir -p $(DIST_FOLDER)
	go build $(FLAGS) $(LD_FLAGS) -o $(DIST_FOLDER)/ctl ./cmd/ctl/main.go
	go build $(FLAGS) $(LD_FLAGS) -o $(DIST_FOLDER)/node ./cmd/node/main.go
	@echo "Binary outputs at $(DIST_FOLDER)"

#build-cuda: clean
#	mkdir -p $(DIST_FOLDER)
#	CGO_LDFLAGS="-L/usr/local/cuda/lib64 -lcudnn -lpthread -lcuda -lcudart -lcublas -lcurand -lcusolver" go build $(FLAGS) $(LD_FLAGS) -o $(BINARY_OUTPUT)
#	@echo "Binary output at $(BINARY_OUTPUT)"
docker-test:
	docker build --build-arg uid=$(shell id -u) \
	--build-arg uname=$(shell id -u -n) \
	-t ghcr.io/eloylp/aton-test -f Dockerfile.integration-test .
init-node:
	docker run --rm -e "NODE_LISTEN_ADDR=0.0.0.0:8082" -e "NODE_DETECTOR_MODEL_DIR=./models" \
	-u $(shell id -u) \
	-v $(shell pwd):/home/$(shell id -u -n)/app \
	-v $(shell go env GOCACHE):$(shell go env GOCACHE) \
	-v $(shell go env GOMODCACHE):$(shell go env GOMODCACHE) \
    -p "8082:8082" -v $(shell pwd):/code ghcr.io/eloylp/aton-test go run ./cmd/node/main.go
docker:
	docker build -t ghcr.io/eloylp/aton:$(BUILD) .
docker-cuda:
	docker build -t ghcr.io/eloylp/aton:$(BUILD) -f Dockerfile.gpu .
clean:
	rm -rf $(DIST_FOLDER)