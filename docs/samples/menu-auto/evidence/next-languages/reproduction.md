# Reproduction

Run from the repository with the existing pinned Go toolchain. No installation is necessary for static tests:

```sh
go test -count=1 -race -coverpkg=./internal/mudarro/... ./...
go vet ./...
go build -o /tmp/mudarro-local ./cmd/mudarro
bash scripts/smoke.sh /tmp/mudarro-local
bash scripts/validate-menu.sh /tmp/mudarro-local
python3 scripts/validate-ui-shell-resize.py /tmp/mudarro-local /tmp/mudarro-resize-owned
```

Use a fresh temporary copy of an individual profile sample; retain any original configuration. `scan --json --root SAMPLE` lists evidence and suggestions. `init --root SAMPLE --select SERVICE:ACTION` requires the exact suggestion names returned by scan, and refuses an existing configuration. `generate --root SAMPLE` twice must retain generated wrapper bytes. `preview SERVICE:ACTION --root SAMPLE` and `run SERVICE:ACTION --dry-run --root SAMPLE` display planned argv without calling unavailable toolchains. Do not call build/test of absent tools until separate installation authorization.

For genuine available-runtime checks use frameworks samples with their explicit mudarro.yaml: Node owned_check.cjs / Python owned_check.py use stdlib. Profile READMEs describe their actual commands, expected controlled failures and fixture limits. To replay the terminal use static-contracts.cast; media-validation.json lists every recorded command, final binary hash and injected q. No browser was used. Canonical coverage merges identical source blocks from mirrored test packages by their maximum counter rather than double-counting statements.
