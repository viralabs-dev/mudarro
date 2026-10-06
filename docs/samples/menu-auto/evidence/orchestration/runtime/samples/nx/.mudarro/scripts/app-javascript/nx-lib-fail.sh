#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'app-javascript:nx-lib-fail' --config 'mudarro.json' --root "$ROOT" "$@"
