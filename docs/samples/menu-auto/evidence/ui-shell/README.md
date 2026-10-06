# MUD-032 — genuine persistent-shell evidence

## Recorded screens

Real CLI execution in a PTY with controlled input. PNGs are rendered recording frames, not desktop screenshots. The banner and footer remain fixed while the center shows menus, preview and execution.

**Recorded revision:** these real screens accompany the persistent shell added by MUD-032. Previous recordings are preserved; validation below was performed locally, independently of CI.

### Dark · English

![Frame of real long output with the banner and footer preserved.](aurora-shell-dark-en.png)

Frame of real long output with the banner and footer preserved.

![Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.](aurora-shell-dark-en.gif)

Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.

### Light · Portuguese

![Frame of real long output with the banner and footer preserved.](aurora-shell-light-ptbr.png)

Frame of real long output with the banner and footer preserved.

![Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.](aurora-shell-light-ptbr.gif)

Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.



Linux validation checkpoint captured on base `9869dc6` before publication; this revision includes the implementation and media, with `[skip ci]`. See [English guide](../../../../ui-shell.md) and [pt-BR guide](../../../../pt-BR/ui-shell.md).

86 top-level tests passed with race and optional recording enabled; no failures/skips. Vet, smoke, automatic-menu/idempotency and UI configuration regressions passed. Unique covered statements: 1,866/2,477 (75.33%). `coverage.out` repeats blocks across tested packages; the summary merges counts by source block before computing this ratio. Darwin arm64 binary and both test packages crosscompiled only.

Two actual CLI PTY captures: dark/en and light/pt-BR, each 24 header/footer comparisons, exit 0 and termios restored. Includes preview, injected SGR mouse bounds, long output with ANSI controls, stderr, error, child stdin, Ctrl+C with descendant termination, execution after cancellation, and explicit destructive confirmation preserving queued child input. Only harmless owned fixture commands were used. No browser, desktop capture or synthetic imagery. PNGs are composed frames from the actual recording GIFs. Physical emulator mouse/graphical acceptance and native macOS remain separate pending work.

Dark GIF: 983×739, 55 frames, 12.31 seconds. Light GIF: 983×739, 50 frames, 12.31 seconds. All frames decode. PNG at 4.2s shows real long output; additional 2.0s/8.0s frames and metadata are retained. The renderer chose DejaVu Sans Mono; backgrounds are renderer themes, not changes to emulator settings.

```bash
go build -o /tmp/mudarro-shell ./cmd/mudarro
python3 scripts/record-ui-shell.py /tmp/mudarro-shell /tmp/mudarro-shell-capture
python3 scripts/validate-ui-shell-resize.py /tmp/mudarro-shell /tmp/mudarro-shell-resize
# Existing agg tool, if available; this step installs nothing.
agg --theme asciinema --font-size 16 --fps-cap 12 --last-frame-duration 1 /tmp/mudarro-shell-capture/aurora-shell-dark-en.cast /tmp/mudarro-shell-capture/aurora-shell-dark-en.gif
agg --theme solarized-light --font-size 16 --fps-cap 12 --last-frame-duration 1 /tmp/mudarro-shell-capture/aurora-shell-light-ptbr.cast /tmp/mudarro-shell-capture/aurora-shell-light-ptbr.gif
```

Raw casts/output, controlled-input event metadata, final test summary/profile and gate logs accompany the media. The command deliberately exits 7 in the error case; this expected error is not a failed gate. Terminal output is bounded/plain and fullscreen children or direct `/dev/tty` access are not supported as contained rich TUIs.

Library: same PNG/GIF identities updated to version 3, preserving previous versions. This does not publish Git changes.
