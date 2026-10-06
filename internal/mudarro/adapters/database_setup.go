package adapters

import "github.com/viralabs-dev/mudarro/internal/mudarro/model"

// Dependency installation is an explicit action, never part of scan or menu rendering.
func SetupActions(s model.Service) []model.Action {
	var install model.Command
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
	var a []model.Action
	if s.Database.Generate && s.Database.Kind == "sqlite" {
		a = append(a, model.Action{Name: "db-init", Group: "banco", Special: "sqlite-init"})
	}
	if len(install.Args) > 0 {
		a = append(a, model.Action{Name: "db-install", Group: "dependencias", Command: install})
	}
	if s.Database.Generate && s.Database.Kind == "postgresql" && s.Infrastructure.Kind == "local" {
		for _, n := range []string{"init", "up", "down", "status", "create"} {
			c := cmd("banco", "bash", "mudarro-postgres.sh", n)
			c.Requires = []string{"initdb", "pg_ctl", "createdb"}
			a = append(a, model.Action{Name: "db-" + n, Group: "banco", Command: c})
		}
	}
	return a
}
