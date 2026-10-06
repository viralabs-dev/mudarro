English · [Português (Brasil)](pt-BR/support-matrix.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

## Current configurable-UI gate

**74 Test functions passed with race; 1,552/2,032 executed statements = 76.38% (Go displays76.4%)**, 480 unexecuted. Terminal package:379/422=89.81%. Vet, Linux build and whole-binary Darwin arm64 cross-build passed; compilation does not validate macOS runtime. Genuine UI E2E passed custom-config wrapper/supervisor propagation, localized scan, read-only preview and secret redaction, plus smoke/menu execution. Final Python/Go language and pnpm real E2Es passed; Dockerfile and eight database combinations with existing dependencies also passed. Previous stage counts/coverage below remain historical with different denominators. Font integration and actual-emulator mouse acceptance remain pending.

Evidence: [tests-race.txt](samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Support inventory and executed matrix

Code inventory reconciled after package extraction, 2026-10-06. Manifest detection and contract tests do not equal complete runtime support.

| Family | Detection/actions | Limits |
|---|---|---|
| JavaScript/TypeScript | package.json; TS config/dependency; install, application dev/start/build, quality test/lint, selectable scripts | npm default; npm/pnpm/yarn single lock or explicit packageManager; arbitrary manager accepted; Bun lock unsupported; next/Nest/vite/express labels only |
| Python | pyproject/requirements/manage.py/Pipfile; pip venv/install, uv sync, poetry install; entrypoint suggestions; Django start/test | uv/poetry by locks; conflicting locks pending; Pipfile does not implement pipenv; non-Django start explicit; no automatic Flask/FastAPI/pytest actions |
| Go | go.mod; main package parser; build/test/mod download; unique main directory start | multiple/no mains pending; go.work evidence only; build tags not evaluated |
| Custom | configured commands, Makefile/shell suggestions | not an automatic language adapter |

JS :/. identifiers normalize to -; collisions now produce pending choices and cannot save ambiguous selection, demonstrated before/after regression. Recursive directory scan respects exclusions/ignored dirs/symlinks and deterministic ordering. JS workspace root/children can become separate services; no workspace/Turbo/Nx resolution or root manager/infra inheritance. Nested Go modules are detected separately; no aggregate go.work execution. Shell/Make suggestions use nearest ancestor owner, fallback index0 outside services.

Infra: local supervisor, Docker/Podman Compose/Dockerfile and Kubernetes manifests/kustomize/Helm. Compose requires engine choice, multiple candidates remain pending, Compose precedes Dockerfile. Kubernetes requires context/namespace. Infra discovery is service-local, not arbitrary nested chart discovery. Podman uses podman-compose. Down preserves storage objects; persistence requires separate data tests.

Databases: PostgreSQL/SQLite; Prisma/Django/Alembic/Goose. Static evidence: schema.prisma/manage.py/alembic.ini/Goose SQL annotation. Only Prisma strings infer kind; other kinds explicit. Multiple tools pending. Migrations/create/seed separate; resets destructive where provided; Alembic has no automatic reset. No invented models.

## Executed package round

| Scenario | Level | Result/limit |
|---|---|---|
| JS/TS npm/pnpm/yarn, Python pip/uv/poetry, Go workspace/multiple main | Detection/contracts | Passed; manager overrides/conflicts, no invented Python start; Go build tags remain limited |
| Seven family/manager combinations | in-process scan/init/generate/menu/dry-run | Passed argv/hash/idempotence; managers not executed |
| Docker/Podman Compose/Dockerfile, Kubernetes three modes | Plans/fake executor | Passed provider/context/ownership/workloads/PVC/DaemonSet/invalid JSON |
| Invalid/large/symlink/excluded manifests, selection/preflight | Negative fixtures | Passed; valid sibling preserved, no partial writes |
| JS normalized collisions | Before/after regression | Failed baseline, passed fix; ambiguous selection does not save |
| Node/npm, Python/Go, existing pnpm | Real E2E | Passed menus/wrappers/supervisor, no dependency install |
| Eight PostgreSQL/SQLite × tool combinations | Real runtime, existing deps | Passed migration/seed/repetition; Goose/Prisma hooks do not prove row insertion |
| Compose/Dockerfile | Real runtime | Passed isolated ownership/lifecycle |
| Kustomize/Helm | Real runtime | Passed mounted Bound PVC, same UID/content |
| Yarn/uv/poetry/Podman runtime | Host tools absent | Not executed; contract tests are not runtime |
| macOS/WSL | Platform unavailable | Not executed; Darwin cross-build only |

Evidence: packages-* in samples/menu-auto/evidence. Aggregate historical profile 72.0% with coverpkg, not individual test-binary percentages. Infra/database integrations used extraction snapshot b21c91f; subsequent JS collision/named-field changes had final suites and Node/Python/Go/pnpm E2Es rerun. The final test-tree gate is recorded below.

## Acceptance for extensions

Valid/invalid/absent/ambiguous manifests, no subprocess/network during scan, root/subdirectory/monorepo, exclusions/symlinks, predictable argv, explicit startup, valid/collision-safe names, selection and repeat generation with hashes, isolated real launcher, preserved overrides/errors/exit codes. Runtime tests require available authorized tools, isolated caches/ports/resources and no resets against user databases. See [coverage](coverage.md) and [roadmap](language-roadmap.md).

## Final test-tree gate — 2026-10-06

All 41 tests passed with race and explicit internal coverpkg instrumentation after relocation. Deduplicated coverage: **1,035/1,427 statements = 72.53% (Go displays 72.5%)**; previous 1,028/1,427=72.04% remains historical. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and diff check passed. Darwin compilation is not macOS runtime. Test tree contains 16 files (nine orchestration test files plus one helper; two terminal test files plus four PTY helpers). No test hooks or overlays added to production.

Evidence: [test-tree-tests.txt](samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](samples/menu-auto/evidence/test-tree-gates.txt).

## Current consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime was not exercised in this historical local stage; current CI runtime failed as recorded above.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
