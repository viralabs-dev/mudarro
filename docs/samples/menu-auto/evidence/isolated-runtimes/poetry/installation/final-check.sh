#!/usr/bin/env bash
set -euo pipefail
base=/tmp/mudarro-runtimes.K7F4J1
export PATH="$base/poetry/bin:$base/poetry/tool-venv/bin:$PATH"
export POETRY_CACHE_DIR="$base/cache/poetry/runtime" POETRY_CONFIG_DIR="$base/poetry/config" POETRY_DATA_DIR="$base/poetry/data" POETRY_VIRTUALENVS_IN_PROJECT=true POETRY_KEYRING_ENABLED=false POETRY_NO_INTERACTION=1
cd "$base/fixtures/poetry"
bin=/tmp/mudarro-timeout-checkpoint/mudarro
trap '"$bin" run app-python:down >/dev/null 2>&1 || true' EXIT
"$bin" run app-python:install
"$bin" run app-python:up
for attempt in 1 2 3 4 5 6 7 8 9 10; do
 if "$bin" run app-python:logs | grep -q 'Poetry fixture ready'; then break; fi
 sleep 0.2
done
"$bin" run app-python:logs | grep 'Poetry fixture ready'
"$bin" run app-python:down
"$bin" run app-python:status
printf '0\n' | script -q -e -O "$base/logs/poetry/session.txt" -T "$base/logs/poetry/session.timing" -c 'poetry run python probe.py; /tmp/mudarro-timeout-checkpoint/mudarro menu'
