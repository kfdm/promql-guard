GO           ?= go
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
PROMU        := $(FIRST_GOPATH)/bin/promu
PREFIX                  ?= $(shell pwd)

.PHONY: common-build
common-build:
	@echo ">> building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) build --prefix $(PREFIX)

.PHONY: run
run:
	./promql-guard --log.format=json