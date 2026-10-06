#!/bin/bash
set -euo pipefail
base=/tmp/mudarro-mud033-runtime
bin=$base/mudarro
root=$base/fixture
export PATH="$base:/home/danielsouza/.bun/bin:$PATH"
cleanup() { "$bin" run app-javascript:down --root "$root" > "$base/logs/cleanup.txt" 2>&1 || true; }
trap cleanup EXIT
step() { printf '\nCOMMAND:'; printf ' %q' "$@"; printf '\n'; "$@"; }
step "$bin" init --root "$root"
sha256sum "$root/mudarro.yaml" > "$base/logs/config-before.sha256"
if "$bin" init --root "$root"; then echo 'ERROR second init overwrote config'; exit 1; fi
sha256sum "$root/mudarro.yaml" > "$base/logs/config-after.sha256"
cmp "$base/logs/config-before.sha256" "$base/logs/config-after.sha256"
step "$bin" generate --root "$root" --dry-run
test ! -e "$root/menu.sh"
step "$bin" generate --root "$root"
find "$root" -type f -exec sha256sum {} + | sort > "$base/logs/generated-before.sha256"
step "$bin" generate --root "$root"
find "$root" -type f -exec sha256sum {} + | sort > "$base/logs/generated-after.sha256"
cmp "$base/logs/generated-before.sha256" "$base/logs/generated-after.sha256"
find "$root/.mudarro/scripts" -name '*.sh' -exec bash -n {} +
bash -n "$root/menu.sh"
step "$bin" doctor --root "$root"
step bash "$root/.mudarro/scripts/app-javascript/test.sh"
step bash "$root/.mudarro/scripts/app-javascript/lint.sh"
step bash "$root/.mudarro/scripts/app-javascript/build.sh"
printf '1\n4\n1\n0\n0\n0\n' | bash "$root/menu.sh" > "$base/logs/menu.txt"
grep -q 'MUD033 Bun check PASS' "$base/logs/menu.txt"
for action in up up status logs restart status logs down status; do
 step bash "$root/.mudarro/scripts/app-javascript/$action.sh"
done
"$bin" run app-javascript:status --root "$root" | grep -E 'stopped|parado'
echo 'BUN REAL E2E PASS'
