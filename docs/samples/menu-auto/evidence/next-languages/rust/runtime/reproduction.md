# Reproduction

Installation is historical, explicitly authorized and fully recorded in download-verification.json/http.headers/rust.sha256/logs/install.*. Do not install globally or invoke rustup. Existing private tools: /tmp/mudarro-approved-runtimes.6AsmDd/rust/prefix/bin/{rustc,cargo}.

Use a fresh copy of sample/ and fresh private cargo/target/tmp cache paths; adjust env.json process-locally without changing HOME. CARGO_NET_OFFLINE=true is required for Mudarro commands, whose configured argv cargo build/test deliberately preserves adapter contract. No negative build.rs execution.

```
cargo generate-lockfile --offline
cargo build --workspace --offline --locked
cargo test --workspace --offline --locked
cargo build --workspace --offline --locked
cargo test --workspace --offline --locked
cargo run --offline --locked --bin owned-app   # output42
cargo run --offline --locked --bin owned-fail  # exit7
cargo run --offline --locked --bin owned-missing # exit101
```

For full scan/select flow, copy source plus Cargo.lock without mudarro.json/generated state, then init --format json --select the six Cargo suggestions. Generate twice and compare hashes. Runtime sample/mudarro.json already contains explicitly configured owned-server and controlled-failure; scan must not overwrite it.

```
<final-mudarro> scan --root <fixture> --json
<final-mudarro> generate --root <fixture>
<final-mudarro> preview --root <fixture> app-rust:test
<final-mudarro> run --root <fixture> app-rust:test
<fixture>/.mudarro/scripts/app-rust/test.sh
<final-mudarro> run --root <fixture> app-rust:controlled-failure # CLIexit1, child7
<final-mudarro> run --root <fixture> app-rust:up
<final-mudarro> run --root <fixture> app-rust:status
<final-mudarro> run --root <fixture> app-rust:logs
<final-mudarro> run --root <fixture> app-rust:restart
<final-mudarro> run --root <fixture> app-rust:down
```

Always down only the owned service after recording. Runtime.py is a historical single-run driver; it creates its fixture once, so adjust base/root to a fresh directory before repeating. compile-failure/ is a separate zero-dependency fixture with syntax error and expected Cargo101, not a production adapter bug. Parent records genuine menu/terminal using env.json; no synthetic image is provided by this front.
