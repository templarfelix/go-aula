PWD = $(shell pwd -L)
GOCMD=go

.PHONY: all

all: help

start: ## start
	$(GOCMD) run ./cmd/api/main.go

update: ## update libs
	$(GOCMD) get -u ./...

tidy: ## Downloads go dependencies
	$(GOCMD) mod tidy

fmt: tidy ## Run go fmt
	go fmt ./...


test: ## Run go fmt
	go test ./...

e2e: ## Run go fmt
	go test ./... -tags e2e

help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help \n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'