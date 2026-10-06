English · [Português (Brasil)](pt-BR/validation.md)

Published base is now `b887f64` after the authorized personal-attribution rewrite of `56dadba`, with an identical tree. [Old/new SHA map and preservation proof](samples/menu-auto/evidence/git-identity-repair/README.md). New robustness/Unicode changes and this reference update remain local.

## Latest integrated checkpoint — official Go workspace parser

165 race PASS (162 substantive groups + three helpers), 0 FAIL/SKIP; **2199/2746 = 80.08%**. Official go.work/go.mod parsing, membership/conflict/exclusion diagnostics, inherited context and explicit isolation pass. All final local gates and genuine Go/five installed runtime rechecks pass. Four cross-built archives and installer fixtures preserve licenses. Native macOS and physical graphical acceptance remain unexecuted. [Evidence / reprodução](samples/menu-auto/evidence/go-work-parser/README.md).

## Historical checkpoint156 — confirmed manager/Go decisions

156 race PASS (154 substantive groups + two helpers), 0 FAIL/SKIP; **79.47%**. Exact SemVer syntax, explicit Go start selection and scoped workspace off pass. Preview/dry-run/doctor show/check isolation. All local gates and current-binary scans/quality for five installed runtimes passed. Full go.work membership parser remains pending; native macOS not executed. [Evidence / reprodução](samples/menu-auto/evidence/manager-go-decisions/README.md).

## Latest local runtime validation — isolated managers

Approved pinned Pipenv2026.8.0, uv0.12.23, Poetry2.5.1 and YarnClassic1.22.22/Modern4.18.1 now have genuine zero-dependency fixture validation. Real locks, generated menus/wrappers, idempotence, preservation negatives and lifecycle passed; Python start is explicit. Installation and media proofs are archived. Production code did not change; the prior143 race tests/79.08% remain the last measured code checkpoint. [Evidence / reprodução](samples/menu-auto/evidence/isolated-runtimes/README.md).

## Latest local checkpoint — Compose doctor timeout

143 race PASS (141 substantive groups + two helpers), 0 FAIL/SKIP; **79.08%**. A genuine Compose probe timed out after 10.05s; configuration and direct process cleanup verified. All local gates passed; native macOS remains unexecuted. [Evidence / reprodução](samples/menu-auto/evidence/doctor-compose-timeout/README.md).

## Previous local checkpoint — doctor special files

141 top-level race PASS (140 substantive groups + one helper), 0 FAIL/SKIP; **79.03%**. Doctor now rejects relative special files with executable bits; genuine before/after CLI preserved fixtures. Vet/builds/smoke/menu/config/resize passed. Native macOS not executed; new changes local. [Evidence / reprodução](samples/menu-auto/evidence/doctor-special/README.md).


## Previous local checkpoint — packageManager names / Go source IO

140 top-level race PASS (139 substantive groups + one subprocess helper), 0 FAIL/SKIP; **2,058/2,604 = 79.03%**. Known-name validation and bounded regular Go-source reading passed; vet/builds/smoke/menu/config/resize passed. Version grammar, Go target selection/context await decision; UV/Poetry runtime absent. Changes local. [Evidence / reprodução](samples/menu-auto/evidence/manager-go-source/README.md).


## Previous local checkpoint — partial MUD-034/038

**131 top-level race tests PASS, 0 FAIL/SKIP; 2,041/2,588 = 78.86%.** Pipenv detection/scripts/database argv and Go host entrypoint regressions pass. Vet/Linux/Darwin cross-build and smoke/menu/config/resize pass. Pipenv runtime is unavailable; Go workspace membership/context/target selection remains unimplemented, with an explicit warning. [Evidence and reproduction](samples/menu-auto/evidence/pipenv-go/README.md). Changes remain local.


## Previous local checkpoint — MUD-033

121 top-level race tests PASS, zero failures/skips; 1,971/2,517 internal statements = **78.31%**. Bun 1.4.0 genuine runtime and menu/wrappers/supervisor passed without installation; bun.lock/bun.lockb detection and cross-manager ambiguity regressions pass. Vet, Linux build, Darwin cross-build and smoke/menu/config/PTY resize gates passed. [Evidence, saved sample and reproduction](samples/menu-auto/evidence/bun/README.md). Native macOS and graphical/physical terminal acceptance remain pending. All new changes are local on b887f64. MUD-035 remains analysis-only.


## Previous local checkpoint — MUD-053

117 top-level race tests passed,zero failures/skips;deduplicated coverage1,965/2,515=78.13%. U+3000 now occupies2cells;clipping respects1–3columns without expanding the budget.17variants and actual PTY resize/NO_COLOR/restoration passed,along with vet/Linuxbuild/Darwincrosscompile and smoke/menu/config/resize E2Es. [Evidence,reproduction and limits](samples/menu-auto/evidence/unicode-resize/README.md). Conservative rune-cell policy,without perfect shaping/indivisible graphemes;physical acceptance026/027 separate,native macOS031pending. Local changes on56dadba,no new commit/push/CI;history056 completed separately with personal attribution. [035analysis-onlytriage](samples/menu-auto/evidence/package-manager-contracts/README.md).

## Previous local robustness checkpoint — MUD-050/051/052

112 top-level tests passed with race, zero failures/skips; fresh internal statement coverage 1,943/2,510 = **77.41%**. Linux supervisor checks actual executable inode and exact argv; corrupt state is rejected. Generation propagates output errors before writes, and manifest reads reject static special files and enforce 4 MiB. Doctor negative checks passed without running application actions. Vet/Linux build, Darwin cross-compilation, genuine smoke/menu idempotence/preservation and selected-config E2Es passed. [Reproduction, logs and limits](samples/menu-auto/evidence/robustness/README.md).

This is local work on published base56dadba, not committed/pushed or validated by CI. Native macOS remains unvalidated; its ps fallback has weaker identity guarantees. IO is atomic per file, without batch rollback or atomic protection against concurrent path replacement. Historical stage measurements below retain their original denominators.

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

## Historical configurable-UI gate

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

## Historical consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime was not exercised in this historical local stage; current CI runtime failed as recorded above.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
