#!/usr/bin/env bash
set -euo pipefail

matrices=(
  "specs/traceability/sep1-sep10-matrix.md"
  "specs/traceability/sep24-matrix.md"
)

for f in "${matrices[@]}"; do
  [[ -f "$f" ]] || { echo "missing traceability matrix: $f"; exit 1; }
  grep -q 'Requirement ID' "$f" || { echo "$f missing Requirement ID column"; exit 1; }
  grep -q 'Status' "$f" || { echo "$f missing Status column"; exit 1; }
  if grep -q '| OPEN |' "$f"; then
    echo "traceability matrix has OPEN items: $f"
    exit 1
  fi
done

echo "traceability checks passed"
