#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'app-space-nested-go:up' --config 'mudarro.json' --root "$ROOT" "$@"
