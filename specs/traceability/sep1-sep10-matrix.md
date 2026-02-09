# SEP-1 + SEP-10 Traceability Matrix

## Normative Baseline

| SEP | Version | Last Updated | Source |
|---|---:|---|---|
| SEP-1 | 2.7.0 | 2025-01-16 | https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0001.md |
| SEP-10 | 3.4.1 | 2024-03-20 | https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0010.md |

## Requirement Mapping

| Requirement ID | SEP Clause | Requirement Summary | Implementation | Test/Check ID | Status |
|---|---|---|---|---|---|
| SEP1-001 | SEP-1 TOML field requirements | `stellar.toml` MUST expose `VERSION` | `reference/go/sep1/toml.go` | `SEP1_TOML_001` | IMPLEMENTED |
| SEP1-002 | SEP-1 service discovery | `stellar.toml` MUST expose `SIGNING_KEY` | `reference/go/sep1/toml.go` | `SEP1_TOML_002` | IMPLEMENTED |
| SEP1-003 | SEP-1 endpoint publication | `stellar.toml` MUST expose `WEB_AUTH_ENDPOINT` for SEP-10 | `reference/go/sep1/toml.go` | `SEP1_TOML_003` | IMPLEMENTED |
| SEP1-004 | SEP-1 transfer endpoint publication | `stellar.toml` MUST expose transfer server endpoint for SEP-24 | `reference/go/sep1/toml.go` | `SEP1_TOML_004` | IMPLEMENTED |
| SEP10-001 | SEP-10 challenge | Challenge transaction source MUST be server signing account | `reference/go/sep10/challenge.go` | `SEP10_CHAL_001` | IMPLEMENTED |
| SEP10-002 | SEP-10 challenge | Challenge transaction sequence MUST be 0 | `reference/go/sep10/challenge.go` | `SEP10_CHAL_002` | IMPLEMENTED |
| SEP10-003 | SEP-10 challenge | First operation MUST be `manageData` with source = client account | `reference/go/sep10/challenge.go` | `SEP10_CHAL_003` | IMPLEMENTED |
| SEP10-004 | SEP-10 challenge | Nonce MUST be 64 bytes of cryptographic randomness | `reference/go/sep10/challenge.go` | `SEP10_CHAL_004` | IMPLEMENTED |
| SEP10-005 | SEP-10 challenge | Challenge MUST be signed by server signing key | `reference/go/sep10/challenge.go` | `SEP10_CHAL_005` | IMPLEMENTED |
| SEP10-006 | SEP-10 verify | Server MUST reject challenge with invalid source account | `reference/go/sep10/verify.go` | `SEP10_VER_001` | IMPLEMENTED |
| SEP10-007 | SEP-10 verify | Server MUST reject expired challenge | `reference/go/sep10/verify.go` | `SEP10_VER_002` | IMPLEMENTED |
| SEP10-008 | SEP-10 verify | Server MUST require server signature presence/validity | `reference/go/sep10/verify.go` | `SEP10_VER_003` | IMPLEMENTED |
| SEP10-009 | SEP-10 verify | Server MUST require client signature for authenticated account | `reference/go/sep10/verify.go` | `SEP10_VER_004` | IMPLEMENTED |
| SEP10-010 | SEP-10 token | Server SHOULD return JWT only after successful verification | `reference/go/sep10/jwt.go` | `SEP10_JWT_001` | IMPLEMENTED |
| SEP10-011 | SEP-10 API | `GET /auth` challenge endpoint response shape | `reference/go/sep10/handler.go` | `SEP10_API_001` | IMPLEMENTED |
| SEP10-012 | SEP-10 API | `POST /auth` token endpoint response shape | `reference/go/sep10/handler.go` | `SEP10_API_002` | IMPLEMENTED |

## Verification Commands

```bash
make spec-lint
make traceability-check
cd reference/go && go test ./sep10/... ./sep1/...
```
