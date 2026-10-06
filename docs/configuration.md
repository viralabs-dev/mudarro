English · [Português (Brasil)](pt-BR/configuration.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Configuration

Init creates mudarro.yaml without overwriting existing config. --interactive selects suggestions, --select service:script,... selects noninteractively and --all-scripts selects all. Detected start/dev/build/test/lint scripts become known actions; other scripts are suggestions. Scanning never executes them.

## Version 1 contract

| Field | Meaning |
|---|---|
| version | Must be 1 |
| name | Banner name; init uses root manifest or directory name |
| exclude | Relative paths/filepath.Match patterns; ** is not recursive |
| services[].id | Unique letters/digits/_/- identifier |
| dir | Relative project directory |
| language, manager, framework | Explicit technology choices |
| infrastructure | kind, mode, file, context, namespace, image, port, generate |
| database | kind, tool, url_env, path, generate |
| commands | Actions with args OR shell, group, requires, destructive |
| pending | Editable scan observations; other actions remain usable |

File/database.path are service-relative. Explicit commands override adapter actions. npm script identifiers normalize :/. to -, keeping original argv; collisions require an explicit command. Args do not use shell expansion; ${VAR} requires an environment variable and {name} uses --name for migration creation. No automatic .env loading. Explicit shell runs Bash; avoid versioned secrets.

Generate creates missing Dockerfile/Compose/environment/Prisma scaffolds. db-install explicitly installs Prisma 7; review existing versions. Existing manifests/ORM versions are not automatically migrated. Configure POSTGRES_DB/USER/PASSWORD and DATABASE_URL. Compose DB host is db; host-side migrations use localhost/POSTGRES_PORT. Python container servers need explicit 0.0.0.0 binding. Database actions run in the service directory; override db-migrate for container execution.

Kubernetes context/namespace must exist and images must be available. Generation provides Deployment/Service, optional PostgreSQL PVC and references to an existing Secret without values. Down scales declared workloads to zero; DaemonSets require custom stop.

SQLite uses database.path (default app.db); db-init preserves an existing file. Prisma/Alembic URL must match. Local PostgreSQL with database.generate:true offers db-init/up/status/create/down, requires installed PostgreSQL tools and POSTGRES_* variables, and stores data in service .mudarro/postgres.

Prisma/Django use existing schema/models; Alembic revisions/target_metadata and Goose SQL remain project responsibilities. No business entities invented. Import mudarro_database.py explicitly into Django settings. Implement Prisma/Alembic/Goose seeds; initial Django fixture is empty. Destructive actions require the literal `APAGAR <service-id>`; db-reset always requires confirmation, even overridden.

## Existing reproduction examples

Commands and configuration identifiers are preserved verbatim; Portuguese comments and user-supplied example values are intentionally retained.

```yaml
version: 1
name: Meu serviço
services:
  - id: api
    dir: .
    language: typescript
    manager: npm
    infrastructure:
      kind: docker              # podman também é aceito
      mode: compose
      file: compose.yaml
      port: 3000
      generate: true
    database:
      kind: postgresql
      tool: prisma
      url_env: DATABASE_URL
      generate: true
    commands:
      start:
        args: [npm, run, dev]
        group: aplicacao
```

```yaml
commands:
  db-migrate:
    args: [docker, compose, -f, compose.yaml, exec, -T, api, npx, --no-install, prisma, migrate, deploy]
    group: banco
```

```yaml
infrastructure:
  kind: kubernetes
  mode: manifests              # ou kustomize; helm integra chart existente
  file: k8s
  context: kind-meu-cluster
  namespace: desenvolvimento
  image: minha-org/minha-app:v1
  port: 3000
  generate: true
```

```yaml
infrastructure:
  kind: custom
commands:
  up:
    args: [make, up]
    group: infraestrutura
  limpar-dados:
    shell: ./scripts/reset-db.sh
    requires: [bash]
    group: banco
    destructive: true
```

## Consumer project identity — MUD-022

The menu title is the consumer project `Config.Name`, not the Mudarro tool name. The controlled sample now uses Aurora (package/config); demo and recording default explicitly to Aurora, overridable with `MUDARRO_SAMPLE_NAME`. The existing CLI already used Config.Name; demonstration inputs were corrected. User visual design acceptance has been received; only the genuine graphical screenshot remains pending.

Text sanitization strips controls and Unicode Cf formatting characters. Conservative terminal-cell estimation counts CJK/fullwidth/emoji as two cells and combining marks as zero; it is not a complete grapheme/terminal-width guarantee. Unsupported bitmap glyphs or an excessively long wordmark fall back to readable text instead of presenting an incorrect identity.

The current name stage has 43 Test functions and 24 rich/narrow/plain/NO_COLOR name-mode cases, covering AtlasAPI, Aurora, Café, 漢字😀, long names and control injection, plus explicit consumer-name CLI integration. Race passed; the final coverage and build gates passed as recorded below. Previous stage coverage remains historical. Genuine latest media belongs in `evidence/project-name-visual/` under the canonical sample; `evidence/packages-visual/` remains historical.

## Current consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime was not exercised in this historical local stage; current CI runtime failed as recorded above.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Current implementation: [UI configuration](ui-configuration.md). Core configuration/i18n/theme/action-preview implemented; fonts and actual-emulator mouse acceptance remain separate pending work. Final gates:74 tests,76.38% aggregate statements; see [coverage](coverage.md).
