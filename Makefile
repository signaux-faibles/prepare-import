.DEFAULT_GOAL := help

install: ## Install dependencies (including development/testing dependencies)
	@go get -t ./...

format: ## Fix the formatting of .go files
	@go fmt

help: ## This help.
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: help install format
