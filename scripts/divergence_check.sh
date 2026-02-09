#!/usr/bin/env bash
set -euo pipefail

# This check enforces that baseline SEP metadata is pinned in matrices.
check_line() {
  local file="$1"
  local needle="$2"
  grep -Fq "$needle" "$file" || {
    echo "missing baseline entry in $file: $needle"
    exit 1
  }
}

check_line "specs/traceability/sep1-sep10-matrix.md" "SEP-1 | 2.7.0 | 2025-01-16"
check_line "specs/traceability/sep1-sep10-matrix.md" "SEP-10 | 3.4.1 | 2024-03-20"
check_line "specs/traceability/sep24-matrix.md" "SEP-24 | 3.7.1 | 2024-08-07"
check_line "specs/traceability/sep24-matrix.md" "SEP-12 | 1.15.0 | 2025-03-03"
check_line "specs/traceability/sep24-matrix.md" "SEP-38 | 2.5.0 | 2024-07-31"

echo "divergence baseline checks passed"
