#!/usr/bin/env bash
# Real CLI UI/configuration regression in an exclusively owned disposable fixture.
set -euo pipefail
bin="$(realpath "${1:?absolute built binary}")"
root="$(mktemp -d '/tmp/mudarro ui.XXXXXX')"
config='selected config.json'
export PATH="$(dirname "$bin"):$PATH"
cleanup() {
  "$bin" run app:down --root "$root" --config "$config" >/dev/null 2>&1 || true
  rm -rf -- "$root"
}
trap cleanup EXIT
cat > "$root/$config" <<'JSON'
{"version":1,"name":"UI regression","ui":{"locale":"en","preview":{"enabled":true,"mouse":"off"}},"services":[{"id":"app","dir":".","language":"bash","manager":"bash","infrastructure":{"kind":"local"},"commands":{"start":{"group":"application","args":["sh","-c","printf 'supervisor ready\n'; exec sleep 30"]},"probe":{"group":"quality","args":["sh","-c","touch preview-only-marker"]},"secret":{"group":"quality","args":["sh","-c","API_TOKEN='fixture private words' printf safe"]}}}]}
JSON
printf 'version: 1\nname: Conventional YAML\nservices: []\n' > "$root/mudarro.yaml"
printf '{"version":1,"name":"Conventional JSON","services":[]}' > "$root/mudarro.json"
reject() { if "$@"; then printf 'Expected rejection\n' >&2; exit 1; fi; }
reject "$bin" menu --root "$root" </dev/null
"$bin" generate --root "$root" --config "$config"
# Both launchers must preserve the selected file despite ambiguous conventional configs.
printf '0\n' | (cd /tmp && bash "$root/menu.sh") > "$root/menu-en.txt"
grep -q 'Services' "$root/menu-en.txt"
printf '0\n' | bash "$root/menu.sh" --locale pt-BR > "$root/menu-pt.txt"
grep -q 'Serviços' "$root/menu-pt.txt"
sha256sum "$root/$config" > "$root/config-before"
"$bin" scan --root "$root" --config "$config" --locale en > "$root/scan-en.txt"
"$bin" scan --root "$root" --config "$config" --locale pt-BR > "$root/scan-pt.txt"
grep -q 'Project:' "$root/scan-en.txt"
grep -q 'Projeto:' "$root/scan-pt.txt"
sha256sum "$root/$config" > "$root/config-after"
cmp "$root/config-before" "$root/config-after"
"$bin" preview app:probe --root "$root" --config "$config" > "$root/preview.txt"
test ! -e "$root/preview-only-marker"
"$bin" preview app:secret --root "$root" --config "$config" > "$root/secret-preview.txt"
! grep -q 'fixture private words' "$root/secret-preview.txt"
grep -q 'REDACTED' "$root/secret-preview.txt"
for action in up up status restart status logs down; do
  printf '\n>>> selected config wrapper: %s\n' "$action"
  (cd /tmp && bash "$root/.mudarro/scripts/app/$action.sh")
done
"$bin" run app:status --root "$root" --config "$config" | grep -E 'stopped|parado'
# Check live supervised state, not just successful dispatch.
bash "$root/.mudarro/scripts/app/up.sh"
"$bin" run app:status --root "$root" --config "$config" | grep -E 'running|rodando'
bash "$root/.mudarro/scripts/app/down.sh"
reject "$bin" menu --root "$root" --config "$config" --theme invalid </dev/null
echo 'UI selected configuration, supervisor, locale and preview end-to-end: OK'
