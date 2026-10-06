English · [Português (Brasil)](pt-BR/ui-shell.md)

[Jump to recorded screens](#recorded-screens).

# Persistent terminal shell — MUD-032

## Recorded screens

Real CLI execution in a PTY with controlled input. PNGs are rendered recording frames, not desktop screenshots. The banner and footer remain fixed while the center shows menus, preview and execution.

**Recorded revision:** these real screens accompany the persistent shell added by MUD-032. Previous recordings are preserved; validation below was performed locally, independently of CI.

### Dark · English

![Frame of real long output with the banner and footer preserved.](samples/menu-auto/evidence/ui-shell/aurora-shell-dark-en.png)

Frame of real long output with the banner and footer preserved.

![Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.](samples/menu-auto/evidence/ui-shell/aurora-shell-dark-en.gif)

Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.

### Light · Portuguese

![Frame of real long output with the banner and footer preserved.](samples/menu-auto/evidence/ui-shell/aurora-shell-light-ptbr.png)

Frame of real long output with the banner and footer preserved.

![Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.](samples/menu-auto/evidence/ui-shell/aurora-shell-light-ptbr.gif)

Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.



Implemented and locally validated starting from base `9869dc6`. This revision includes the implementation and visible media; publication uses `[skip ci]`, without claiming a new CI pass. Linux: **86 tests passed, no failures/skips**, race and vet passed. Coverage: **75.33% (1,866/2,477 statements)** across internal packages, deduplicating repeated coverpkg blocks. Smoke, automatic-menu/idempotency and selected-configuration regressions passed. Darwin arm64 binary and both test packages compiled; native macOS runtime remains unvalidated. The earlier CI failure remains historical and distinct.

## Frame and capabilities

Supported TTY menus retain header/project banner/footer across service, group, action, preview and action execution. The initial project context is fixed; center content uses absolute cursor positions and DECSTBM scrolling margins. Resizing selects compact layout. With no available content height, Enter execution is disabled. NO_COLOR retains the frame without SGR colors. Non-TTY, dumb terminals or initial dimensions smaller than 3×4 use fallback rather than an unsupported frame.

Action output is bounded: last 256 lines, maximum 4 KiB per line. Incremental capture strips all ANSI sequences before displaying the latest viewport. Execution output is not an interactive scrollback viewer. This is a generic terminal UI, not a sandbox or support for arbitrary fullscreen children; programs that access /dev/tty directly can bypass piped capture.

## Child lifecycle and input

Children receive piped input while output stays in the central viewport. Ctrl-C cancels through context and a dedicated child process group: TERM then KILL after 300 ms if required; WaitDelay 750 ms bounds inherited-pipe shutdown. Errors remain visible until Enter returns to the menu. External SIGTERM currently cancels the action and returns to the menu; it is not described as immediate process exit. Terminal state/margins are restored at completion and return paths. Destructive confirmation remains mandatory and can be cancelled before child execution.

## Genuine validation and reproduction

Dark-English and light-pt-BR isolated PTY captures each passed 24 chrome checks, covering output, error, stdin, explicit destructive confirmation with queued child input preserved, repeated child execution and restoration. These are real application output with controlled PTY input, not actual-emulator acceptance or graphical desktop screenshots. Canonical media destination: docs/samples/menu-auto/evidence/ui-shell, with profile.cast, pty.txt, metadata.json, GIF and composed PNG. Media is shared by both documentation languages ; both GIFs decoded and composed PNG frames were visually inspected.

```bash
python3 scripts/record-ui-shell.py /absolute/path/to/mudarro /tmp/mudarro-ui-shell-reproduction
```

Run from repository root with an actual compiled binary. Output resources are isolated; recording opens no GUI. Renderer theme backgrounds are separate from foreground UI palettes and do not alter emulator settings. Visible graphical tests remain restricted to the previously confirmed workspace 4/integrated notebook display. No font-family integration or installation is included.

Resize/NO_COLOR were checked in 100×32, 12×6 and 8×2 PTYs; Enter was disabled at zero central height, and SIGTERM during navigation restored termios/margins/alternate screen. The resize check evaluates new geometry after its reset; bytes already queued for the previous geometry may arrive first.

```bash
python3 scripts/validate-ui-shell-resize.py /absolute/path/to/mudarro /tmp/mudarro-shell-resize
```

See [raw evidence and media](samples/menu-auto/evidence/ui-shell/README.md).
