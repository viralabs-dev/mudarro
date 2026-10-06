#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'app-space-nested-go:install' --config 'mudarro.json' --root "$ROOT" "$@"
