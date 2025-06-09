export SHELL := /bin/sh
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

.ONESHELL:

# Include .env file
include .env
export $(shell sed 's/=.*//' .env)

# Include .envrc file if it exists
ifneq (,$(wildcard .envrc))
  include .envrc
  # Commented out problematic line
  # export $(shell grep -v '^#' .envrc | sed "s/export //")
endif

default: testacc

.PHONY: help
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install-githooks
install-githooks: ## Install git hooks
	git config --local core.hooksPath .githooks

.PHONY: build
build:	## Build the provider
	go build

.PHONY: deploy-locally
deploy-locally: ## Install the provider locally in ~/.terraform.d/plugins. Optional set VERSION arg to use specific version, otherwise 0.0.1 will be used
	"$(CURDIR)/scripts/deploy_locally.sh" $(VERSION)

.PHONY: generate
generate: ## Run go generate
	go generate ./...

.PHONY: go-fmt
go-fmt: ## Run go fmt
	go fmt ./...

.PHONY: tf-fmt
tf-fmt: ## Run terraform fmt
	terraform fmt -recursive

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5

.PHONY: go-lint
go-lint: tools ## Run Golang linters
	@echo "==> Run Golang CLI linter..."
	@golangci-lint config verify
	@golangci-lint run

.PHONY: pull_test_assets
pull_test_assets: ## Pull test docker images
	@docker compose -f "$(CURDIR)/docker-compose.yaml" pull

.PHONY: start_test_env
start_test_env: ## Start test environment
	"$(CURDIR)/scripts/start_test_env.sh"
	"$(CURDIR)/scripts/wait_for_test_env_ready.sh"
	"$(CURDIR)/scripts/setup_test_env.sh"

.PHONY: test
test: ## Run acceptance tests only (no setup or cleanup)
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run acceptance tests
.PHONY: testacc
testacc: start_test_env ## Start test environment, run acceptance tests and clean up
	@trap '$(MAKE) clean' EXIT

	$(MAKE) test

.PHONY: clean
clean: ## Clean up test environment
	"$(CURDIR)/scripts/stop_test_env.sh"
