include Makefile.common

DOCKER_IMAGE_NAME ?= promql-guard

test-flags = -coverprofile=coverage.out # -v

style:
	@echo skip common-style

check_license:
	@echo skip common-check_license

lint:
	@echo skip common-lint

.PHONY: clean
clean:
	rm -rf vendor promql-guard

.PHONY: build
build:
	goreleaser build --snapshot --rm-dist

.PHONY: run
run:
	go run cmd/promql-guard/main.go --log.level=debug --log.format=json

.PHONY: ship
ship:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: cover
cover:	test
	${GO} tool cover -html=coverage.out
