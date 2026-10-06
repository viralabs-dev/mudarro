# Pipenv and Go — partial local queue checkpoint

Base main `b887f64`, new code/tests/docs remain uncommitted. Two independent agents owned new test files and exclusive fixtures; root integrated shared code sequentially and was the only Vault writer. No installs, dependency additions, Actions, publication, browser or desktop capture.

Integrated final: **131 top-level race tests PASS, 0 FAIL, 0 SKIP; 2,041/2,588 statements = 78.86%**, deduplicating coverage source blocks. [Summary](summary.json), [events](tests.jsonl), [coverage](coverage.out), [functions](coverage-functions.txt), [genuine PTY](interactive.cast). Vet, Linux build, Darwin arm64 cross-build and smoke/menu/config/resize NO_COLOR/SIGTERM passed. Native macOS not executed. Final Linux binary SHA256 `2d3b3408e5af427558c1d318eda1d3e716f540eeca663e921a96d9d7325b88d9`.

MUD-034: valid Pipfile selects Pipenv; lock present proposes pipenv sync, absent proposes pipenv install. Scan executes neither. TOML/JSON parsing errors and orphan lock are actionable. Conflicting uv/Poetry locks keep manager unresolved. String scripts are opt-in suggestions by original key; unsupported forms/collisions stay pending. Django and Alembic use pipenv run python for database commands; dependency actions use pipenv install. [Baseline](pipenv/before-evidence.txt), [database baseline](pipenv/database-before-evidence.txt), [seven groups/nine subcases after](pipenv/final-after-evidence.txt), [contract](pipenv/contract.md). **Runtime not executed**: [Pipenv absent from PATH and six known locations](pipenv/runtime-triage.json). No full Pipfile/lock schema or hash synchronization validation is claimed. Card remains incomplete pending actual runtime, without installing tools.

MUD-038: infer Go start only for an active host build file with a valid func main declaration. Ignore build-excluded files and nested-module boundaries. This static check does not prove import resolution, compilation or package consistency, and does not infer custom build tags. go.work emits an explicit limitation warning even when valid; **membership, execution context and target selection are still unimplemented**. Existing module discovery remains independent. [Baseline](go/baseline.txt), [three groups/twelve subcases after](go/final-focused.txt), [contract](go/contract.txt). Genuine offline Go workspace fixture passed scan/init, generation hashes, build/test and supervisor up/status/logs/restart/down on the Go-fix checkpoint before the later Python-only database helper change: [runtime](go/runtime-final.json), [cleanup](go/cleanup.json). Initial environmental VCS stamping failure preserved [here](go/runtime-vcs-limitation.json); rerun scoped GOFLAGS=-buildvcs=false to the subprocess without changing Git settings. Workspace support remains partial.

Reproduce from repo root with existing Go:

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run '^(TestPipenvSupport|TestGoWorkspace)' -v
GOCACHE=/tmp/mudarro-go-cache MUDARRO_INTERACTIVE_CAST=/tmp/mudarro-language-repro.cast /home/danielsouza/sdk/go1.27.1/bin/go test -race -coverpkg=./internal/mudarro/... -coverprofile=/tmp/mudarro-language-repro.cover ./...
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go vet ./...
```

Sources: [Pipenv manifest](https://pipenv.pypa.io/en/latest/pipfile.html), [Go workspace reference](https://go.dev/ref/mod#workspaces). Historical Bun, robustness and Unicode evidence remains unchanged in neighboring directories. MUD-035 remains analysis-only with fourteen CLI cases rechecked; no manager/version validation fix yet.
