package mudarro

func (prismaAdapter) Actions(s Service) []Action {
	var a []Action

	a = append(a, action("db-migrate", "banco", "npx", "--no-install", "prisma", "migrate", "deploy"), action("db-migration-new", "banco", "npx", "--no-install", "prisma", "migrate", "dev", "--create-only", "--name", "{name}"), action("db-seed", "banco", "npx", "--no-install", "prisma", "db", "seed"), action("db-reset", "banco", "npx", "--no-install", "prisma", "migrate", "reset", "--force"))

	for j := range a {
		if a[j].Name == "db-reset" {
			a[j].Command.Destructive = true
		}
	}
	return a
}
func (djangoAdapter) Actions(s Service) []Action {
	var a []Action

	for _, v := range []struct {
		n    string
		args []string
	}{{"db-migrate", []string{"manage.py", "migrate"}}, {"db-migration-new", []string{"manage.py", "makemigrations", "--name", "{name}"}}, {"db-seed", []string{"manage.py", "loaddata", "seeds/initial.json"}}, {"db-reset", []string{"manage.py", "flush", "--noinput"}}} {
		a = append(a, action(v.n, "banco", pythonArgs(s.Manager, v.args...)...))
	}

	for j := range a {
		if a[j].Name == "db-reset" {
			a[j].Command.Destructive = true
		}
	}
	return a
}
func (alembicAdapter) Actions(s Service) []Action {
	var a []Action

	a = append(a, action("db-migrate", "banco", pythonArgs(s.Manager, "-m", "alembic", "upgrade", "head")...), action("db-migration-new", "banco", pythonArgs(s.Manager, "-m", "alembic", "revision", "-m", "{name}")...), action("db-seed", "banco", pythonArgs(s.Manager, "seeds/seed.py")...))

	for j := range a {
		if a[j].Name == "db-reset" {
			a[j].Command.Destructive = true
		}
	}
	return a
}
func (gooseAdapter) Actions(s Service) []Action {
	var a []Action
	d := s.Database

	driver := "postgres"
	if d.Kind == "sqlite" {
		driver = "sqlite3"
	}
	env := d.URLenv
	if env == "" {
		env = "DATABASE_URL"
	}
	url := "${" + env + "}"
	if d.Kind == "sqlite" {
		url = d.Path
		if url == "" {
			url = "app.db"
		}
	}
	a = append(a, action("db-migrate", "banco", "goose", "-dir", "migrations", driver, url, "up"), action("db-migration-new", "banco", "goose", "-dir", "migrations", "create", "{name}", "sql"), action("db-seed", "banco", "bash", "seeds/seed.sh"), action("db-reset", "banco", "goose", "-dir", "migrations", driver, url, "reset"))

	for j := range a {
		if a[j].Name == "db-reset" {
			a[j].Command.Destructive = true
		}
	}
	return a
}
