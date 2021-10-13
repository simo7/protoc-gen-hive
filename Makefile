SHELL := /usr/bin/env bash -o pipefail
.DEFAULT_GOAL := help

RED=$(shell tput -T xterm setaf 1)
GREEN=$(shell tput -T xterm setaf 2)
YELLOW=$(shell tput -T xterm setaf 3)
RESET=$(shell tput -T xterm sgr0)

export GO_VERSION=1.17
export PROTOC_VERSION=3.18.0
export PROTOC_GEN_GO_VERSION=1.27.1

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: generate
generate: ## Build protobuf stubs in docker container
generate:
	docker run -it --rm \
		-v $$(pwd):/src \
		-w /src \
		-e PROTOC_VERSION \
		-e PROTOC_GEN_GO_VERSION \
		golang:${GO_VERSION} \
		make build

.PHONY: build
build: ## Build
build: check-inside-container clean dependencies
	protoc \
    --go_out=./ \
    --go_opt=paths=source_relative \
		hive_options/hive_options.proto
	go build -o /go/bin/protoc-gen-hive


.PHONY: dependencies
dependencies: ## Get remote proto dependencies
dependencies: check-inside-container
	apt update && apt install -y clang-format unzip
	wget -P /tmp https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip
	unzip -o -d /usr/local/ /tmp/protoc-3.18.0-linux-x86_64.zip
	wget -P /tmp https://github.com/protocolbuffers/protobuf-go/releases/download/v${PROTOC_GEN_GO_VERSION}/protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.linux.amd64.tar.gz
	tar -xzvf /tmp/protoc-gen-go.v${PROTOC_GEN_GO_VERSION}.linux.amd64.tar.gz -C /usr/local/bin

.PHONY: check-inside-container
check-inside-container:
ifeq (,$(wildcard /.dockerenv))
	@echo ""
	@echo "${YELLOW} 👉 Only use 'make build' for testing, use 'generate' otherwise 👈"
	@echo "${RESET}"
endif

PHONY: clean
clean: ## Delete all stub files
	rm -rf ./hive_options/*.pb.go