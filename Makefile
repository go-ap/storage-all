SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO ?= go
TEST := $(GO) test
TEST_FLAGS ?= -v
TEST_TARGET ?= .
GO111MODULE = on
PROJECT_NAME := $(shell basename $(PWD))

.PHONY: test coverage clean download

download: go.sum

go.sum: go.mod
	$(GO) mod tidy

test: go.sum clean
	@
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_fs" -cover $(TEST_TARGET) -json > tests.json
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_sqlite" -cover $(TEST_TARGET) -json >> tests.json
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_boltdb" -cover $(TEST_TARGET) -json >> tests.json
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_badger" -cover $(TEST_TARGET) -json >> tests.json
	$(GO) tool tparse -file tests.json

bench: go.sum clean
	$(TEST) $(TEST_FLAGS) -test.bench=. -test.run=xxxxx -cover $(TEST_TARGET)

coverage: go.sum clean
	@mkdir ./_coverage
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_fs" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_sqlite" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_boltdb" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(TEST) $(TEST_FLAGS) -tags "conformance storage_badger" -covermode=count -args -test.gocoverdir="$(PWD)/_coverage" $(TEST_TARGET) > /dev/null || true
	$(GO) tool covdata percent -i=./_coverage/ -o $(PROJECT_NAME).coverprofile
	@$(RM) -r ./_coverage

clean:
	@$(RM) -r ./_coverage
	$(RM) -v *.coverprofile
	@$(RM) -v tests.json

