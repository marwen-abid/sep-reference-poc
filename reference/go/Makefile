SHELL := /bin/bash
GOCACHE ?= /tmp/go-build-cache

.PHONY: test run

test:
	GOCACHE=$(GOCACHE) go test ./...

run:
	go run ./cmd/server
