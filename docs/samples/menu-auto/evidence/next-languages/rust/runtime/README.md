# Rust/Cargo1.99.0 genuine private runtime

PASS:46 actual commands, with Rust/Cargo1.99.0 installed only in the authorized /tmp prefix. Official dated HTTPS archive SHA256891c6366d7100feda0bca4c03ce63f3c9ac827cbebbc283e7433061d42c6a376 matched both approved pin and official checksum before extraction. All67711 archive members were inspected; no external links/special files. Standalone install.sh selected rustc,cargo,rust-std with --disable-ldconfig; no rustup,sudo,profile/global configuration.

Cargo used private CARGO_HOME/CARGO_TARGET_DIR/RUSTUP_HOME/TMPDIR and CARGO_NET_OFFLINE=true. Direct Cargo commands also passed --offline and --locked after real offline generate-lockfile. Workspace contains only local library/app crates; registry files remained absent. Genuine Cargo.lock is archived. No application dependency downloads.

Workspace build/test and repeats passed. Library unit test verifies21; explicitly selected app bin printed42. Explicit failure bin returned7; Mudarro/generated wrapper returned1 with child exit status7. Missing bin and intentional compilation error returned101. Negative build.rs sample was not executed. Runtime validates this sample, not Cargo graph/plugin/build-script support generally.

Mudarro scan returned six opt-in build/test suggestions with no auto commands/start. Explicit selection, init overwrite rejection, generation twice identical, previews, six selected run commands and six generated wrappers passed. User-defined manual start selects owned-server; controlled-failure selects owned-fail. Own service up/status/logs/restart/down/stopped passed, logs contained OWNED_SERVER_42. Final config preserved. No persistent fixture process remained before parent recording; no personal processes were signaled.

Authoritative runtime config: sample/mudarro.json. Legacy manual-start.example.yaml illustrates a finite app bin and is not the long-running runtime configuration. Artifacts exclude prefix/toolchain, downloaded tar, src toolchain extraction, caches,target outputs,binaries,.mudarro state and supervisor identity tokens. Full local originals stay under /tmp/mudarro-approved-runtimes.6AsmDd/rust.

Parent owns PTY screenshots/GIF, full gates, Vault and publication. This worker did not change production/test files during runtime and did not use browser/GUI.
