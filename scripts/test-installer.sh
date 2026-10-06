#!/usr/bin/env bash
set -euo pipefail
repo_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp="$(mktemp -d)"
trap 'rm -rf -- "$tmp"' EXIT
mkdir -p "$tmp/fakebin" "$tmp/payload" "$tmp/assets" "$tmp/install"
printf '#!/usr/bin/env bash\necho test-version\n' > "$tmp/payload/mudarro"
chmod +x "$tmp/payload/mudarro"
case "$(uname -s)" in Linux) platform=linux ;; Darwin) platform=darwin ;; esac
case "$(uname -m)" in x86_64) arch=amd64 ;; arm64|aarch64) arch=arm64 ;; esac
asset="mudarro_${platform}_${arch}.tar.gz"
tar -czf "$tmp/assets/$asset" -C "$tmp/payload" mudarro
(cd "$tmp/assets" && shasum -a 256 "$asset" > checksums.txt)
cat > "$tmp/fakebin/curl" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
url=''
while (($#)); do
  case "$1" in
    https://*) url="$1"; shift ;;
    -o) output="$2"; shift 2 ;;
    *) shift ;;
  esac
done
cp "$TEST_ASSETS/${url##*/}" "$output"
SH
chmod +x "$tmp/fakebin/curl"
export TEST_ASSETS="$tmp/assets" MUDARRO_INSTALL_DIR="$tmp/install" MUDARRO_VERSION=v0.1.0
export PATH="$tmp/fakebin:$PATH"
bash "$repo_dir/install.sh"
bash "$repo_dir/install.sh"
test "$("$tmp/install/mudarro" version)" = test-version
printf corrupt >> "$tmp/assets/$asset"
if bash "$repo_dir/install.sh"; then echo 'Aceitou checksum inválido' >&2; exit 1; fi
test "$("$tmp/install/mudarro" version)" = test-version
rm -- "$tmp/install/mudarro"
test ! -e "$tmp/install/mudarro"
echo 'Instalação, atualização e remoção simuladas: OK'
