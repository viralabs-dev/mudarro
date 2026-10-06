# Doctor Compose timeout — MUD-061 (local)

A genuine CLI baseline hung while a controlled provider executed sleep60. After correction, the version probe timed out in 10.05 seconds, exited 1 and printed TIMEOUT rather than AUSENTE. Only compose version was called; configuration remained intact and the direct provider process ended. [Baseline](baseline.json), [after CLI](cli-after.json), [focused race test](after.txt).

Contract: ten seconds per Compose-provider version probe. This does not bound the entire doctor or application actions, prove provider health, or guarantee cleanup of arbitrary descendants. The isolated test uses exec sleep and an exclusive process group; existing successful/failed provider cases remain covered by MUD-051.

Integrated checkpoint: **143 top-level race PASS (141 substantive groups + two isolated helpers), 0 FAIL/0 SKIP; 2064/2610 = 79.08%**. [Summary](summary.json), [events](tests.jsonl), [coverage](coverage.out), [genuine PTY cast](interactive.cast). Vet, Linux build, Darwin arm64 cross-build, smoke, menu generation, configuration and PTY resize passed. Native macOS not executed. Linux SHA256 `5428ee7ee71e5ac1d92110b0170e9dc2efff34a13bcb41eb08025efc2e472607`. Published base b887f64; new changes remain local, with no installation, push or Actions invocation.

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run '^TestDoctorComposeVersion' -v
```
