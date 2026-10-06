# Doctor special files — MUD-060 (local)

Relative FIFO with execute bits0755 previously passed doctor as executable: [focused baseline](before.txt), [genuine CLI before](cli-before.json). Guard now requires a regular file plus execute permission; [focused after](after.txt), [genuine CLI after](cli-after.json) reject FIFO with exit1 while regular0755 passes and regular0644 fails. No application action executed and config/source/FIFO inode, mode, size and mtime stayed intact. This diagnostic does not evaluate file contents or guarantee runtime success.

Integrated checkpoint: **141 top-level race PASS (140 substantive groups + one helper), 0 FAIL/SKIP; 2,058/2,604 = 79.03%**. [Summary](summary.json), [test events](tests.jsonl), [coverage](coverage.out), [functions](coverage-functions.txt), [real PTY](interactive.cast). Vet/Linuxbuild/Darwin arm64 cross-build and smoke/menu/config/resize passed; native macOS not executed. Binary SHA2563dd67bf6d6d13281266bb253bf4a630ade6fc909e9a5795a1a08a022e243b9a1. Published baseb887f64; new changes uncommitted, no installs/push/Actions/browser/desktop capture. Compose-provider timeout is separate MUD-061 and not fixed by this checkpoint.

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro -run '^TestDoctorSpecial' -v
```
