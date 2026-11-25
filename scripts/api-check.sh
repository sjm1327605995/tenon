#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

go run scripts/api_snapshot.go "$TMP/api.txt"
baseline="${TENON_API_BASELINE:-docs/public-api.txt}"
if ! diff -u "$baseline" "$TMP/api.txt"; then
    echo "Public API changed. If intentional, run 'make api-update' and document the migration." >&2
    exit 1
fi
