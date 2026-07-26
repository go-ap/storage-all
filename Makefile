SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO ?= go
TEST := $(GO) test
TEST_FLAGS ?= -timeout 20s -count=1 -tags '$(TAGS)'
TEST_TARGET ?= .
GO111MODULE = on
PROJECT_NAME := $(shell basename $(PWD))

PODMAN_SOCKET ?= /run/user/$(shell id -u)/podman/podman.sock

.PHONY: test coverage clean download

download: go.sum

go.sum: go.mod
	$(GO) mod tidy

test: go.sum clean
	@export DOCKER_HOST="unix://$(PODMAN_SOCKET)"
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_sqlite" -cover $(TEST_TARGET) -json >> tests.json || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_boltdb" -cover $(TEST_TARGET) -json >> tests.json || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_badger" -cover $(TEST_TARGET) -json >> tests.json || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_postgres" -cover $(TEST_TARGET) -json >> tests.json || true
	$(GO) tool tparse -file tests.json

bench: go.sum clean
	$(TEST) $(TEST_FLAGS) -test.bench=. -test.run=xxxxx -cover $(TEST_TARGET)

coverage: go.sum clean
	@export DOCKER_HOST="unix://$(PODMAN_SOCKET)"
	mkdir ./_coverage
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_fs" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_sqlite" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_boltdb" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_badger" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -test.bench=xxxxx -tags "conformance storage_postgres" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(GO) tool covdata percent -i=./_coverage/ -o $(PROJECT_NAME).coverprofile

clean:
	@$(RM) -r ./_coverage
	$(RM) -v *.coverprofile
	@$(RM) -v tests.json

