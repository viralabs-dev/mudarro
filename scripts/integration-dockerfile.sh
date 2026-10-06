#!/usr/bin/env bash
set -euo pipefail
bin="$(realpath "${1:?binary}")"
root="$(mktemp -d /tmp/mudarro-dockerfile.XXXXXX)"
project="$(printf '%s/.' "$root" | sha256sum | cut -c1-12)"
container="mudarro-$project-app"
image="mudarro-dockerfile-test:$project"
cleanup() { docker rm -f "$container" >/dev/null 2>&1 || true; docker image rm "$image" >/dev/null 2>&1 || true; rm -rf -- "$root"; }
trap cleanup EXIT
cat > "$root/Dockerfile" <<'DOCKER'
FROM alpine:3.21
CMD ["sh", "-c", "echo dockerfile-ready; sleep 300"]
DOCKER
cat > "$root/mudarro.yaml" <<YAML
version: 1
name: Dockerfile test
services:
- id: app
  dir: .
  infrastructure:
    kind: docker
    mode: dockerfile
    file: Dockerfile
    image: $image
YAML
"$bin" generate --root "$root"
"$bin" run app:image-build --root "$root"
"$bin" run app:up --root "$root"
id="$(docker inspect -f '{{.Id}}' "$container")"
"$bin" run app:up --root "$root"
test "$id" = "$(docker inspect -f '{{.Id}}' "$container")"
"$bin" run app:status --root "$root" >/dev/null
"$bin" run app:logs --root "$root" | grep dockerfile-ready
"$bin" run app:restart --root "$root"
"$bin" run app:down --root "$root"
test "$(docker inspect -f '{{.State.Running}}' "$container")" = false
"$bin" run app:up --root "$root"
test "$id" = "$(docker inspect -f '{{.Id}}' "$container")"
"$bin" run app:down --root "$root"
docker rm "$container" >/dev/null
# Foreign ownership refusal: this container belongs only to this disposable fixture.
docker run -d --name "$container" --label dev.viralabs.mudarro=foreign "$image" >/dev/null
if "$bin" run app:down --root "$root"; then echo 'Foreign ownership accepted' >&2; exit 1; fi
test "$(docker inspect -f '{{.State.Running}}' "$container")" = true
echo 'Dockerfile: OK; repeat/resume and ownership refusal verified'
