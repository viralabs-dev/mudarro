# Local robustness — MUD-050 / MUD-051 / MUD-052

Base main `56dadbaf62e3515e4bb31bb3c3ff2a583cd642ed`; this checkpoint is local and uncommitted. No installation, Actions run, commit or push in this round.

| Scope | Before | Final focused result | Evidence |
|---|---|---|---|
| Supervisor identity | Corrupt JSON and forged project state failed; separate forged argv[0] failed after first fix | 9 groups passed with race | [before](supervisor/supervisor-before.txt), [spoof before](supervisor/spoof-before.txt), [final](supervisor/supervisor-final.txt) |
| IO/preservation | Writer failure ignored and FIFO read blocked (owned helper terminated/joined) | 13 groups / 6 subcases passed with race | [before](io/io-before.log), [final](io/io-final-expanded.log) |
| Doctor negatives | No new production defect reproduced | 4 groups passed with race | [focused](doctor/doctor-focused.txt) |

Integrated race: **112 top-level tests passed, zero failures/skips**, including optional real PTY recording. Deduplicated internal statement coverage **1,943 / 2,510 = 77.41%**. Repeated coverpkg blocks are merged by source range and covered if any binary executed them. Supervisor subprocesses use an ordinary built CLI: their behavior is real E2E evidence, not instrumentation of child functions. [Summary](summary.json), [JSON test events](tests.jsonl), [profile](coverage.out), [functions](coverage-functions.txt), [PTY cast](interactive.cast).

Vet, Linux build, Darwin arm64 binary/test cross-compilation passed. [Smoke](smoke.txt), [automatic generation/idempotence/manual-edit preservation](menu.txt), and [selected-config supervisor/wrappers](ui.txt) passed on final Linux binary SHA256 `8350359bafb6d4591efbdf347311acd6abef54b46fb8f42e804867ef3b97183b`; see [exit results](e2e-results.json). Temporary processes/fixtures belong only to these tests.

Linux ownership now compares actual executable inode and exact NUL-separated argv (binary/root/service/selected config/token). JSON decoding errors invalidate state. `/proc` availability is required for Linux ownership detection; this is not protection against a hostile process modifying its own address space or PID reuse between check and signal. Darwin fallback compares a full `ps` command string; argument-boundary/inode guarantees and native runtime are **not validated** (MUD-031). The Linux inode regression explicitly skips Darwin for that documented platform-specific mechanism.

Generation checks output errors before mutations and bounds ownership manifests to regular files of at most 4 MiB. Manifest reads reject static FIFO/directory/symlink fixtures and enforce a bounded read. Writes are atomic **per file**, not a batch transaction: earlier writes may remain on a later IO failure. SafePath/stat/open checks do not provide an atomic filesystem boundary against concurrent replacement; a FIFO swapped after stat can still race opening. No injected abstractions solely for coverage. Doctor checks executables/provider version, never application/database health; Compose provider tests are controlled stubs, not engine runtime.

Reproduce from repository root with existing Go:

```bash
export GOCACHE=/tmp/mudarro-go-cache
export MUDARRO_INTERACTIVE_CAST=/tmp/mudarro-robustness-interactive.cast
/home/danielsouza/sdk/go1.27.1/bin/go test -race -coverpkg=./internal/mudarro/... -coverprofile=/tmp/mudarro-robustness-coverage.out -json ./...
/home/danielsouza/sdk/go1.27.1/bin/go vet ./...
/home/danielsouza/sdk/go1.27.1/bin/go build -o /tmp/mudarro-robustness-cli ./cmd/mudarro
bash scripts/smoke.sh /tmp/mudarro-robustness-cli
bash scripts/validate-menu.sh /tmp/mudarro-robustness-cli
bash scripts/validate-ui.sh /tmp/mudarro-robustness-cli
```

Existing graphical/macOS/mouse/font blockers remain separate. No browser, desktop window or new screenshot was used. Prior genuine dark/light recordings and all historical evidence are preserved.
