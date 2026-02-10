# SEP-24: Interactive Deposit and Withdrawal

## Overview

This specification defines implementation guidance for SEP-24 interactive flows. It covers discovery, interactive session creation, transaction status retrieval, and fee calculation with SEP-38 compatibility hooks.

## Quick Reference

- Depends on: SEP-1, SEP-10, SEP-12
- Endpoints: `GET /info`, `POST /transactions/deposit/interactive`, `POST /transactions/withdraw/interactive`, `GET /transaction`, `GET /transactions`, `GET /fee`
- Authentication: SEP-10 JWT for protected routes (`403` when missing/invalid)

## Implementation Requirements

### Server MUST

- [ ] Publish SEP-24 endpoint in `stellar.toml`.
- [ ] Require SEP-10 bearer token on protected SEP-24 endpoints.
- [ ] Support `GET /info` with supported assets and fee details.
- [ ] Create transaction record and return interactive URL for deposit interactive requests.
- [ ] Create transaction record and return interactive URL for withdrawal interactive requests.
- [ ] Expose transaction lookup endpoint by `id`, `external_transaction_id`, and `stellar_transaction_id`.
- [ ] Expose transaction list endpoint scoped to authenticated account and support `asset_code`, `kind`, `limit`, and `no_older_than`.
- [ ] Use SEP-24 state values and reject invalid transitions.
- [ ] Return deterministic error payloads.

### Server MUST NOT

- [ ] Allow unauthorized access to account-specific transactions.
- [ ] Return invalid status transitions.

### Server SHOULD

- [ ] Support SEP-12 field extension patterns for KYC requirements.
- [ ] Support SEP-38 compatible fee strategy when quote service is enabled.

## Endpoint Specifications

### GET /info

Returns anchor capabilities, supported assets, and feature flags.
`fee_fixed` and `fee_percent` are numeric values.

### POST /transactions/deposit/interactive

Creates interactive deposit session and returns transaction id + URL.
`account` is optional, but if present must be a valid Stellar account and match JWT subject.

### POST /transactions/withdraw/interactive

Creates interactive withdrawal session and returns transaction id + URL.
`account` is optional, but if present must be a valid Stellar account and match JWT subject.

### GET /transaction

Returns one transaction for the authenticated account, queried by one of:
- `id`
- `external_transaction_id`
- `stellar_transaction_id`

Response transaction shape includes SEP-24 fields such as `more_info_url`, `kind` (`deposit` or `withdrawal`), and required `to`/`from` fields depending on kind.

### GET /transactions

Returns transaction list for authenticated account with support for:
- `asset_code`
- `kind` (`deposit`, `withdrawal`)
- `limit`
- `no_older_than`

Results are returned in descending `started_at` order.

### GET /fee

Returns fee for operation and asset pair as a numeric `fee` field. May integrate quote logic.

## Security Considerations

- Enforce bearer token auth on protected routes.
- Scope transaction visibility to token subject.
- Validate asset and amount inputs before transaction creation.

## Validation

```bash
npx @stellar/anchor-tests --home-domain http://localhost:8080 --seps 1 10 24
```
