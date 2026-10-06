#!/usr/bin/env bash
# Real CLI + generated launchers; fixtures never mutate the saved sample.
set -euo pipefail
bin="$(realpath "${1:?absolute binary}")"
repo="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
root="$(mktemp -d '/tmp/mudarro menu.XXXXXX')"
export PATH="$(dirname "$bin"):$PATH"
trap '"$bin" run app-javascript:down --root "$root" >/dev/null 2>&1 || true; rm -rf -- "$root"' EXIT
cp "$repo/docs/samples/menu-auto/package.json" "$repo/docs/samples/menu-auto/app.js" "$root/"
step() { printf '\n>>> %s\n' "$*"; "$@"; }
fail() { if "$@"; then echo 'Expected rejection' >&2; exit 1; fi; }
step "$bin" scan --root "$root" --json
step "$bin" init --root "$root" --select app-javascript:hello
before="$(sha256sum "$root/mudarro.yaml")"
fail "$bin" init --root "$root"
test "$before" = "$(sha256sum "$root/mudarro.yaml")"
step "$bin" generate --root "$root" --dry-run
test ! -e "$root/menu.sh"
step "$bin" generate --root "$root"
find "$root/.mudarro/scripts" -name '*.sh' -exec bash -n {} +
bash -n "$root/menu.sh"
find "$root" -type f -exec sha256sum {} + | sort > "$root.before"
step "$bin" generate --root "$root"
find "$root" -type f -exec sha256sum {} + | sort > "$root.after"
cmp "$root.before" "$root.after"
rm "$root.before" "$root.after"
step "$bin" doctor --root "$root"
printf 'abc\n1\n0\n0\n' | (cd /tmp && bash "$root/menu.sh") | tee "$root/menu-output"
grep -q 'Invalid option' "$root/menu-output"
grep -q 'quality' "$root/menu-output"
printf '1\n4\n1\n0\n0\n0\n' | bash "$root/menu.sh" | tee "$root/menu-action-output"
grep -q 'sample quality OK' "$root/menu-action-output"
step bash "$root/.mudarro/scripts/app-javascript/test.sh"
step bash "$root/.mudarro/scripts/app-javascript/hello.sh"
step "$bin" run app-javascript:up --root "$root"
step "$bin" run app-javascript:up --root "$root"
step "$bin" run app-javascript:status --root "$root"
step "$bin" run app-javascript:logs --root "$root"
step "$bin" run app-javascript:restart --root "$root"
step "$bin" run app-javascript:down --root "$root"
"$bin" run app-javascript:status --root "$root" | grep -E 'stopped|parado'
fail "$bin" run app-javascript:nonexistent --root "$root"
printf '# user edit\n' >> "$root/menu.sh"
before="$(sha256sum "$root/menu.sh" "$root/.mudarro/generated.json")"
fail "$bin" generate --root "$root"
test "$before" = "$(sha256sum "$root/menu.sh" "$root/.mudarro/generated.json")"
echo 'Automatic menu end-to-end: OK'
