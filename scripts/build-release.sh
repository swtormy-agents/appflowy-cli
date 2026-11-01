#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
go build -o appflowy ./cmd/appflowy
echo "built ./appflowy"
