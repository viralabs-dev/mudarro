# MUD-038 official go.work parser focused proof

Before: `before.txt`, `before-complete.txt`, `before-contract-adjusted.txt`. These preserve successive baseline refinements before production integration. Canonical duplicate use was moved from the valid case into a diagnosed invalid case after genuine Go runtime evidence.

After: `after.txt`: 7 top-level groups (6 substantive + 1 isolated helper), 12 subcases; all PASS with race. Fixtures include quotes/comments/use block/toolchain, canonical normalization, replacements not treated as members, invalid syntax/missing/not-module/invalid module/missing declaration/outside paths/4 MiB limits/symlink/FIFO. Valid nonmembers remain independently discovered. Scanner never invokes the controlled Go executable; snapshots preserve source files.

Reproduce from repository root:

```sh
GOCACHE=/tmp/mudarro-go-cache GOMODCACHE=/tmp/mudarro-mud038-parser/module-cache GOPROXY=off GOTOOLCHAIN=local /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run '^TestGoWorkParser' -v
```

Workspace and member-module FIFOs are scanned in owned test helpers with 3-second watchdogs; every child is awaited. Final helpers ended normally; no watchdog or process termination was necessary after the fix. No user process was signalled. These parser tests run on Linux/Darwin; native Darwin runtime was not exercised here.

Membership/module declarations are parsed, not the dependency graph. Replacement targets, toolchain availability, imports, module downloads and cross-target compilation remain separate runtime concerns. ReadManifest has a 4 MiB per-file limit and static symlink/special-file checks; concurrent path replacement remains outside its atomic guarantees. No dependency or tool was installed by this test agent; the root integrator supplied the authorized project library in the private module cache.

Read-only follow-up evidence: `duplicate-module-path-go.json` proves real Go rejects distinct member directories sharing one module path. Root notified; proposed explicit conflict warning without hiding member descriptors. At this checkpoint the focused tests above do not assert this additional case.

## Final follow-up

`before-followup.txt` reproduced duplicate module-path silence and excluded-member IO validation. `excluded-manifest-baseline.txt` captures the intermediate checkpoint: directory exclusion passed, manifest-specific exclusion still failed. Final `after-followup-fresh.txt` uses `-count=1`: 9 top-level groups (8 substantive + 1 helper), 14 subcases all passed with race. Canonical use duplicates and duplicate module declarations are diagnosed; member metadata is preserved. Explicit directory/manifest exclusions plus default .mudarro exclusion are reported in Excluded without reading private FIFO manifests. Four exclusive FIFO helper invocations ended normally and were awaited; no watchdog or process signalling was required. The prior duplicate-module-path review gap is resolved.
