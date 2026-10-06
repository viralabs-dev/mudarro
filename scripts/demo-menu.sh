#!/usr/bin/env bash
# Observer demo: disposable fixture, no source edits or shared containers/clusters.
set -euo pipefail
bin="$(realpath "${1:?absolute checkpoint binary}")"
repo="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
root="$(mktemp -d /tmp/mudarro-observer.XXXXXX)"
export PATH="$(dirname "$bin"):$PATH"
trap '"$bin" run app-javascript:down --root "$root" >/dev/null 2>&1 || true; rm -rf -- "$root"' EXIT
cp "$repo/docs/samples/menu-auto/"{package.json,app.js} "$root/"
"$bin" init --root "$root" --name "${MUDARRO_SAMPLE_NAME:-Aurora}" --select app-javascript:hello >/dev/null
"$bin" generate --root "$root" >/dev/null
bash "$root/menu.sh"
