# SEP-24 Traceability Matrix

## Normative Baseline

| SEP | Version | Last Updated | Source |
|---|---:|---|---|
| SEP-24 | 3.7.1 | 2024-08-07 | https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0024.md |
| SEP-12 | 1.15.0 | 2025-03-03 | https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0012.md |
| SEP-38 | 2.5.0 | 2024-07-31 | https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0038.md |

## Requirement Mapping

| Requirement ID | SEP Clause | Requirement Summary | Implementation | Test/Check ID | Status |
|---|---|---|---|---|---|
| SEP24-001 | SEP-24 Authentication | Protected endpoints MUST require SEP-10 bearer token | `reference/go/internal/middleware/auth.go` | `SEP24_AUTH_001` | IMPLEMENTED |
| SEP24-002 | SEP-24 `/info` | Anchor MUST provide `/info` with supported assets and features | `reference/go/sep24/info.go` | `SEP24_INFO_001` | IMPLEMENTED |
| SEP24-003 | SEP-24 interactive deposit | Anchor MUST create transaction record for interactive deposit session | `reference/go/sep24/deposit.go` | `SEP24_DEP_001` | IMPLEMENTED |
| SEP24-004 | SEP-24 interactive withdrawal | Anchor MUST create transaction record for interactive withdraw session | `reference/go/sep24/withdraw.go` | `SEP24_WDR_001` | IMPLEMENTED |
| SEP24-005 | SEP-24 transaction query | Anchor MUST expose transaction lookup by id | `reference/go/sep24/transaction.go` | `SEP24_TX_001` | IMPLEMENTED |
| SEP24-006 | SEP-24 transactions list | Anchor MUST expose paginated transaction listing | `reference/go/sep24/transaction.go` | `SEP24_TX_002` | IMPLEMENTED |
| SEP24-007 | SEP-24 status model | Transaction status values MUST follow SEP-24 states | `reference/go/sep24/state.go` | `SEP24_STATE_001` | IMPLEMENTED |
| SEP24-008 | SEP-24 status transitions | Invalid state transitions MUST be rejected | `reference/go/sep24/state.go` | `SEP24_STATE_002` | IMPLEMENTED |
| SEP24-009 | SEP-24 fees | Anchor MUST expose `/fee` for interactive flow fee discovery | `reference/go/sep24/transaction.go` | `SEP24_FEE_001` | IMPLEMENTED |
| SEP24-010 | SEP-24 + SEP-38 alignment | Fee response SHOULD align with quote-based logic when enabled | `reference/go/sep24/transaction.go` | `SEP24_FEE_002` | IMPLEMENTED |
| SEP24-011 | SEP-24 + SEP-12 fields | KYC related fields MUST be represented in interactive flow model | `reference/go/sep24/deposit.go` | `SEP24_KYC_001` | IMPLEMENTED |
| SEP24-012 | SEP-24 API | `POST /transactions/*/interactive` MUST return interactive URL and id | `reference/go/sep24/handler.go` | `SEP24_API_001` | IMPLEMENTED |

## Verification Commands

```bash
make spec-lint
make traceability-check
cd reference/go && go test ./sep24/... ./internal/...
```
