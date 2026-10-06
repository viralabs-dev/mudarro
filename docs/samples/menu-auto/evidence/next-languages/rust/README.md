# Rust/Cargo: static minimum (MUD039)

Ten groups/25 subcases passed fresh race. Seventeen genuine Mudarro CLI commands verified scan, explicit selection, preserving existing config, generate twice, preview and dry-run. Cargo/rustc sentinels were never invoked. This is static adapter/CLI coverage, not Rust compilation or runtime acceptance; no toolchain was installed.

Regular bounded contained Cargo.toml is parsed with the existing TOML library. Named packages, explicit lib/bin paths and literal workspace members/default-members are validated. Header/source AST, implicit Cargo targets, dependency graphs, inherited manifest values, resolver/toolchain compatibility, globs and Cargo metadata are not resolved. Glob paths are diagnosed; ambiguous/outside/excluded/invalid declarations leave pending configuration and no automatic executable suggestions for the unresolved declaration. Valid member services remain independent. build.rs is never executed during scan.

Services use language rust/manager cargo with empty Commands; build/test are Purpose cargo opt-in suggestions. Each command retains its service cwd. No install/run/start/server/framework/DB guess. Startup must be explicitly configured, including chosen bin. Multiple explicit bins do not choose startup. Workspace ownership is informational Pending metadata; the JS-only WorkspaceRoot schema is unchanged. Virtual workspace IDs use a workspace suffix; root+app-member collision is covered by scan/init regression.

Sample workspace in sample/: root literal members library/app, default app, local path dependency only, library value21 and app expected42, alternate failure bin exits7 if run in a separately authorized future runtime. Neither output nor exit was executed in this delivery. manual-start.example.yaml is an explicit configuration example, not generated inference. negative-build-script/ proves manifest discovery without evaluating build.rs.

Samples contain no compiled target/, Cargo.lock, caches or installed tools; Cargo.lock cannot honestly be represented as generated without Cargo execution. Future runtime needs a pinned toolchain/installation plan and authorization. Preserve static evidence distinction when copying this artifact.

Final owner regression: nonempty package.workspace is Pending with no suggestions; no owner file is read and no nearest-ancestor ownership is inferred, including contained, oversized, excluded, symlink and outside paths.
