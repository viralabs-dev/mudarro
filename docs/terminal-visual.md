English · [Português (Brasil)](pt-BR/terminal-visual.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Terminal visuals

TTY semantic colors: cyan titles, orange choices/prompts in 256-color terminals (basic yellow otherwise), gray hints, red errors. Numbers/text retain meaning without color. Default terminal background/foreground remain unchanged. Original 5x7 bitmap lettering with depth/sidebar falls back to embedded 3x5 ASCII; no figlet/network/mandatory visual dependency.

Presentation lives in internal/mudarro/terminal. JSON/noninteractive logs/generated files stay plain. ioctl checks the actual output stream; pipes and /dev/null are not TTYs. Empty/dumb TERM or any nonempty NO_COLOR (including 0) disables colors, following [NO_COLOR](https://no-color.org/). Width comes from window or explicit COLUMNS. Long text is clipped; wide CJK/emoji can exceed rune-count estimates.

Public inspirations, approximate stars checked 2026-10-06: [Lip Gloss](https://github.com/charmbracelet/lipgloss) 11.9k MIT, [Gum](https://github.com/charmbracelet/gum) 24.5k MIT, [Fastfetch](https://github.com/fastfetch-cli/fastfetch) 24.9k MIT, [ascii-image-converter](https://github.com/TheZoraiz/ascii-image-converter) 3.5k Apache-2.0, [VHS](https://github.com/charmbracelet/vhs) 21.1k MIT. Inspiration concerns hierarchy/spacing/colors; no code, logo, font or brand copied.

Official [agg 1.9.0](https://github.com/asciinema/agg/releases/tag/v1.9.0) was authorized only under /tmp/mudarro-media-tools/agg without sudo. It renders genuine PTY recordings, not AI images or desktop screenshots. GIF decode and composed frames were inspected. Final package recording lives in samples/menu-auto/evidence/packages-visual; earlier versions remain. Tests cover plain/color policy/narrow/EOF. User reference pixels were inspected; original Mudarro lettering, orange sidebar/context/responsive footer were implemented. q returns/exits; rich actions wait Enter before redraw, plain output remains sequential.

Graphical screenshot remains pending; user visual design acceptance has been received. Previous portal capture was cancelled. Visible execution is restricted to the already-confirmed workspace 4 on the integrated notebook screen; do not reopen portal without a request. Headless PTY recordings open no desktop window. Composed PNG frames are explicitly identified as recordings, not screenshots.

## Consumer project identity — MUD-022

The menu title is the consumer project `Config.Name`, not the Mudarro tool name. The controlled sample now uses Aurora (package/config); demo and recording default explicitly to Aurora, overridable with `MUDARRO_SAMPLE_NAME`. The existing CLI already used Config.Name; demonstration inputs were corrected. User visual design acceptance has been received; only the genuine graphical screenshot remains pending.

Text sanitization strips controls and Unicode Cf formatting characters. Conservative terminal-cell estimation counts CJK/fullwidth/emoji as two cells and combining marks as zero; it is not a complete grapheme/terminal-width guarantee. Unsupported bitmap glyphs or an excessively long wordmark fall back to readable text instead of presenting an incorrect identity.

The current name stage has 43 Test functions and 24 rich/narrow/plain/NO_COLOR name-mode cases, covering AtlasAPI, Aurora, Café, 漢字😀, long names and control injection, plus explicit consumer-name CLI integration. Race passed; the final coverage and build gates passed as recorded below. Previous stage coverage remains historical. Genuine latest media belongs in `evidence/project-name-visual/` under the canonical sample; `evidence/packages-visual/` remains historical.

## Current consumer-name validation — MUD-022

43 Test functions passed with race, including 24 consumer-name/mode cases. Vet, Linux build, Darwin arm64 terminal-test cross-compilation and Bash syntax passed. Coverage: **1,063/1,451 statements = 73.26% (Go displays 73.3%)**, 388 unexecuted. Prior MUD-017 1,035/1,427=72.53% is historical; changed code and denominator preclude interpreting the difference as equivalent requirement coverage. macOS runtime was not exercised in this historical local stage; current CI runtime failed as recorded above.

Real final GIF: 983×739, 8 frames, 15.04 s. Composed frame inspected: AURORA lettering and PROJETO / Aurora. This is a real PTY recording frame, not the still-pending graphical screenshot. User design approval received. Binary SHA-256: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidence: [project-name-tests.txt](samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Current implementation: [UI configuration](ui-configuration.md). Core configuration/i18n/theme/action-preview implemented; fonts and actual-emulator mouse acceptance remain separate pending work. Final gates:74 tests,76.38% aggregate statements; see [coverage](coverage.md).
