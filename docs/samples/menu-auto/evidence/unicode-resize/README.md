Current published base is `b887f64`, tree-identical to the recorded `56dadba`. History repair056 subsequently completed; this recorded test checkpoint remains historical. Latest local validation: [Bun033](../bun/README.md).

# Unicode and resize — MUD-053 (local)

Base main56dadba with local robustness050/051/052 changes preserved; no new commit/push, history rewrite, Actions, installation or desktop/browser operation. Git identity has separately been corrected only in repository config; published attribution incident056 still awaits approval.

Before: [Unicode baseline](unicode-before.txt) shows U+3000 ideographic spaces counted as one cell, producing28–37 cells in24 columns. [Tiny-width baseline](unicode-tiny-before.txt) shows1/3-column frames emitted4-cell text;4 columns passed. A separate NO_COLOR test assertion initially rejected the harmless ESC[0m cleanup reset; corrected oracle permits only that reset and rejects color SGR.

After: [focused final](unicode-final.txt), five groups and17 variants passed with race (14 Unicode/color variants, three tiny widths). Unicode alphabet covers CJK/fullwidth, combining accents, emoji components, regional pairs and ideographic spaces. Incremental split UTF-8 combining input preserved; bidi/format controls filtered. Actual isolated PTY preview80×24→12×6→8×2→32×12 exited normally and restored exact termios; tiny1/3/4-column frame text and absolute cursor addresses remained within dimensions.

Fix: U+3000 consumes two cells. fitText returns empty at nonpositive width, clips without ellipsis at1–3, and uses reserved ellipsis from4 onwards. It never expands the width budget to four. This remains a conservative rune-cell policy, **not a grapheme shaping engine**: ZWJ/format controls are removed, emoji component sums can overestimate rendered width, and truncation can split an emoji/grapheme sequence. Font/emulator behavior needs separate physical acceptance in026/027; no perfect width guarantee.

Integrated final: **117 top-level race tests PASS, zero failures/skips;1,965/2,515 internal statements=78.13%**. [Summary](summary.json), [test events](tests.jsonl), [coverage](coverage.out), [functions](coverage-functions.txt), [real PTY cast](interactive.cast). Vet/Linuxbuild/Darwinarm64binary+terminaltestcompile passed; native macOS not executed. [Gates](gates.json). Linux binary SHA2560cc82b21da942853aec6d2cf518544d8fcd4e00bdc20161299fe95945de35b5b.

[Smoke](smoke.txt), [automatic menu/idempotence/preservation](menu.txt), [selected configuration](ui.txt), [CLI resize/NO_COLOR/SIGTERM](resize.txt) passed on that final binary. Robustness112/77.41% evidence remains historical in the sibling directory.

Reproduce with existing Go from repo root:

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race ./test/internal/mudarro/terminal -run '^TestUnicodeResize' -v
GOCACHE=/tmp/mudarro-go-cache MUDARRO_INTERACTIVE_CAST=/tmp/mudarro-unicode-interactive.cast /home/danielsouza/sdk/go1.27.1/bin/go test -race -coverpkg=./internal/mudarro/... -coverprofile=/tmp/mudarro-unicode-coverage.out ./...
python3 scripts/validate-ui-shell-resize.py /tmp/mudarro-unicode/bin/mudarro /tmp/mudarro-unicode-resize-reproduction
```
