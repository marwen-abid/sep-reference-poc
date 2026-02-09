# SEP-{N}: {Title}

## Overview

{2-3 sentence summary}

## Quick Reference

- Depends on: {SEP list}
- Endpoints: {endpoint list}
- Authentication: {mode}

## Implementation Requirements

### Server MUST
- [ ] {testable requirement}

### Server MUST NOT
- [ ] {forbidden behavior}

### Server SHOULD
- [ ] {recommended behavior}

## Endpoint Specifications

### {METHOD} {PATH}

**Purpose**: {purpose}

**Request**:

| Parameter | Type | Required | Description |
|---|---|---|---|
| {name} | {type} | {yes/no} | {description} |

**Response (200)**

```json
{}
```

**Error Response**

```json
{"error":"..."}
```

## Security Considerations

{security notes}

## Validation

```bash
npx @stellar/anchor-tests --home-domain http://localhost:8080 --seps {N}
```
