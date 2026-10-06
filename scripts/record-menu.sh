#!/usr/bin/env bash
# Run under util-linux script to record genuine terminal output.
set -euo pipefail
bin="$(realpath "${1:?absolute binary}")"
repo="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
export PATH="$(dirname "$bin"):$PATH" TERM=xterm-256color COLUMNS=100
root="$(mktemp -d /tmp/mudarro-visual.XXXXXX)"
trap 'rm -rf -- "$root"' EXIT
cp "$repo/docs/samples/menu-auto/"{package.json,app.js} "$root/"
stty cols 100 rows 32
screen() { printf '\033[2J\033[H'; printf 'Mudarro — sessão real de terminal / %s\n\n' "$1"; }
screen 'detecção offline'
"$bin" scan --root "$root"
sleep 2
screen 'criação automática'
"$bin" init --root "$root" --name "${MUDARRO_SAMPLE_NAME:-Aurora}" --select app-javascript:hello
"$bin" generate --root "$root"
sleep 2
screen 'menu, submenus e ação de qualidade'
{ sleep 2; printf '1\n'; sleep 2; printf '4\n'; sleep 2; printf '1\n'; sleep 2; printf '\n0\n0\n0\n'; } | bash "$root/menu.sh"
screen 'idempotência e proteção'
"$bin" generate --root "$root"
printf 'Geração repetida: sem arquivos alterados\n'
printf '# edição manual\n' >> "$root/menu.sh"
if "$bin" generate --root "$root"; then exit 1; fi
printf '\nEdição preservada; recusa esperada.\n'
sleep 3
