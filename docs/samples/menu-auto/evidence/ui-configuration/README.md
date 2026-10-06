# Genuine configurable UI recordings

These files capture the actual Mudarro CLI with the binary SHA-256 in each metadata JSON. `*.pty.txt` and `*.cast` contain real output; input event/timing metadata is separate. The reproduction helper creates nonsecret Aurora fixtures and injects keys/SGR mouse reports through an isolated PTY. `validation.json` confirms expansion, scrolling, wheel, scrollbar drag, collapse, cleanup and absence of the execution marker. This demonstrates terminal protocol handling, not physical mouse acceptance in a graphical emulator.

`aurora-dark-en.gif` and `aurora-light-ptbr.gif` were rendered from their casts with the authorized existing agg. Both are 983×739, 12.90 seconds including final hold. The corrected dark recording has 17 GIF frames, light 16. `frame-dark-en.png` and `frame-light-ptbr.png` are inspected composited recording frames, not graphical screenshots. The renderer uses asciinema/solarized-light backgrounds; the application changes foreground palettes only and does not change terminal settings or font family.

From repository root, with an already-built absolute binary path:

```bash
python3 scripts/record-ui-preview.py /absolute/path/mudarro /tmp/mudarro-ui-reproduction
agg --theme asciinema /tmp/mudarro-ui-reproduction/aurora-dark-en.cast /tmp/mudarro-dark.gif
agg --theme solarized-light /tmp/mudarro-ui-reproduction/aurora-light-ptbr.cast /tmp/mudarro-light.gif
```

Use `gates.txt` for commands/results and limits, `tests-race.txt` for individual tests and `coverage-summary.tsv` for deduplicated statement coverage. The optional API recording is enabled with `MUDARRO_INTERACTIVE_CAST`; it is separate from these full CLI recordings. Existing-dependency DB mode avoids installs; verify available tool paths before reproduction. Do not run unapproved installers or assume native PostgreSQL/macOS/WSL coverage.
