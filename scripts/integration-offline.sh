#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário absoluto}"
root="$(mktemp -d)"
trap 'rm -rf -- "$root"' EXIT
printf '{"name":"offline-demo","scripts":{"start":"node app.js"}}\n' > "$root/package.json"
docker run --rm --network none --user "$(id -u):$(id -g)" \
  -v "$bin:/usr/local/bin/mudarro:ro" -v "$root:/project" \
  alpine:3.21 sh -ec 'mudarro scan --root /project --json; mudarro init --root /project; mudarro generate --root /project; mudarro generate --root /project'
test -f "$root/menu.sh"
echo 'Detecção e geração com rede desabilitada: OK'
