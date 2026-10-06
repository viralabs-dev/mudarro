#!/usr/bin/env bash
set -euo pipefail
base=/tmp/mudarro-runtimes.K7F4J1
export PATH="$base/poetry/bin:$base/poetry/tool-venv/bin:$PATH"
export POETRY_CACHE_DIR="$base/cache/poetry/runtime"
export POETRY_CONFIG_DIR="$base/poetry/config"
export POETRY_DATA_DIR="$base/poetry/data"
export POETRY_VIRTUALENVS_IN_PROJECT=true
export POETRY_KEYRING_ENABLED=false
export POETRY_NO_INTERACTION=1
cd "$base/fixtures/poetry"
poetry --version
poetry env use /home/linuxbrew/.linuxbrew/bin/python3
poetry lock
poetry install --no-root
poetry run python probe.py
bin=/tmp/mudarro-timeout-checkpoint/mudarro
"$bin" scan --json > "$base/logs/poetry/scan.json"
if [ ! -f mudarro.yaml ]; then "$bin" init --name 'Poetry runtime'; fi
cat > mudarro.yaml <<'YAML'
version: 1
name: Poetry runtime
services:
  - id: app-python
    dir: .
    language: python
    manager: poetry
    infrastructure: {kind: local}
    commands:
      install: {args: [poetry, install], group: dependencias}
      start: {args: [poetry, run, python, app.py], group: aplicacao}
      test: {args: [poetry, run, python, probe.py], group: qualidade}
YAML
"$bin" generate
find .mudarro/scripts -type f -print
sha256sum menu.sh .mudarro/generated.json > "$base/logs/poetry/hashes-before.txt"
"$bin" generate
sha256sum menu.sh .mudarro/generated.json > "$base/logs/poetry/hashes-after.txt"
cmp "$base/logs/poetry/hashes-before.txt" "$base/logs/poetry/hashes-after.txt"
"$bin" run app-python:test
trap '"$bin" run app-python:down >/dev/null 2>&1 || true' EXIT
"$bin" run app-python:up
"$bin" run app-python:up
"$bin" run app-python:status
"$bin" run app-python:logs
"$bin" run app-python:restart
"$bin" run app-python:down
printf '0\n' | ./menu.sh
.mudarro/scripts/app-python/test.sh
if "$bin" init; then echo 'FAIL repeated init accepted'; exit 1; fi
if "$bin" run app-python:missing; then echo 'FAIL unknown action accepted'; exit 1; fi
printf '\n# user change\n' >> menu.sh
sha256sum menu.sh > "$base/logs/poetry/manual-before.txt"
if "$bin" generate; then echo 'FAIL manual edit accepted'; exit 1; fi
sha256sum menu.sh > "$base/logs/poetry/manual-after.txt"
cmp "$base/logs/poetry/manual-before.txt" "$base/logs/poetry/manual-after.txt"
echo 'Poetry runtime E2E OK'
