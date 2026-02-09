#!/usr/bin/env bash
set -euo pipefail

required_files=(
  "specs/sep10/openapi.yaml"
  "specs/sep24/openapi.yaml"
  "specs/sep10/ai-spec.md"
  "specs/sep24/ai-spec.md"
  "specs/sep10/test-vectors.json"
  "specs/sep24/test-vectors.json"
)

for f in "${required_files[@]}"; do
  [[ -f "$f" ]] || { echo "missing required spec file: $f"; exit 1; }
done

for f in specs/sep10/openapi.yaml specs/sep24/openapi.yaml; do
  grep -q '^openapi: 3.1.0' "$f" || { echo "$f must declare OpenAPI 3.1.0"; exit 1; }
  grep -q '^paths:' "$f" || { echo "$f missing paths section"; exit 1; }
  grep -q '^components:' "$f" || { echo "$f missing components section"; exit 1; }
done

for f in specs/sep10/test-vectors.json specs/sep24/test-vectors.json; do
  python3 -m json.tool "$f" >/dev/null
done

echo "spec lint checks passed"
