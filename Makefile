SHELL := /bin/bash
GOCACHE ?= /tmp/go-build-cache

.PHONY: spec-lint traceability-check divergence-check test compliance-sep10 compliance-sep24 compliance-all

spec-lint:
	./scripts/spec_lint.sh

traceability-check:
	./scripts/traceability_check.sh

divergence-check:
	./scripts/divergence_check.sh

test:
	cd reference/go && GOCACHE=$(GOCACHE) go test ./...

compliance-sep10:
	./scripts/run_anchor_tests.sh "1 10"

compliance-sep24:
	./scripts/run_anchor_tests.sh "1 10 24"

compliance-all: spec-lint traceability-check divergence-check test compliance-sep24
