# Copyright (c) 2022 Cisco and/or its affiliates.

include common.mk

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy-all
tidy-all:	## go mod tidy all go modules
	./scripts/for_all_go_modules.sh --with-file Makefile -- make tidy

.PHONY: license-cache-all
license-cache-all: ${REPO_ROOT}/bin/licensei
	./scripts/for_all_go_modules.sh --with-file Makefile -- make license-cache

.PHONY: license-check-all
license-check-all: ${REPO_ROOT}/bin/licensei
	./scripts/for_all_go_modules.sh --with-file Makefile -- make license-check

.PHONY: fmt-all
fmt-all:	## go fmt all go modules
	./scripts/for_all_go_modules.sh --with-file Makefile -- make fmt

.PHONY: vet-all
vet-all:	## go vet all go modules
	./scripts/for_all_go_modules.sh --with-file Makefile -- make vet

.PHONY: test-all
test-all:	## go test all go modules
	./scripts/for_all_go_modules.sh --with-file Makefile -- make test

.PHONY: build-all
build-all:	## go build all go modules
	./scripts/for_all_go_modules.sh --with-file Makefile -- make build

.PHONY: lint-all
lint-all: ${REPO_ROOT}/bin/golangci-lint ## lint the whole repo
	./scripts/for_all_go_modules.sh --parallel 1 --with-file Makefile -- make lint

.PHONY: lint-fix-all
lint-fix-all: ${REPO_ROOT}/bin/golangci-lint ## lint --fix the whole repo
	./scripts/for_all_go_modules.sh --parallel 1 --with-file Makefile -- make lint-fix

.PHONY: mod-download-all
mod-download-all:	## go mod download all go modules
	./scripts/for_all_go_modules.sh -- go mod download all
