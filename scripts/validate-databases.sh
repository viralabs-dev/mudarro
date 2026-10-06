#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário}"
name="mudarro-db-validation-$$"
trap 'docker rm -f "$name" >/dev/null' EXIT
docker run -d --name "$name" -e POSTGRES_PASSWORD=mudarro_test_only -e POSTGRES_DB=mudarro_test -p 127.0.0.1::5432 postgres:17 >/dev/null
ready=false
for attempt in $(seq 1 60); do
  if docker exec "$name" pg_isready -U postgres >/dev/null 2>&1; then ready=true; break; fi
  sleep 1
done
"$ready" || { echo 'PostgreSQL não ficou pronto' >&2; exit 1; }
port="$(docker port "$name" 5432/tcp)"
export TEST_POSTGRES_URL="postgresql://postgres:mudarro_test_only@127.0.0.1:${port##*:}/mudarro_test"
bash scripts/integration-databases.sh "$bin"
