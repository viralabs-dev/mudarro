package mudarro

const localPostgresScript = `#!/usr/bin/env bash
set -euo pipefail
cd -- "$(dirname -- "${BASH_SOURCE[0]}")"
data="$PWD/.mudarro/postgres"
mkdir -p "$PWD/.mudarro"
[[ ! -L "$data" && ! -L "$PWD/.mudarro" ]] || { echo 'Symlink não permitido no estado PostgreSQL' >&2; exit 1; }
case "${1:-}" in
  init)
    : "${POSTGRES_USER:?configure POSTGRES_USER}" "${POSTGRES_PASSWORD:?configure POSTGRES_PASSWORD}"
    if [[ -f "$data/PG_VERSION" ]]; then echo 'PostgreSQL já inicializado'; exit 0; fi
    pwfile="$(mktemp)"
    trap 'rm -f -- "$pwfile"' EXIT
    chmod 600 "$pwfile"
    printf '%s' "$POSTGRES_PASSWORD" > "$pwfile"
    initdb -D "$data" -U "$POSTGRES_USER" --pwfile="$pwfile" --auth=scram-sha-256
    ;;
  up)
    : "${POSTGRES_PORT:?configure POSTGRES_PORT}"
    [[ "$POSTGRES_PORT" =~ ^[0-9]+$ ]] && ((POSTGRES_PORT > 0 && POSTGRES_PORT < 65536)) || exit 1
    pg_ctl -D "$data" status >/dev/null 2>&1 && exit 0
    pg_ctl -D "$data" -l "$PWD/.mudarro/postgres.log" -o "-h 127.0.0.1 -p $POSTGRES_PORT -k ''" start
    ;;
  down) pg_ctl -D "$data" stop -m fast ;;
  status) pg_ctl -D "$data" status ;;
  create)
    : "${POSTGRES_USER:?}" "${POSTGRES_PASSWORD:?}" "${POSTGRES_DB:?}" "${POSTGRES_PORT:?}"
    PGPASSWORD="$POSTGRES_PASSWORD" createdb -h 127.0.0.1 -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -- "$POSTGRES_DB"
    ;;
  *) echo 'uso: mudarro-postgres.sh init|up|down|status|create' >&2; exit 1 ;;
esac
`
