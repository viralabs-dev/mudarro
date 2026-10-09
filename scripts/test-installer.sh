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
for bad in 1.0.0 v1.0 'v1.0.0;id' ../v1.0.0; do
  if MUDARRO_VERSION="$bad" bash "$repo_dir/install.sh"; then echo "Aceitou versão inválida: $bad" >&2; exit 1; fi
done
for bad in 'evil/../x' 'a b/c' owner owner/.. https://x/y a/b/c; do
  if MUDARRO_REPOSITORY="$bad" bash "$repo_dir/install.sh"; then echo "Aceitou repositório inválido: $bad" >&2; exit 1; fi
done
# New archives install their complete notices without affecting legacy support.
cp "$repo_dir/LICENSE" "$repo_dir/THIRD_PARTY_NOTICES.md" "$tmp/payload/"
cp -R "$repo_dir/LICENSES" "$tmp/payload/"
tar -czf "$tmp/assets/$asset" -C "$tmp/payload" mudarro LICENSE THIRD_PARTY_NOTICES.md LICENSES
(cd "$tmp/assets" && shasum -a 256 "$asset" > checksums.txt)
bash "$repo_dir/install.sh"
bash "$repo_dir/install.sh"
cmp "$tmp/payload/LICENSE" "$tmp/install/mudarro-licenses/LICENSE"
cmp "$tmp/payload/THIRD_PARTY_NOTICES.md" "$tmp/install/mudarro-licenses/THIRD_PARTY_NOTICES.md"
for license in golang.org-x-mod.txt go-toml-v2.txt gopkg.in-yaml.v3.txt; do
  cmp "$tmp/payload/LICENSES/$license" "$tmp/install/mudarro-licenses/LICENSES/$license"
done
# Incomplete new payload must preserve both binary and existing notices.
cp "$tmp/install/mudarro" "$tmp/installed-before"
printf '#!/usr/bin/env bash\necho forbidden-partial-update\n' > "$tmp/payload/mudarro"
tar -czf "$tmp/assets/$asset" -C "$tmp/payload" mudarro THIRD_PARTY_NOTICES.md LICENSES/golang.org-x-mod.txt
(cd "$tmp/assets" && shasum -a 256 "$asset" > checksums.txt)
if bash "$repo_dir/install.sh"; then echo 'Aceitou payload de licenças incompleto' >&2; exit 1; fi
cmp "$tmp/installed-before" "$tmp/install/mudarro"
cmp "$tmp/payload/THIRD_PARTY_NOTICES.md" "$tmp/install/mudarro-licenses/THIRD_PARTY_NOTICES.md"
printf '#!/usr/bin/env bash\necho test-version\n' > "$tmp/payload/mudarro"
tar -czf "$tmp/assets/$asset" -C "$tmp/payload" mudarro THIRD_PARTY_NOTICES.md LICENSES
(cd "$tmp/assets" && shasum -a 256 "$asset" > checksums.txt)
# Refuse only symlinks created in this test's private destinations.
mkdir -p "$tmp/foreign"
printf 'owned sentinel' > "$tmp/foreign/sentinel"
mv "$tmp/install/mudarro-licenses" "$tmp/saved-licenses"
ln -s "$tmp/foreign" "$tmp/install/mudarro-licenses"
if bash "$repo_dir/install.sh"; then echo 'Aceitou symlink de avisos' >&2; exit 1; fi
cmp "$tmp/installed-before" "$tmp/install/mudarro"
test "$(cat "$tmp/foreign/sentinel")" = 'owned sentinel'
test "$(find "$tmp/foreign" -type f | wc -l)" -eq 1
rm "$tmp/install/mudarro-licenses"
mv "$tmp/saved-licenses" "$tmp/install/mudarro-licenses"
mv "$tmp/install/mudarro-licenses/LICENSES" "$tmp/saved-inner-licenses"
ln -s "$tmp/foreign" "$tmp/install/mudarro-licenses/LICENSES"
if bash "$repo_dir/install.sh"; then echo 'Aceitou symlink de licenças' >&2; exit 1; fi
cmp "$tmp/installed-before" "$tmp/install/mudarro"
test "$(find "$tmp/foreign" -type f | wc -l)" -eq 1
rm "$tmp/install/mudarro-licenses/LICENSES"
mv "$tmp/saved-inner-licenses" "$tmp/install/mudarro-licenses/LICENSES"
printf corrupt >> "$tmp/assets/$asset"
if bash "$repo_dir/install.sh"; then echo 'Aceitou checksum inválido' >&2; exit 1; fi
test "$("$tmp/install/mudarro" version)" = test-version
cmp "$tmp/payload/THIRD_PARTY_NOTICES.md" "$tmp/install/mudarro-licenses/THIRD_PARTY_NOTICES.md"
for license in golang.org-x-mod.txt go-toml-v2.txt gopkg.in-yaml.v3.txt; do
  cmp "$tmp/payload/LICENSES/$license" "$tmp/install/mudarro-licenses/LICENSES/$license"
done
rm -- "$tmp/install/mudarro"
test ! -e "$tmp/install/mudarro"
echo 'Instalação, atualização e remoção simuladas: OK'
