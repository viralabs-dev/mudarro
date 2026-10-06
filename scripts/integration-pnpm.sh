#!/usr/bin/env bash
set -euo pipefail
bin="$(realpath "${1:?binary}")"
command -v pnpm >/dev/null || { echo 'pnpm runtime not executed: unavailable'; exit 1; }
repo="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
root="$(mktemp -d '/tmp/mudarro pnpm.XXXXXX')"
trap '"$bin" run app-javascript:down --root "$root" >/dev/null 2>&1 || true; rm -rf -- "$root"' EXIT
mkdir "$root/tools"
ln -s "$bin" "$root/tools/mudarro"
export PATH="$root/tools:$PATH"
cp "$repo/docs/samples/menu-auto/"{package.json,app.js} "$root/"
python3 - "$root/package.json" "$(pnpm --version)" <<'PY'
import json,sys
p,version=sys.argv[1:]
c=json.load(open(p)); c['packageManager']='pnpm@'+version
open(p,'w').write(json.dumps(c))
PY
"$bin" scan --root "$root" --json
"$bin" init --root "$root" --select app-javascript:hello
"$bin" generate --root "$root"
bash "$root/.mudarro/scripts/app-javascript/test.sh" | grep 'sample quality OK'
bash "$root/.mudarro/scripts/app-javascript/hello.sh" | grep 'sample custom OK'
"$bin" run app-javascript:up --root "$root"
"$bin" run app-javascript:up --root "$root"
"$bin" run app-javascript:status --root "$root"
for attempt in $(seq 1 20); do
  if "$bin" run app-javascript:logs --root "$root" | grep 'sample ready'; then break; fi
  if [[ "$attempt" = 20 ]]; then exit 1; fi
  sleep 1
done
"$bin" run app-javascript:restart --root "$root"
"$bin" run app-javascript:down --root "$root"
echo 'pnpm real generated wrappers/supervisor: OK; no install action executed'
