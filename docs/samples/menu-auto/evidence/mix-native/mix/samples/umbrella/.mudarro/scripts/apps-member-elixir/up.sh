#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../../.." && pwd)"
exec mudarro run 'apps-member-elixir:up' --config 'mudarro.json' --root "$ROOT" "$@"
