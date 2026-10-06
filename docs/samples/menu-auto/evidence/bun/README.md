# Bun — MUD-033 (local)

Published base `b887f64`; all new robustness, Unicode and Bun changes remain uncommitted. No installation, browser, desktop capture, Actions or further publication was performed.

Detection recognizes `bun.lock` and legacy `bun.lockb`. Both files count as one manager; locks belonging to different managers retain the existing Pending behavior unless an explicit packageManager selects one. Scan never invokes Bun and does not migrate locks. [Official Bun lockfile reference](https://bun.sh/docs/pm/lockfile).

[Before](before.log): locks were ignored, conflicts silently selected another manager and init selected npm. [After](after.log): four top-level regression groups, six subcases passed with race, covering both formats, deduplication, npm/pnpm/yarn conflicts, explicit Bun override, init and generation idempotence.

Integrated: **121 top-level race tests PASS, 0 FAIL, 0 SKIP; 1,971/2,517 statements = 78.31%**, deduplicating repeated coverpkg blocks. [Summary](summary.json), [events](tests.jsonl), [coverage](coverage.out), [functions](coverage-functions.txt), [genuine integrated PTY](interactive.cast). Vet, Linux build and Darwin arm64 cross-build passed; native macOS was not executed. Smoke, menu and selected-config E2Es and PTY resize/NO_COLOR/SIGTERM passed. Binary SHA256: `6214994c4e1dfcc06acfbd7beae14e9b88e592b7e25049fbb06a55e3a8c12033`.

Existing Bun **1.4.0** executed genuine scan/init, dry-run, two identical generate passes, doctor, generated test/lint/build wrappers, menu check, up twice, status/logs/restart/down/stopped. [Runtime log](runtime/e2e.txt), [menu execution](runtime/menu.txt), [real menu PTY cast](runtime/menu.cast), [owned-process cleanup](runtime/process-cleanup.json). This fixture uses explicit packageManager; lock-only detection is covered separately by regression tests. Runtime required no dependencies or install; [saved sample](sample/package.json). The cast records an actual headless PTY with injected exit key, not a graphical screenshot. Existing dark/light PNG/GIF recordings remain in [ui-shell](../ui-shell/).

From repo root, using existing tools:

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run '^TestBun' -v
GOCACHE=/tmp/mudarro-go-cache MUDARRO_INTERACTIVE_CAST=/tmp/mudarro-bun-repro.cast /home/danielsouza/sdk/go1.27.1/bin/go test -race -coverpkg=./internal/mudarro/... -coverprofile=/tmp/mudarro-bun-repro.cover ./...
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go build -o /tmp/mudarro-bun-repro ./cmd/mudarro
bash scripts/smoke.sh /tmp/mudarro-bun-repro
bash scripts/validate-menu.sh /tmp/mudarro-bun-repro
bash scripts/validate-ui.sh /tmp/mudarro-bun-repro
python3 scripts/validate-ui-shell-resize.py /tmp/mudarro-bun-repro /tmp/mudarro-bun-resize-repro
```

Runtime helper scripts preserve the original exclusive `/tmp/mudarro-mud033-runtime` paths and capture recipe. For independent reproduction, copy `sample` to a fresh directory, build the binary, put its directory and existing Bun on PATH, then run scan/init/generate/doctor and generated wrappers. Do not execute the archived helper against an unrelated directory. packageManager validation MUD-035 remains analysis-only; unknown managers/version grammar are not corrected by this change. Supervisor, manifest and terminal limitations remain documented in the preceding robustness/Unicode evidence.
