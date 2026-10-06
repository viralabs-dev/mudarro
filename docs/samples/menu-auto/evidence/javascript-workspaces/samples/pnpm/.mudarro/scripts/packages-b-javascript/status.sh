#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'packages-b-javascript:status' --config 'mudarro.json' --root "$ROOT" "$@"
