# Go Reference Implementation

This module provides a runnable SEP reference server for:
- SEP-1 (`/.well-known/stellar.toml`)
- SEP-10 (`/auth`)
- SEP-24 (`/sep24/*`)

## Run

```bash
cp .env.example .env
go run ./cmd/server
```

## Test

```bash
go test ./...
```

## Wallet CLI

```bash
go run ./cmd/wallet-cli --base-url http://localhost:8080
```
