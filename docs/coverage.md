English · [Português (Brasil)](pt-BR/coverage.md)

## Current configurable-UI gate

**74 Test functions passed with race; 1,552/2,032 executed statements = 76.38% (Go displays76.4%)**, 480 unexecuted. Terminal package:379/422=89.81%. Vet, Linux build and whole-binary Darwin arm64 cross-build passed; compilation does not validate macOS runtime. Genuine UI E2E passed custom-config wrapper/supervisor propagation, localized scan, read-only preview and secret redaction, plus smoke/menu execution. Final Python/Go language and pnpm real E2Es passed; Dockerfile and eight database combinations with existing dependencies also passed. Previous stage counts/coverage below remain historical with different denominators. Font integration and actual-emulator mouse acceptance remain pending.

Evidence: [tests-race.txt](samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Interpreting measured coverage

Current measurement after test migration, 2026-10-06: **1,035/1,427 executed statements = 72.53%**, displayed as 72.5%; 392 statements unexecuted. All 41 tests passed.

Historical package measurement recorded on 2026-10-06: **1,028/1,427 executed statements = 72.04%**, displayed as 72.0%; 399 statements unexecuted. The earlier 71.8% checkpoint was superseded after the JS collision regression/fix. This measures statements, not lines, branches, requirements or platforms. The final test-tree gate below supersedes this historical measurement.

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race \
  -coverpkg=./internal/mudarro/... \
  -coverprofile=docs/samples/menu-auto/evidence/project-name-coverage.out -v ./...
/home/danielsouza/sdk/go1.27.1/bin/go tool cover \
  -func=docs/samples/menu-auto/evidence/project-name-coverage.out
```

Internal packages are instrumented; cmd/mudarro is outside coverpkg. Model has no executable statements. Atomic profile totals count NumStmt once per block, covered when Count>0; do not average package or test-binary percentages.

Historical MUD-017 package totals:

| Package | Covered / total | Percentage |
|---|---:|---:|
| Root | 589/935 | 62.99% |
| adapters | 279/304 | 91.78% |
| executor | 10/13 | 76.92% |
| projectfs | 19/23 | 82.61% |
| terminal | 138/152 | 90.79% |
| model | No executable statements | N/A |

Go tests include detection, contracts, fake execution, config/generation/preservation, in-process CLI and terminal. Race detects races only in exercised runs. External scripts execute ordinary uninstrumented binaries: Node/npm/pnpm/Python/Go, supervisor, containers/Kubernetes and database E2Es are independent evidence. OSExecutor tests cover parent code, not subprocess applications. A 0% local supervisor profile therefore does not negate its real E2E. Previous 62.8% described the old root package; the denominator changed after extraction and new tests.

## Priorities by risk

- Supervisor/identity/signals: stale/corrupt state, foreign PID/token/project refusal, startup failure, process trees and concurrency; supervisor functions were 0% in this profile despite passing E2Es.
- Doctor: absent/nonexecutable tools, blocked config, symlinks and exit status; doctor was 0% but exercised in available E2E environments.
- IO/execution failures: atomicWrite 53.3%, Output 40%, SafePath 85.7%, ReadManifest 77.8%; preserve files and errors under create/rename/read failures.
- TTY: final terminalSize 81.8%, Choose 97.3%, HasColor 100%; real isolated PTYs cover the rich path. Resize, wide Unicode and cancellation still require explicit acceptance. Legacy unused ascii/terminalColumns remain 0%.
- Templates/config: djangoDatabase 33.3%, scaffoldFiles 74.8%, Validate 71.0%; cover meaningful variants without inventing models.

No automatic 100% target. Coverage does not prove security or all managers/platforms. Yarn/uv/poetry/Podman and macOS/WSL runtime remain untested.

Evidence: [profile](samples/menu-auto/evidence/test-tree-coverage.out), [functions](samples/menu-auto/evidence/test-tree-coverage-functions.txt), [totals](samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [matrix](support-matrix.md).

## Final test-tree gate — 2026-10-06

All 41 tests passed with race and explicit internal coverpkg instrumentation after relocation. Deduplicated coverage: **1,035/1,427 statements = 72.53% (Go displays 72.5%)**; previous 1,028/1,427=72.04% remains historical. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and diff check passed. Darwin compilation is not macOS runtime. Test tree contains 16 files (nine orchestration test files plus one helper; two terminal test files plus four PTY helpers). No test hooks or overlays added to production.

Evidence: [test-tree-tests.txt](samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](samples/menu-auto/evidence/test-tree-gates.txt).

## Current consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime remains unexecuted.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Current terminal package totals: 164/176 statements (93.18%). Function profile: cleanText100%, fitText94.1%, runeColumns80%, wordmark95.7%, HasColor100%, terminalSize81.8%. Remaining width/grapheme/resize limits still apply.
