include ../common.mk

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test -tags=wasmtime ./...
	go test -tags=wasmer ./...

bench: ## Run benchmark
	go test -tags=wasmtime -bench=. ./...
	go test -tags=wasmer -bench=. ./...
