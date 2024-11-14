.DEFAULT_GOAL := build
PROJECT_BIN = $(shell pwd)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
GOBIN = go
PATH := $(PROJECT_BIN):$(PATH)
GOARCH = amd64
LDFLAGS = -extldflags '-static' -w -s -buildid=
GCFLAGS = "all=-trimpath=$(shell pwd) -dwarf=false -l"
ASMFLAGS = "all=-trimpath=$(shell pwd)"
APP = gokesl
DOCKER_BUILD = -f docker-compose.test.yaml up --build --force-recreate
TEST_HOST = astra@192.168.122.106

build: build-linux .crop ## Build all

release: build-linux .crop ## Build release

docker: ## Build with docker
	docker compose up --build --force-recreate || docker-compose up --build --force-recreate


build-linux: ## Build for linux
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) \
	  $(GOBIN) build -ldflags="$(LDFLAGS)" -trimpath -gcflags=$(GCFLAGS) -asmflags=$(ASMFLAGS) \
	  -o $(PROJECT_BIN)/$(APP) cmd/$(APP)/main.go

.crop:
	strip $(PROJECT_BIN)/$(APP)
	objcopy --strip-unneeded $(PROJECT_BIN)/$(APP)

test: build
	@if [ -f /usr/bin/docker-compose ]; then \
	    docker-compose $(DOCKER_BUILD); \
	else \
	    docker compose $(DOCKER_BUILD); \
	fi

test2: build
	scp $(PROJECT_BIN)/$(APP) $(TEST_HOST):/tmp
	ssh $(TEST_HOST) 'echo astraastra | sudo -S /tmp/$(APP)'
	
help:
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

