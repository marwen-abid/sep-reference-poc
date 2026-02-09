# SEP-10: Stellar Web Authentication

## Overview

This document defines an implementation-focused specification for SEP-10 authentication in anchor services. It is derived from the authoritative SEP-10 text and is optimized for both human implementation and LLM-assisted generation.

## Quick Reference

- Depends on: SEP-1
- Endpoints: `GET /auth`, `POST /auth`
- Authentication: None for challenge request; signed challenge for token issuance

## Implementation Requirements

### Server MUST

- [ ] Publish `WEB_AUTH_ENDPOINT` in `stellar.toml`.
- [ ] Build challenge transactions with source account set to server signing key.
- [ ] Set challenge sequence number to `0`.
- [ ] Include first `manageData` operation with source set to client account.
- [ ] Generate 64-byte cryptographically secure nonce in first `manageData` value.
- [ ] Sign challenge transaction with server signing key.
- [ ] Verify challenge source account, time bounds, operations, and signatures.
- [ ] Return JWT only after successful verification.

### Server MUST NOT

- [ ] Accept expired challenges.
- [ ] Accept challenges with missing or invalid server signatures.
- [ ] Issue JWT for account that did not sign challenge.

### Server SHOULD

- [ ] Use short challenge validity windows (for example, 5 minutes).
- [ ] Use short JWT TTL and rotate signing keys.

## Endpoint Specifications

### GET /auth

**Purpose**: Create challenge transaction for account authentication.

**Authentication**: None

**Request Query Parameters**:

| Parameter | Type | Required | Description |
|---|---|---|---|
| `account` | string | yes | Stellar public key (`G...`) |
| `home_domain` | string | no | Optional home domain override |
| `client_domain` | string | no | Optional client domain |
| `memo` | string | no | Optional memo for shared accounts |

**Response (200)**

```json
{
  "transaction": "AAAA...",
  "network_passphrase": "Test SDF Network ; September 2015"
}
```

### POST /auth

**Purpose**: Verify signed challenge and issue JWT.

**Authentication**: Signed SEP-10 challenge in request body

**Request Body**

```json
{
  "transaction": "AAAA..."
}
```

**Response (200)**

```json
{
  "token": "eyJ..."
}
```

**Error Response (400/401)**

```json
{
  "error": "invalid challenge"
}
```

## Security Considerations

- Enforce strict validation for source account, operation ordering, timebounds, and signatures.
- Reject replay attempts using nonce and short expiration.
- Use dedicated JWT signing secret and rotate regularly.

## Common Implementation Mistakes

1. **Incorrect challenge source account**: source must be server signing key.
2. **Weak nonce**: nonce must be 64 random bytes.
3. **Missing signature checks**: both server and client signature requirements must be enforced.

## Validation

```bash
npx @stellar/anchor-tests --home-domain http://localhost:8080 --seps 1 10
```
