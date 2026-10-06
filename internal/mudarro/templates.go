package mudarro

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
)

func scaffoldFiles(s Service) (map[string]GeneratedFile, error) {
	files := map[string]GeneratedFile{}
	put := func(p, b string) { files[p] = GeneratedFile{[]byte(b), 0644} }
	i, d := s.Infrastructure, s.Database
	if i.Generate {
		if i.Kind != "docker" && i.Kind != "podman" && i.Kind != "kubernetes" {
			return nil, fmt.Errorf("generate de infraestrutura suporta docker, podman ou kubernetes")
		}
		start, ok := s.Commands["start"]
		if !ok {
			return nil, fmt.Errorf("commands.start obrigatório para gerar imagem")
		}
		var dockerfile string
		switch s.Language {
		case "javascript", "typescript":
			switch s.Manager {
			case "npm":
				dockerfile = "FROM node:24-bookworm-slim\nWORKDIR /app\nCOPY . .\nRUN npm install\n"
			case "pnpm", "yarn":
				dockerfile = "FROM node:24-bookworm-slim\nWORKDIR /app\nRUN corepack enable\nCOPY . .\nRUN " + s.Manager + " install\n"
			default:
				return nil, fmt.Errorf("manager JS/TS deve ser npm, pnpm ou yarn")
			}
		case "python":
			dockerfile = "FROM python:3.13-slim\nWORKDIR /app\nCOPY . .\n"
			switch s.Manager {
			case "pip":
				dockerfile += "RUN python -m venv .venv\n"
				install, ok := s.Commands["install"]
				if !ok || len(install.Args) == 0 {
					return nil, fmt.Errorf("declare commands.install.args para Python pip")
				}
				dockerfile += "RUN " + jsonString(install.Args) + "\n"
			case "uv":
				dockerfile += "RUN pip install --no-cache-dir uv && uv sync\n"
			case "poetry":
				dockerfile += "RUN pip install --no-cache-dir poetry && poetry install\n"
			default:
				return nil, fmt.Errorf("manager Python não suportado")
			}
		case "go":
			dockerfile = "FROM golang:1.26-bookworm\nWORKDIR /app\nCOPY . .\nRUN go mod download\n"
		default:
			return nil, fmt.Errorf("declare Dockerfile existente para linguagem %s", s.Language)
		}
		args := start.Args
		if start.Shell != "" {
			args = []string{"/bin/sh", "-c", start.Shell}
		}
		dockerfile += "CMD " + jsonString(args) + "\n"
		dockerPath := "Dockerfile"
		if i.Mode == "dockerfile" && i.File != "" {
			dockerPath = i.File
		}
		put(dockerPath, dockerfile)
		put(".dockerignore", ".git\n.mudarro\n.env\n.env.*\n!.env.example\nnode_modules\n.venv\n__pycache__\n")
		if i.Mode == "compose" && (i.Kind == "docker" || i.Kind == "podman") {
			if i.File == "" {
				return nil, fmt.Errorf("declare infrastructure.file")
			}
			service := map[string]any{"build": map[string]any{"context": ".", "dockerfile": dockerPath}}
			if i.Port > 0 {
				service["ports"] = []string{fmt.Sprintf("127.0.0.1:%d:%d", i.Port, i.Port)}
			}
			services := map[string]any{s.ID: service}
			doc := map[string]any{"services": services}
			if d.Generate && d.Kind == "postgresql" {
				services["db"] = map[string]any{"image": "postgres:17", "environment": map[string]string{"POSTGRES_DB": "${POSTGRES_DB:?configure POSTGRES_DB}", "POSTGRES_USER": "${POSTGRES_USER:?configure POSTGRES_USER}", "POSTGRES_PASSWORD": "${POSTGRES_PASSWORD:?configure POSTGRES_PASSWORD}"}, "ports": []string{"127.0.0.1:${POSTGRES_PORT:-5432}:5432"}, "volumes": []string{"mudarro-db:/var/lib/postgresql/data"}}
				doc["volumes"] = map[string]any{"mudarro-db": map[string]any{}}
				env := d.URLenv
				if env == "" {
					env = "DATABASE_URL"
				}
				service["environment"] = map[string]string{env: "${" + env + ":?configure conexão usando host db}"}
			}
			b, _ := yaml.Marshal(doc)
			put(i.File, string(b))
		}
		if i.Kind == "kubernetes" {
			if i.Mode != "manifests" && i.Mode != "kustomize" {
				return nil, fmt.Errorf("geração Kubernetes usa manifests ou kustomize; Helm integra charts existentes")
			}
			if i.Image == "" || i.File == "" || i.Port == 0 || i.Context == "" || i.Namespace == "" {
				return nil, fmt.Errorf("geração Kubernetes exige image, file (diretório), port, context e namespace")
			}
			labels := map[string]string{"app.kubernetes.io/instance": s.ID}
			container := map[string]any{"name": s.ID, "image": i.Image, "ports": []any{map[string]any{"containerPort": i.Port}}}
			if d.Generate && d.Kind == "postgresql" {
				env := d.URLenv
				if env == "" {
					env = "DATABASE_URL"
				}
				container["env"] = []any{map[string]any{"name": env, "valueFrom": map[string]any{"secretKeyRef": map[string]string{"name": s.ID + "-database", "key": "url"}}}}
			}
			deployment := map[string]any{"apiVersion": "apps/v1", "kind": "Deployment", "metadata": map[string]any{"name": s.ID, "namespace": i.Namespace, "labels": labels}, "spec": map[string]any{"replicas": 1, "selector": map[string]any{"matchLabels": labels}, "template": map[string]any{"metadata": map[string]any{"labels": labels}, "spec": map[string]any{"containers": []any{container}}}}}
			b, _ := yaml.Marshal(deployment)
			put(filepath.Join(i.File, "deployment.yaml"), string(b))
			service := map[string]any{"apiVersion": "v1", "kind": "Service", "metadata": map[string]any{"name": s.ID, "namespace": i.Namespace}, "spec": map[string]any{"selector": labels, "ports": []any{map[string]any{"port": i.Port, "targetPort": i.Port}}}}
			b, _ = yaml.Marshal(service)
			put(filepath.Join(i.File, "service.yaml"), string(b))
			resources := []string{"deployment.yaml", "service.yaml"}
			if d.Generate && d.Kind == "postgresql" {
				put(filepath.Join(i.File, "database.yaml"), postgresKubernetes(s))
				resources = append(resources, "database.yaml")
				put("database-secret.example.txt", "Crie o Secret "+s.ID+"-database no namespace "+i.Namespace+" com chaves user, password, database e url. Use seu gerenciador de segredos; não versione valores. A URL da aplicação usa o Service "+s.ID+"-db.\n")
			}
			if i.Mode == "kustomize" {
				b, _ = yaml.Marshal(map[string]any{"apiVersion": "kustomize.config.k8s.io/v1beta1", "kind": "Kustomization", "resources": resources})
				put(filepath.Join(i.File, "kustomization.yaml"), string(b))
			}
		}
	}
	if d.Generate {
		if d.Kind == "postgresql" && i.Kind == "local" {
			put("mudarro-postgres.sh", localPostgresScript)
		}
		if d.Kind != "postgresql" && d.Kind != "sqlite" {
			return nil, fmt.Errorf("database.kind deve ser postgresql ou sqlite")
		}
		if d.Tool == "" {
			return nil, fmt.Errorf("declare database.tool")
		}
		env := d.URLenv
		if env == "" {
			env = "DATABASE_URL"
		}
		put(".env.example", env+"=\nPOSTGRES_DB=\nPOSTGRES_USER=\nPOSTGRES_PASSWORD=\nPOSTGRES_PORT=5432\n")
		switch d.Tool {
		case "prisma":
			put("prisma/schema.prisma", "generator client {\n  provider = \"prisma-client-js\"\n}\n\ndatasource db {\n  provider = "+jsonString(d.Kind)+"\n}\n\n// Declare os modelos de domínio antes de gerar uma migration.\n")
			put("prisma.config.ts", "import { defineConfig, env } from 'prisma/config';\nexport default defineConfig({\n  schema: 'prisma/schema.prisma',\n  migrations: { path: 'prisma/migrations', seed: 'node prisma/seed.cjs' },\n  datasource: { url: env("+jsonString(env)+") },\n});\n")
			put("prisma/seed.cjs", "throw new Error('Implemente o seed idempotente do seu domínio antes de executar.');\n")
		case "alembic":
			put("alembic.ini", "[alembic]\nscript_location = migrations\nprepend_sys_path = .\n")
			put("migrations/env.py", "import os\nfrom alembic import context\nfrom sqlalchemy import create_engine, pool\nurl = os.environ["+jsonString(env)+"]\ntarget_metadata = None  # Importe o metadata dos modelos para autogenerate.\nif context.is_offline_mode():\n    context.configure(url=url, target_metadata=target_metadata, literal_binds=True)\n    with context.begin_transaction():\n        context.run_migrations()\nelse:\n    with create_engine(url, poolclass=pool.NullPool).connect() as connection:\n        context.configure(connection=connection, target_metadata=target_metadata)\n        with context.begin_transaction():\n            context.run_migrations()\n")
			put("migrations/script.py.mako", "\"\"\"${message}\"\"\"\nfrom alembic import op\nimport sqlalchemy as sa\n${imports if imports else \"\"}\nrevision = ${repr(up_revision)}\ndown_revision = ${repr(down_revision)}\nbranch_labels = ${repr(branch_labels)}\ndepends_on = ${repr(depends_on)}\ndef upgrade():\n    ${upgrades if upgrades else \"pass\"}\ndef downgrade():\n    ${downgrades if downgrades else \"pass\"}\n")
			put("migrations/versions/.gitkeep", "")
			put("seeds/seed.py", "raise SystemExit('Implemente o seed idempotente do domínio antes de executar.')\n")
		case "django":
			put("seeds/initial.json", "[]\n")
			put("mudarro_database.py", djangoDatabase(d, env))
			put("DATABASE-SETUP.md", "# Banco Django\n\nImporte `DATABASES` de `mudarro_database.py` no settings existente. Declare seus models e execute db-migration-new. O seed inicial vazio não inventa dados. Instale psycopg para PostgreSQL.\n")
		case "goose":
			put("migrations/.gitkeep", "")
			put("seeds/seed.sh", "#!/usr/bin/env bash\nset -euo pipefail\necho 'Implemente o seed idempotente do domínio antes de executar.' >&2\nexit 1\n")
		default:
			return nil, fmt.Errorf("database.tool não suportada: %s", d.Tool)
		}
	}
	return files, nil
}
func djangoDatabase(d Database, env string) string {
	if d.Kind == "sqlite" {
		p := d.Path
		if p == "" {
			p = "app.db"
		}
		return "from pathlib import Path\nDATABASES = {'default': {'ENGINE': 'django.db.backends.sqlite3', 'NAME': Path(__file__).resolve().parent / " + jsonString(p) + "}}\n"
	}
	return "import os\nfrom urllib.parse import urlparse, unquote\nu = urlparse(os.environ[" + jsonString(env) + "])\nDATABASES = {'default': {'ENGINE': 'django.db.backends.postgresql', 'NAME': u.path.lstrip('/'), 'USER': unquote(u.username or ''), 'PASSWORD': unquote(u.password or ''), 'HOST': u.hostname, 'PORT': u.port or 5432}}\n"
}
func postgresKubernetes(s Service) string {
	// No Secret is generated. The user supplies credentials through their existing secret workflow.
	return strings.NewReplacer("NAMESPACE", s.Infrastructure.Namespace, "NAME", s.ID).Replace(`apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: NAME-db
  namespace: NAMESPACE
spec:
  accessModes: [ReadWriteOnce]
  resources:
    requests:
      storage: 1Gi
---
apiVersion: v1
kind: Service
metadata:
  name: NAME-db
  namespace: NAMESPACE
spec:
  selector:
    app.kubernetes.io/instance: NAME-db
  ports:
    - port: 5432
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: NAME-db
  namespace: NAMESPACE
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app.kubernetes.io/instance: NAME-db
  template:
    metadata:
      labels:
        app.kubernetes.io/instance: NAME-db
    spec:
      containers:
        - name: postgres
          image: postgres:17
          env:
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef: {name: NAME-database, key: user}
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef: {name: NAME-database, key: password}
            - name: POSTGRES_DB
              valueFrom:
                secretKeyRef: {name: NAME-database, key: database}
          volumeMounts:
            - {name: data, mountPath: /var/lib/postgresql/data}
      volumes:
        - name: data
          persistentVolumeClaim: {claimName: NAME-db}
`)
}
