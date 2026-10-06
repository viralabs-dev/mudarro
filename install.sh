#!/usr/bin/env bash
set -euo pipefail

repo="${MUDARRO_REPOSITORY:-viralabs-dev/mudarro}"
install_dir="${MUDARRO_INSTALL_DIR:-$HOME/.local/bin}"
version="${MUDARRO_VERSION:-latest}"
case "$(uname -s)" in Linux) platform=linux ;; Darwin) platform=darwin ;; *) echo 'Sistema não suportado; Windows deve usar WSL.' >&2; exit 1 ;; esac
case "$(uname -m)" in x86_64|amd64) arch=amd64 ;; aarch64|arm64) arch=arm64 ;; *) echo 'Arquitetura não suportada.' >&2; exit 1 ;; esac
for tool in curl tar mktemp; do command -v "$tool" >/dev/null || { echo "Requisito ausente: $tool" >&2; exit 1; }; done
if [[ "$version" != latest && ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+([.-][A-Za-z0-9.-]+)?$ ]]; then
  echo 'MUDARRO_VERSION deve ser latest ou uma tag vX.Y.Z.' >&2; exit 1
fi
asset="mudarro_${platform}_${arch}.tar.gz"
if [[ "$version" == latest ]]; then
  base="https://github.com/$repo/releases/latest/download"
else
  base="https://github.com/$repo/releases/download/$version"
fi
tmp="$(mktemp -d)"
trap 'rm -rf -- "$tmp"' EXIT
curl --proto '=https' --tlsv1.2 -fsSL "$base/$asset" -o "$tmp/$asset"
curl --proto '=https' --tlsv1.2 -fsSL "$base/checksums.txt" -o "$tmp/checksums.txt"
expected="$(awk -v file="$asset" '$2 == file { print $1 }' "$tmp/checksums.txt")"
[[ "$expected" =~ ^[0-9a-f]{64}$ ]] || { echo 'Checksum ausente ou inválido.' >&2; exit 1; }
if command -v sha256sum >/dev/null; then
  actual="$(sha256sum "$tmp/$asset")"
elif command -v shasum >/dev/null; then
  actual="$(shasum -a 256 "$tmp/$asset")"
else
  echo 'Instale sha256sum ou shasum.' >&2; exit 1
fi
[[ "${actual%% *}" == "$expected" ]] || { echo 'Checksum divergente; instalação interrompida.' >&2; exit 1; }
tar -xzf "$tmp/$asset" -C "$tmp" mudarro
mkdir -p "$install_dir"
[[ ! -L "$install_dir/mudarro" ]] || { echo 'Destino é symlink; escolha outro diretório.' >&2; exit 1; }
target="$(mktemp "$install_dir/.mudarro-install.XXXXXX")"
cp "$tmp/mudarro" "$target"
chmod 755 "$target"
mv -f "$target" "$install_dir/mudarro"
printf 'Mudarro instalado: %s/mudarro\n' "$install_dir"
case ":$PATH:" in *":$install_dir:"*) ;; *) printf 'Adicione ao PATH: export PATH="%s:$PATH"\n' "$install_dir" ;; esac
"$install_dir/mudarro" version
