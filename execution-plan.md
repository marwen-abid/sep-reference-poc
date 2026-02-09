# SEP Reference Execution Backlog (Phases 1-2)

This file operationalizes the approved Phase 1-2 plan into ticket tracking.

## Ticket Status Legend
- TODO: not started
- IN_PROGRESS: active
- DONE: completed and verified

## Phase 1

| ID | Title | Status | Evidence |
|---|---|---|---|
| P1-001 | Repo scaffolding and CI skeleton | DONE | `.github/workflows/ci.yml`, `Makefile`, `scripts/*` |
| P1-002 | SEP-1 + SEP-10 traceability matrix | DONE | `specs/traceability/sep1-sep10-matrix.md` |
| P1-003 | SEP-1 `stellar.toml` module | DONE | `reference/go/sep1/toml.go` |
| P1-004 | SEP-10 OpenAPI 3.1 spec | DONE | `specs/sep10/openapi.yaml` |
| P1-005 | SEP-10 AI spec with RFC2119 language | DONE | `specs/sep10/ai-spec.md` |
| P1-006 | SEP-10 test vectors | DONE | `specs/sep10/test-vectors.json` |
| P1-007 | SEP-10 Go auth implementation | DONE | `reference/go/sep10/*` |
| P1-008 | Runnable server bootstrap | DONE | `reference/go/cmd/server/main.go` |
| P1-009 | Local compose + anchor-tests harness | DONE | `reference/go/docker-compose.yml`, `scripts/run_anchor_tests.sh` |
| P1-010 | Compliance audit and remediation tracking | DONE | `compliance-report.md` |

## Phase 2

| ID | Title | Status | Evidence |
|---|---|---|---|
| P2-001 | SEP-24 traceability matrix | DONE | `specs/traceability/sep24-matrix.md` |
| P2-002 | SEP-24 AI spec + state machine | DONE | `specs/sep24/ai-spec.md`, `specs/sep24/state-machine.md` |
| P2-003 | SEP-24 OpenAPI 3.1 spec | DONE | `specs/sep24/openapi.yaml` |
| P2-004 | SEP-24 test vectors | DONE | `specs/sep24/test-vectors.json` |
| P2-005 | SEP-24 Go handlers + state machine | DONE | `reference/go/sep24/*` |
| P2-006 | Storage abstraction + in-memory store | DONE | `reference/go/internal/db/*` |
| P2-007 | Minimal interactive templates | DONE | `reference/go/sep24/interactive/*` |
| P2-008 | Full compliance integration + audit closure | DONE | `compliance-report.md`, CI workflow |
