# SEP-24: Interactive Deposit and Withdrawal

## Overview

This specification defines implementation guidance for SEP-24 interactive flows. It covers discovery, interactive session creation, transaction status retrieval, and fee calculation with SEP-38 compatibility hooks.

## Quick Reference

- Depends on: SEP-1, SEP-10, SEP-12
- Endpoints: `GET /info`, `POST /transactions/deposit/interactive`, `POST /transactions/withdraw/interactive`, `GET /transaction`, `GET /transactions`, `GET /fee`
- Authentication: SEP-10 JWT for protected routes

## Implementation Requirements

### Server MUST

- [ ] Publish SEP-24 endpoint in `stellar.toml`.
- [ ] Require SEP-10 bearer token on protected SEP-24 endpoints.
- [ ] Support `GET /info` with supported assets and fee details.
- [ ] Create transaction record and return interactive URL for deposit interactive requests.
- [ ] Create transaction record and return interactive URL for withdrawal interactive requests.
- [ ] Expose transaction lookup endpoint by id.
- [ ] Expose transaction list endpoint scoped to authenticated account.
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

### POST /transactions/deposit/interactive

Creates interactive deposit session and returns transaction id + URL.

### POST /transactions/withdraw/interactive

Creates interactive withdrawal session and returns transaction id + URL.

### GET /transaction

Returns one transaction by id and authenticated account.

### GET /transactions

Returns paginated transaction list for authenticated account.

### GET /fee

Returns fee for operation and asset pair. May integrate quote logic.

## Security Considerations

- Enforce bearer token auth on protected routes.
- Scope transaction visibility to token subject.
- Validate asset and amount inputs before transaction creation.

## Validation

```bash
npx @stellar/anchor-tests --home-domain http://localhost:8080 --seps 1 10 24
```
