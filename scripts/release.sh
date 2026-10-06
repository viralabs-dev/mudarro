#!/usr/bin/env bash
set -euo pipefail
version="${1:?uso: scripts/release.sh vX.Y.Z}"
mkdir -p dist
for platform in linux darwin; do
  for arch in amd64 arm64; do
    dir="$(mktemp -d)"
    CGO_ENABLED=0 GOOS="$platform" GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$dir/mudarro" ./cmd/mudarro
    tar -czf "dist/mudarro_${platform}_${arch}.tar.gz" -C "$dir" mudarro
    rm -rf -- "$dir"
  done
done
(cd dist && shasum -a 256 mudarro_*.tar.gz > checksums.txt)
