"""Disposable integration projects; downloads only happen when explicitly running this script."""
from contextlib import nullcontext
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile

BIN = str(Path(sys.argv[1]).resolve())


def run(args, cwd, env=None, capture=False):
    return subprocess.run(args, cwd=cwd, env=env, check=True, text=True,
                          stdout=subprocess.PIPE if capture else None).stdout


with (nullcontext(os.environ["MUDARRO_TEST_ROOT"]) if os.environ.get("MUDARRO_TEST_ROOT") else tempfile.TemporaryDirectory(prefix="mudarro-databases-")) as temp:
    base = Path(temp)
    # Reuse mode never installs; both paths must point to existing dependencies.
    reuse_python = os.environ.get("MUDARRO_PYTHON_ENV")
    reuse_js = os.environ.get("MUDARRO_JS_DEPS")
    if bool(reuse_python) != bool(reuse_js):
        raise SystemExit("Configure MUDARRO_PYTHON_ENV e MUDARRO_JS_DEPS juntos")
    if not shutil.which("goose"):
        raise SystemExit("Goose ausente; nenhum pacote instalado. Forneça um Goose existente no PATH")
    if reuse_python:
        venv = Path(reuse_python).resolve()
        js = Path(reuse_js).resolve()
        python = str(venv / "bin/python")
        run([python, "-c", "import alembic, sqlalchemy, django, psycopg"], base)
        if not (js / "node_modules/prisma").is_dir():
            raise SystemExit("Prisma ausente nas dependências fornecidas; nenhuma instalação feita")
    else:
        venv = base / "python-env"
        run([sys.executable, "-m", "venv", str(venv)], base)
        python = str(venv / "bin/python")
        run([python, "-m", "pip", "install", "alembic", "sqlalchemy", "django", "psycopg[binary]"], base)
        js = base / "js-deps"
        js.mkdir()
        (js / "package.json").write_text('{"private":true}')
        run(["npm", "install", "--no-audit", "--no-fund", "prisma@7", "@prisma/client@7"], js)
        run(["node", "node_modules/@prisma/engines/scripts/postinstall.js"], js)
    kinds = ["sqlite"]
    if os.environ.get("TEST_POSTGRES_URL"):
        kinds.append("postgresql")
    for kind in kinds:
        for tool in os.environ.get("MUDARRO_TEST_TOOLS", "goose,alembic,django,prisma").split(","):
            root = base / f"{kind}-{tool}"
            root.mkdir()
            (root / ".venv").symlink_to(venv, target_is_directory=True)
            (root / "node_modules").symlink_to(js / "node_modules", target_is_directory=True)
            (root / "package.json").write_text('{"private":true}')
            config = {"version": 1, "name": "Database test", "services": [{
                "id": "app", "dir": ".", "language": "python" if tool in ("django", "alembic") else "go" if tool == "goose" else "javascript",
                "manager": "pip" if tool in ("django", "alembic") else "npm",
                "infrastructure": {"kind": "custom"},
                "database": {"kind": kind, "tool": tool, "generate": True, "path": "app.db"}}]}
            (root / "mudarro.yaml").write_text(json.dumps(config))
            run([BIN, "generate", "--root", str(root)], root)
            env = os.environ.copy()
            url = os.environ.get("TEST_POSTGRES_URL", "")
            if kind == "sqlite":
                url = {"prisma": "file:./app.db", "alembic": "sqlite:///app.db"}.get(tool, "app.db")
            elif tool == "alembic":
                url = url.replace("postgresql://", "postgresql+psycopg://")
            elif tool == "prisma":
                url += "?schema=mudarro_prisma"
            env["DATABASE_URL"] = url

            def action(name, *args):
                return run([BIN, "run", "app:" + name, "--root", str(root), *args], root, env)

            if kind == "sqlite":
                action("db-init")
            if tool == "goose":
                action("db-migration-new", "--name", "probe")
                migration = next((root / "migrations").glob("*.sql"))
                migration.write_text("-- +goose Up\nCREATE TABLE goose_probe (id INTEGER PRIMARY KEY);\n-- +goose Down\nDROP TABLE goose_probe;\n")
                action("db-migrate")
                # Exercise a user-provided idempotent seed through the generated action.
                (root / "seeds/seed.sh").write_text("#!/usr/bin/env bash\nset -euo pipefail\nprintf 'goose seed hook invoked\\n'\n")
                action("db-seed")
            elif tool == "alembic":
                action("db-migration-new", "--name", "probe")
                migration = next((root / "migrations/versions").glob("*.py"))
                text = migration.read_text().replace("def upgrade():\n    pass", "def upgrade():\n    op.create_table('alembic_probe', sa.Column('id', sa.Integer(), primary_key=True))")
                migration.write_text(text)
                action("db-migrate")
                (root / "seeds/seed.py").write_text("import os\nfrom sqlalchemy import create_engine, text\nwith create_engine(os.environ['DATABASE_URL']).begin() as c:\n c.execute(text('INSERT INTO alembic_probe (id) VALUES (1) ON CONFLICT DO NOTHING'))\n")
                action("db-seed")
                action("db-seed")
            elif tool == "django":
                (root / "manage.py").write_text("import os, sys\nos.environ.setdefault('DJANGO_SETTINGS_MODULE','settings')\nfrom django.core.management import execute_from_command_line\nexecute_from_command_line(sys.argv)\n")
                (root / "settings.py").write_text("from mudarro_database import DATABASES\nSECRET_KEY='integration-test-only'\nINSTALLED_APPS=['probe']\nDEFAULT_AUTO_FIELD='django.db.models.AutoField'\n")
                (root / "probe/migrations").mkdir(parents=True)
                (root / "probe/__init__.py").write_text("")
                (root / "probe/migrations/__init__.py").write_text("")
                (root / "probe/models.py").write_text("from django.db import models\nclass Record(models.Model):\n name = models.CharField(max_length=40)\n")
                action("db-migration-new", "--name", "probe")
                action("db-migrate")
                (root / "seeds/initial.json").write_text('[{"model":"probe.record","pk":1,"fields":{"name":"seed"}}]')
                action("db-seed")
                action("db-seed")
            else:
                schema = root / "prisma/schema.prisma"
                schema.write_text(schema.read_text() + '\nmodel PrismaProbe {\n id Int @id\n}\n')
                action("db-migration-new", "--name", "probe")
                action("db-migrate")
                (root / "prisma/seed.cjs").write_text("console.log('prisma seed hook invoked');\n")
                action("db-seed")
            action("db-migrate")
            print(f"PASS {kind}/{tool}: migration + seed + repeated migrate", flush=True)
