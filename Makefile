-include .env
export

# disables built-in rules
.SUFFIXES: 

###############################################################################
#
# Dependencies
#
###############################################################################

GOLANGCI-LINT_VERSION := v1.52.2
NATS_VERSION := v2.9.21
OAPICODEGEN_VERSION := v1.12.4

GOPATH := $(shell go env GOPATH)

golangci-lint := $(GOPATH)/bin/golangci-lint
oapi-codegen := $(GOPATH)/bin/oapi-codegen

$(oapi-codegen):
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$(OAPICODEGEN_VERSION)

$(golangci-lint):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI-LINT_VERSION)

.PHONY: install	
install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI-LINT_VERSION)

###############################################################################
#
# Configure local environment
#
###############################################################################

	
###############################################################################
#
# Docker compose commands
#
###############################################################################

.PHONY: clean-all
clean-all: down-natsbox
	docker compose down --volumes

.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: up-natsbox
up-natsbox:
	docker compose -f docker-natsbox.yaml up -d 

.PHONY: down-natsbox
down-natsbox:
	docker compose -f docker-natsbox.yaml down 

.PHONY: run-world
run-world:
	NATS_URL=$(NATS_URL) NATS_PASSWORD=$(NATS_PASSWORD) NATS_USER=$(NATS_USER) go run cmd/*.go

.PHONY: init
init:
	NATS_URL=$(NATS_URL) NATS_PASSWORD=$(NATS_PASSWORD) NATS_USER=$(NATS_USER) go run cmd/init/*.go

.PHONY: natsbox
natsbox: up-natsbox
	docker exec -it greenisland-nats-box-1 /bin/sh

###############################################################################
#
# Build commands 
#
###############################################################################


###############################################################################
#
# Linting & testing
#
###############################################################################

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: test
test:
	go test ./... -timeout 30s -failfast
