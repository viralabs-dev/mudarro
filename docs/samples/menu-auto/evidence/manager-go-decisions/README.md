# Confirmed manager/Go decisions — MUD-035 and partial MUD-038

**Historical checkpoint/proposal:** subsequent explicit parser authorization superseded the earlier dependency constraint. The official pinned parser is now implemented and verified; [current integrated result](../go-work-parser/README.md). Previous baselines and approval state below are retained as historical evidence.


Integrated local checkpoint: **156 race PASS (154 substantive groups + two isolated helpers), 0 FAIL/SKIP; 2114/2660 = 79.47%**. [Summary](summary.json), [events](tests.jsonl), [coverage](coverage.out), [genuine PTY](interactive.cast). Vet, Linux build, Darwin arm64 cross-build, smoke, menu/configuration and actual PTY resize passed. Native macOS not executed. SHA256 `eee2d8f973c282c5744a220f17686258f0a67f04b0cb970a325315c43fd72dde`; base b887f64, new work local with no new commit/push/Actions.

MUD-035: legacy known manager names without versions remain accepted. Supplied versions require exact SemVer2.0 syntax, including optional prerelease and build metadata. Tags, ranges, URLs, prefixes, malformed identifiers, whitespace and invalid leading zeros are rejected with a packageManager diagnostic; valid siblings and user/generated files remain preserved. Three groups/44 cases and seven actual CLI scans passed. [Before](semver/before.txt), [after](semver/after.txt), [actual CLI](semver/cli-after.json).

Corepack-style `+sha224.text` is valid build-metadata syntax, but scanning does not verify hash length/content/artifact integrity or enforce/download the named version. Execution still uses the selected manager on PATH. Explicit custom commands remain available. These are syntax/detection contracts, not dependency provisioning. [SemVer2.0](https://semver.org/) · [Corepack packageManager documentation](https://github.com/nodejs/corepack#when-authoring-packages).

MUD-038 delivered portions: active host Go entries become sorted `go-entry-1`, `go-entry-2`, … suggestions whose purpose is go-start. Inspect the suggestion's argv/evidence: numbering may shift as sources change. A selected entry is promoted to commands.start, not an extra synthetic action; multiple selections for one service, including all-scripts/interactive ambiguity, reject before writing. Existing configuration is preserved by init's no-overwrite guard. The unique-entry default remains automatic. Seven groups/17 cases and 15 actual offline Go checks passed. [Baseline](go-selection/baseline-full.log), [after](go-selection/after-race.log), [runtime](go-selection/runtime-summary.json).

Go configuration example:

```yaml
services:
  - id: app-go
    dir: app
    language: go
    manager: go
    go_workspace: off
    infrastructure: {kind: local}
    commands:
      start: {args: [go, run, ./cmd/api], group: aplicacao}
```

Empty/omitted or inherit preserves Go's discovery and inherited GOWORK. Explicit off adds `env GOWORK=off` only to this Go service's child command, including shell and supervisor starts; parent environment remains unchanged. Off is validated only for Go services. Preview and dry-run expose the direct-command override; doctor requires env when needed. Local special operations retain planned-operation display rather than pretending resolved child argv. Three diagnostic groups/three cases and three real read-only CLI checks passed. [Before](go-diagnostics/before.txt), [after](go-diagnostics/after.txt), [CLI](go-diagnostics/cli-after.json).

**MUD-038 remains partial:** go.work membership/parser validation is not implemented, and warnings remain. Real tests show inherited nonmember/missing-member failures and explicitly isolated success, without inventing inheritance policy. [Concrete remaining parser proposal](go-work-parser-plan/README.md) records the dependency constraint and an official pinned candidate; metadata/checksums only were read.

[Five installed runtimes](../isolated-runtimes/README.md) were each rechecked with this final binary: scan and actual quality actions passed. [Recheck summary](runtime-recheck/summary.json). Their original complete fixtures/media remain on the preceding143-test binary and are labelled accordingly; this recheck does not invent new unit-test counts.

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run 'TestPackageManagerSemver|TestGoSelection|TestGoWorkspace' -v
```
