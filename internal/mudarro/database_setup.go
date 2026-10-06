package mudarro

// Dependency installation is an explicit action, never part of scan or menu rendering.
func databaseSetupActions(s Service) []Action {
	var install Command
	switch s.Database.Tool {
	case "prisma":
		switch s.Manager {
		case "pnpm":
			install = cmd("dependencias", "pnpm", "add", "-D", "prisma@7", "@prisma/client@7")
		case "yarn":
			install = cmd("dependencias", "yarn", "add", "-D", "prisma@7", "@prisma/client@7")
		default:
			install = cmd("dependencias", "npm", "install", "--save-dev", "prisma@7", "@prisma/client@7")
		}
	case "django", "alembic":
		packages := []string{s.Database.Tool}
		if s.Database.Tool == "alembic" {
			packages = append(packages, "sqlalchemy")
		}
		if s.Database.Kind == "postgresql" {
			packages = append(packages, "psycopg[binary]")
		}
		prefix := []string{".venv/bin/python", "-m", "pip", "install"}
		if s.Manager == "uv" {
			prefix = []string{"uv", "add"}
		}
		if s.Manager == "poetry" {
			prefix = []string{"poetry", "add"}
		}
		install = cmd("dependencias", append(prefix, packages...)...)
	case "goose":
		install = cmd("dependencias", "go", "install", "github.com/pressly/goose/v3/cmd/goose@v3.24.1")
	}
	var a []Action
	if s.Database.Generate && s.Database.Kind == "sqlite" {
		a = append(a, Action{Name: "db-init", Group: "banco", Special: "sqlite-init"})
	}
	if len(install.Args) > 0 {
		a = append(a, Action{Name: "db-install", Group: "dependencias", Command: install})
	}
	if s.Database.Generate && s.Database.Kind == "postgresql" && s.Infrastructure.Kind == "local" {
		for _, n := range []string{"init", "up", "down", "status", "create"} {
			c := cmd("banco", "bash", "mudarro-postgres.sh", n)
			c.Requires = []string{"initdb", "pg_ctl", "createdb"}
			a = append(a, Action{Name: "db-" + n, Group: "banco", Command: c})
		}
	}
	return a
}

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
