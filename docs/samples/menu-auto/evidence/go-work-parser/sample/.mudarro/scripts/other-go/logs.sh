#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'other-go:logs' --config 'mudarro.json' --root "$ROOT" "$@"
