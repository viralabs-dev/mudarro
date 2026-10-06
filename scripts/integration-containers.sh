#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário}"
engine="${2:?docker ou podman}"
root="$(mktemp -d)"
trap '"$bin" run web:down --root "$root" || true; rm -rf -- "$root"' EXIT
cat > "$root/compose.yaml" <<'YAML'
services:
  web:
    image: alpine:3.21
    command: [sh, -c, 'echo ready; sleep 300']
YAML
cat > "$root/mudarro.yaml" <<YAML
version: 1
name: Container smoke
services:
  - id: web
    dir: .
    infrastructure:
      kind: $engine
      mode: compose
      file: compose.yaml
YAML
"$bin" doctor --root "$root"
"$bin" generate --root "$root"
"$bin" run web:up --root "$root"
"$bin" run web:status --root "$root"
"$bin" run web:logs --root "$root"
"$bin" run web:restart --root "$root"
"$bin" run web:down --root "$root"
echo "$engine: OK"
