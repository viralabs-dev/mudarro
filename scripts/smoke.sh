#!/usr/bin/env bash
set -euo pipefail
bin="${1:?caminho absoluto do binário}"
root="$(mktemp -d "${TMPDIR:-/tmp}/mudarro smoke.XXXXXX")"
trap '"$bin" run app:down --root "$root" >/dev/null 2>&1 || true; rm -rf -- "$root"' EXIT
cat > "$root/mudarro.yaml" <<'YAML'
version: 1
name: Smoke Local
services:
  - id: app
    dir: .
    infrastructure:
      kind: local
    commands:
      start:
        args: [bash, -c, 'echo ready; exec sleep 60']
        group: aplicacao
      check:
        args: [printf, '%s', 'literal $(touch injected)']
        group: qualidade
YAML
"$bin" generate --root "$root"
"$bin" generate --root "$root"
"$bin" run app:up --root "$root"
"$bin" run app:up --root "$root"
"$bin" run app:status --root "$root" | grep -E 'running|em execução'
"$bin" run app:logs --root "$root" | grep ready
"$bin" run app:check --root "$root" | grep 'literal'
test ! -e "$root/injected"
"$bin" run app:restart --root "$root"
"$bin" run app:down --root "$root"
"$bin" run app:status --root "$root" | grep -E 'stopped|parado'
printf '1\n0\n0\n' | "$bin" menu --root "$root" > "$root/menu-output"
grep -E 'infrastructure|infraestrutura' "$root/menu-output"
"$bin" scan --root "$root" --json > "$root/scan.json"
echo 'Smoke local: OK'
