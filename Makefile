################################################################################
##                             VERSION PARAMS                                 ##
################################################################################

## Docker Build Versions
DOCKER_BUILD_IMAGE = golang:1.14.6
DOCKER_BASE_IMAGE = alpine:3.12

## Tool Versions
TERRAFORM_VERSION=0.11.14
KOPS_VERSION=v1.17.1
HELM_VERSION=v2.16.9
KUBECTL_VERSION=v1.18.3

################################################################################

GO ?= $(shell command -v go 2> /dev/null)
DOCKER_IMAGE ?= mattermost/channel-spammer:0.0.4
MACHINE = $(shell uname -m)
GOFLAGS ?= $(GOFLAGS:)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
BUILD_HASH := $(shell git rev-parse HEAD)

################################################################################

BUILD_HASH = $(shell git rev-parse HEAD)

# Binaries.
TOOLS_BIN_DIR := $(abspath bin)
GO_INSTALL = ./scripts/go_install.sh

MOCKGEN_VER := v1.4.3
MOCKGEN_BIN := mockgen
MOCKGEN := $(TOOLS_BIN_DIR)/$(MOCKGEN_BIN)-$(MOCKGEN_VER)

OUTDATED_VER := master
OUTDATED_BIN := go-mod-outdated
OUTDATED_GEN := $(TOOLS_BIN_DIR)/$(OUTDATED_BIN)

GOVERALLS_VER := master
GOVERALLS_BIN := goveralls
GOVERALLS_GEN := $(TOOLS_BIN_DIR)/$(GOVERALLS_BIN)

GOLINT_VER := master
GOLINT_BIN := golint
GOLINT_GEN := $(TOOLS_BIN_DIR)/$(GOLINT_BIN)

export GO111MODULE=on

## Checks the code style, tests, builds and bundles.
all: check-style dist

## Runs govet and gofmt against all packages.
.PHONY: check-style
check-style: govet lint
	@echo Checking for style guide compliance

## Runs lint against all packages.
.PHONY: lint
lint: $(GOLINT_GEN)
	@echo Running lint
	$(GOLINT_GEN) -set_exit_status ./...
	@echo lint success

## Runs govet against all packages.
.PHONY: vet
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

## Builds and thats all :)
.PHONY: dist
dist:	build

.PHONY: build
build: ## Build
	@echo Building
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -gcflags all=-trimpath=$(PWD) -asmflags all=-trimpath=$(PWD) -a -installsuffix cgo -o build/_output/bin/main  ./

build-image:  ## Build the docker image
	@echo Building Docker Image
	docker build \
	--build-arg DOCKER_BUILD_IMAGE=$(DOCKER_BUILD_IMAGE) \
	--build-arg DOCKER_BASE_IMAGE=$(DOCKER_BASE_IMAGE) \
	. -f build/Dockerfile -t $(DOCKER_IMAGE) \
	--no-cache
