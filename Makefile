.DEFAULT_GOAL := prepare-import

prepare-import: *.go
	@make install
	@go build

install: ## Install dependencies (including development/testing dependencies)
	@go get -t ./...

test: ## Run tests, including end-to-end binary tests
	@make prepare-import
	@go test

format: ## Fix the formatting of .go files
	@go fmt

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install test format help
