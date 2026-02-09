#!/usr/bin/env bash
set -euo pipefail

if ! command -v npx >/dev/null 2>&1; then
  echo "npx is required to run @stellar/anchor-tests"
  exit 1
fi

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <sep-list>"
  echo "example: $0 '1 10 24'"
  exit 1
fi

SEPS="$1"
HOME_DOMAIN="${HOME_DOMAIN:-http://localhost:8080}"

npx --yes @stellar/anchor-tests --home-domain "$HOME_DOMAIN" --seps $SEPS --verbose
