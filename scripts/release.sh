#!/usr/bin/env bash
set -euo pipefail
version="${1:?uso: scripts/release.sh vX.Y.Z}"
command -v zip >/dev/null || { echo 'Requisito ausente: zip' >&2; exit 1; }
mkdir -p dist
dist="$(cd dist && pwd)"
for platform in linux darwin; do
  for arch in amd64 arm64; do
    dir="$(mktemp -d)"
    CGO_ENABLED=0 GOOS="$platform" GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$dir/mudarro" ./cmd/mudarro
    cp LICENSE THIRD_PARTY_NOTICES.md "$dir/"
    cp -R LICENSES "$dir/"
    tar -czf "dist/mudarro_${platform}_${arch}.tar.gz" -C "$dir" mudarro LICENSE THIRD_PARTY_NOTICES.md LICENSES
    rm -rf -- "$dir"
  done
done
# Windows ships a zip with mudarro.exe and the same notices; install.ps1 consumes it.
for arch in amd64 arm64; do
  dir="$(mktemp -d)"
  CGO_ENABLED=0 GOOS=windows GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$dir/mudarro.exe" ./cmd/mudarro
  cp LICENSE THIRD_PARTY_NOTICES.md "$dir/"
  cp -R LICENSES "$dir/"
  archive="$dist/mudarro_windows_${arch}.zip"
  rm -f -- "$archive"
  # Fixed member order and timestamps, no directory entries or extra attributes: reproducible zip.
  (
    cd "$dir"
    find . -type f -exec env TZ=UTC touch -h -t 198001010000.00 {} +
    { printf '%s\n' mudarro.exe LICENSE THIRD_PARTY_NOTICES.md; find LICENSES -type f | LC_ALL=C sort; } |
      TZ=UTC zip -q -X -D "$archive" -@
  )
  rm -rf -- "$dir"
done
(cd dist && shasum -a 256 mudarro_*.tar.gz mudarro_*.zip > checksums.txt)
