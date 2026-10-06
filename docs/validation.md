English · [Português (Brasil)](pt-BR/validation.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

## Current configurable-UI gate

**74 Test functions passed with race; 1,552/2,032 executed statements = 76.38% (Go displays76.4%)**, 480 unexecuted. Terminal package:379/422=89.81%. Vet, Linux build and whole-binary Darwin arm64 cross-build passed; compilation does not validate macOS runtime. Genuine UI E2E passed custom-config wrapper/supervisor propagation, localized scan, read-only preview and secret redaction, plus smoke/menu execution. Final Python/Go language and pnpm real E2Es passed; Dockerfile and eight database combinations with existing dependencies also passed. Previous stage counts/coverage below remain historical with different denominators. Font integration and actual-emulator mouse acceptance remain pending.

Evidence: [tests-race.txt](samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Validation

Historical initial measurements were recorded for 2026-10-05; expanded local round on 2026-10-06 used main base 049a405 with uncommitted changes. Results distinguish automated contracts, genuine runtime execution and unavailable platforms.

Go/race/vet, Linux build, Darwin arm64 cross-build, shell syntax, installer simulation, Node/npm/pnpm/Python/Go menu/wrappers/local lifecycle passed. Aggregate recorded package coverage: 1,028/1,427 statements (72.04%, Go displays72.0%). Earlier root62.8% has a different denominator. See [coverage](coverage.md) and [support matrix](support-matrix.md). The separate final migration gate is recorded below.

Docker Compose/Dockerfile, offline --network none, kind manifests/kustomize and Helm passed isolated runtime tests. Eight PostgreSQL/SQLite combinations with Goose/Alembic/Django/Prisma were rerun using existing dependencies and passed. Django/Alembic tested inserted/idempotent data; Goose/Prisma tested configured seed hooks. Initial PVC fixture Pending proved object preservation only; expanded mounted Bound fixture proved same UID/content after down/up/restart. This is fixture-cycle persistence, not backup/disaster recovery or retention after deleting a cluster.

Native host PostgreSQL, Podman, macOS and WSL runtime were not exercised locally: unavailable binaries/platforms. Current CI macOS executed and failed; container CI passed. Yarn/uv/poetry runtime also unavailable; contract tests remain separate. Doctor is executable/provider checking, not application/database health. Prisma7 SQLite requires initial file creation, covered by db-init. Models remain project-owned; container migrations require explicit commands when tools are not on host.

Tools previously exercised locally: Go1.27.1, Node26.7.0, Python3.14.7, Prisma7.10.0, Alembic1.20.0, Django6.1.1, Goose3.24.1, PostgreSQL17. CI versions may differ.

Historical CI [37405606285](https://github.com/viralabs-dev/mudarro/actions/runs/37405606285) succeeded at f62ce9d and does not validate new local changes. Existing public [v0.1.0](https://github.com/viralabs-dev/mudarro/releases/tag/v0.1.0) was previously published with four binaries/checksums and installation verified. This historical local round made no commit/push/publication; the subsequent published checkpoint and its failing CI are recorded above.

Final recorded v4 snapshot: /tmp/mudarro-checkpoint-v4/bin/mudarro, SHA25676ccc665970c48a7f7c95f9a9052599cd23aee4580060b11a265a5b1bd91d942. Genuine terminal GIF/composed frames validated in [sample](samples/menu-auto/README.md); graphical screenshot pending; user visual design acceptance received, existing workspace4/notebook confirmation preserved.

## Existing reproduction examples

Commands and configuration identifiers are preserved verbatim; Portuguese comments and user-supplied example values are intentionally retained.

```bash
go test -race ./...
go vet ./...
go build -o bin/mudarro ./cmd/mudarro
bash scripts/smoke.sh "$PWD/bin/mudarro"
bash scripts/test-installer.sh
bash scripts/integration-containers.sh "$PWD/bin/mudarro" docker
bash scripts/integration-containers.sh "$PWD/bin/mudarro" podman
bash scripts/integration-offline.sh "$PWD/bin/mudarro"
bash scripts/validate-kind.sh "$PWD/bin/mudarro"
go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
bash scripts/validate-databases.sh "$PWD/bin/mudarro"
```

## Final test-tree gate — 2026-10-06

All 41 tests passed with race and explicit internal coverpkg instrumentation after relocation. Deduplicated coverage: **1,035/1,427 statements = 72.53% (Go displays 72.5%)**; previous 1,028/1,427=72.04% remains historical. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and diff check passed. Darwin compilation is not macOS runtime. Test tree contains 16 files (nine orchestration test files plus one helper; two terminal test files plus four PTY helpers). No test hooks or overlays added to production.

Evidence: [test-tree-tests.txt](samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](samples/menu-auto/evidence/test-tree-gates.txt).

Plain go test also passed: [test-tree-plain-tests.txt](samples/menu-auto/evidence/test-tree-plain-tests.txt).

## Current consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime was not exercised in this historical local stage; current CI runtime failed as recorded above.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
