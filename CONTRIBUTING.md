# Contributing

## Principles

- SEP text is authoritative for behavior.
- Any behavior difference between implementation and SEP MUST be documented in traceability files.
- All PRs must pass CI gates for tests, spec checks, and traceability checks.

## Development

```bash
cd reference/go
go test ./...
```

## Validation

```bash
make spec-lint
make traceability-check
make compliance-sep10
make compliance-sep24
```

## Pull Requests

Each PR must include:
1. Updated spec or code.
2. Traceability updates linking SEP clause to tests.
3. Test evidence (`go test` and compliance command outputs).
