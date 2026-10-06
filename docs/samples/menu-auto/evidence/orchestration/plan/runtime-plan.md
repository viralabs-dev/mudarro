# MUD-037: isolated Turbo/Nx runtime plan

Status: reviewable plan; installation and runtime execution require the user's next explicit approval. Only official registry metadata and documentation were read. No tarballs, package code, installers or binaries were downloaded. No repository, Vault, Git or global configuration changes.

## Verified inputs

Node `/home/danielsouza/.nvm/versions/node/v24.15.0/bin/node`: v24.15.0. Its sibling npm: 11.12.1. Use these exact existing tools; fixture packageManager must use `npm@11.12.1`, not the other PATH npm version. Registry snapshots and exact dist integrity/tarball URLs are in `*-metadata.json` and `summary.json` beside this plan.

| Tool | Exact package | Linux x64 optional native package |
|---|---|---|
| Turbo | turbo@2.11.7 | @turbo/linux-64@2.11.7 |
| Nx | nx@23.2.1 | @nx/nx-linux-x64-gnu@23.2.1 |

The selected official package metadata declares no engines constraint for these two main packages; this is not proof of successful runtime compatibility. Nx declares a postinstall hook; it will deliberately not execute. Explicit native package installation avoids relying on that hook to obtain a platform package, but successful native binding loading remains a runtime gate. GNU native package assumes this glibc Linux machine, not musl or another OS. Nx transitive dependencies are numerous and not pinned by the main package alone: archive the resolved package-lock, all exact versions and integrity entries before executing anything.

## Proposed authorized download/install stage (NOT EXECUTED)

Run a new private directory, not a reused personal project. Every npm invocation has explicit empty user/global configuration and a private cache. No npx, global install, npm lifecycle hooks, daemon, generator or shell-RC change.

```bash
runtime_root=$(mktemp -d /tmp/mudarro-orchestrators.XXXXXX)
chmod 700 "$runtime_root"
mkdir -p "$runtime_root"/{config,npm-cache,turbo,nx,fixtures,logs}
: > "$runtime_root/config/npm-user.rc"
: > "$runtime_root/config/npm-global.rc"
node_bin=/home/danielsouza/.nvm/versions/node/v24.15.0/bin
export PATH="$node_bin:/usr/bin:/bin"
# Repeat common flags verbatim for both commands.
npm --userconfig "$runtime_root/config/npm-user.rc" --globalconfig "$runtime_root/config/npm-global.rc" --cache "$runtime_root/npm-cache" --registry https://registry.npmjs.org install --prefix "$runtime_root/turbo" --package-lock-only --ignore-scripts --no-audit --no-fund --save-exact turbo@2.11.7 @turbo/linux-64@2.11.7
npm --userconfig "$runtime_root/config/npm-user.rc" --globalconfig "$runtime_root/config/npm-global.rc" --cache "$runtime_root/npm-cache" --registry https://registry.npmjs.org install --prefix "$runtime_root/nx" --package-lock-only --ignore-scripts --no-audit --no-fund --save-exact nx@23.2.1 @nx/nx-linux-x64-gnu@23.2.1
# Inspect lockfiles: root versions exact; every resolved tarball official npm registry;
# every downloaded dependency has integrity; preserve hashes of both locks.
# Stop rather than accept an unexpected source or missing integrity.
npm --userconfig "$runtime_root/config/npm-user.rc" --globalconfig "$runtime_root/config/npm-global.rc" --cache "$runtime_root/npm-cache" --registry https://registry.npmjs.org ci --prefix "$runtime_root/turbo" --ignore-scripts --no-audit --no-fund
npm --userconfig "$runtime_root/config/npm-user.rc" --globalconfig "$runtime_root/config/npm-global.rc" --cache "$runtime_root/npm-cache" --registry https://registry.npmjs.org ci --prefix "$runtime_root/nx" --ignore-scripts --no-audit --no-fund
```

The package-lock-only stage fetches metadata; ci downloads executable package/native code and changes only the private prefix/cache. npm verifies registry integrity against the lock, which provides corruption detection rather than independent authenticity. Compare the four direct lock entries with the captured official integrity; preserve the complete transitive lock and npm logs. Do not enable omitted lifecycle hooks as a fallback without separately reporting the reason and obtaining authorization. npm may retain platform optional packages in the lock while only compatible ones are installed; report actual installed versions independently of lock entries.

## Scoped runtime environment

Use a process-local allowlisted environment for the fixture, without logging inherited secrets. Do not repurpose HOME. If runtime requires HOME-owned persistence despite these settings, stop and report it; do not silently change personal settings.

Turbo: `TURBO_TELEMETRY_DISABLED=1`; invoke absolute private `.bin/turbo`, with `--cache=local:rw --cache-dir "$runtime_root/cache/turbo"`, no remote cache flags/tokens, no `turbo telemetry disable` (that writes persistent settings). Verify `--version` equals 2.11.7. First run with `--force` proves actual fixture execution; repeat without force demonstrates cache semantics separately.

Nx: `NX_DAEMON=false NX_NO_CLOUD=true NX_SKIP_REMOTE_CACHE=true NX_INTERACTIVE=false NX_TUI=false NX_LOAD_DOT_ENV_FILES=false`; set `NX_CACHE_DIRECTORY`, `NX_WORKSPACE_DATA_DIRECTORY` and `NX_NATIVE_FILE_CACHE_DIRECTORY` to distinct directories under runtime_root, native cache mode0700. Use `node "$runtime_root/nx/node_modules/nx/bin/nx.js"`, never registry-resolving npx. Verify version23.2.1; run target with `--skip-nx-cache --outputStyle=static --no-cloud`. Do not call nx reset, init, migrate, generators, or connect-to-cloud. No Nx plugin packages or cloud runner configured. Explicit plugins[] plus standalone project.json avoids intentionally adding inferred plugins; observe the graph/runtime rather than claim plugin execution is impossible.

## Genuine fixtures and gates

Create zero application dependency fixtures inside the private root. Root package.json/workspaces and real npm package-lock generated with the same private npm flags; no framework guessed.

* Turbo: root workspaces packages/*, two members with explicit scripts check/build that execute small local Node scripts and write a marker/artifact. turbo.json tasks includes check/build with explicit outputs and dependency edge if deliberately testing ordering. Run both selected tasks directly; capture output, exit status, cache repeat, changed-input cache invalidation, deliberate nonzero target and missing target. Do not use --parallel, which bypasses its graph.
* Nx: nx.json without cloud configuration or plugins; standalone project.json named sample with target executor `nx:run-commands`, options command `node check.js`, cwd the explicit fixture directory. This built-in executor is documented by Nx and needs no framework plugin. Node script prints a marker and writes an artifact; another target exits nonzero. Invoke sample:check, verify output/artifact, repeat, failure and unknown target. Never treat shell commands read from an untrusted manifest as safe automatically.
* Mudarro integration after parent-owned implementation is ready: scan must perform no tool execution; init explicitly selects supported targets; generate twice compares bytes; menu/wrapper exact root cwd and argv; selected check target executes genuine pinned tool; invalid JSON/missing target cannot silently become a different executable action; cleanup leaves no process. Capture genuine PTY output/frame where viable; label an asciinema-derived frame accurately. No native macOS/WSL claims.

Archive fixture manifests/locks/source, tools/version output, lock/direct integrity comparisons, commands/exit statuses, runtime logs, Mudarro scan/config output and cleanup proof. Samples exclude node_modules, cache, state, secrets and installed binaries. Tool absence or ignored-hook incompatibility is a blocker, not a successful mocked runtime. No installation has happened under this plan.

## Scanner contract review (proposal, no new code available at this read)

Current inspected scan.go/javascript_workspace.go had no Turbo/Nx parser references yet. Parent owns implementation. A bounded typed JSON subset may support explicit turbo.json tasks and standalone Nx project.json targets; it must not claim the complete dependency graph, plugins, inferred targets, JSONC/extends resolution or Nx package.json target merging unless implemented and tested. Parse errors, duplicate/ambiguous targets and unsupported declarations must be visible, not silently reinterpreted as npm scripts. Reject or explicitly skip JSONC instead of stripping comments with regex (strings can contain comment-like characters). Contained regular-file reads and existing exclude-before-read semantics should apply. Root task ownership must survive normalization: same display names cannot merge distinct project:target owners; installation belongs at workspace root while child scripts retain their own cwd. Never execute an orchestrator/plugin during scan, and don't convert arbitrary executor options into shell text. Unsupported configuration is a disclosed validation limitation, not proof of graph correctness.

## Primary sources

* Registry snapshots: https://registry.npmjs.org/turbo/2.11.7 and https://registry.npmjs.org/nx/23.2.1; native snapshots corresponding scoped registry endpoints are captured beside this plan.
* Turbo telemetry: https://turborepo.dev/docs/telemetry
* Turbo run/cache and graph flags: https://turborepo.dev/docs/reference/run
* Nx built-in executor: https://nx.dev/docs/reference/nx/executors
* Nx command shorthand: https://nx.dev/docs/kb/executors-and-configurations
* Nx process/cache/data/native/.env settings: https://nx.dev/docs/reference/environment-variables
* Nx cloud settings: https://nx.dev/docs/kb/config
* Nx daemon: https://nx.dev/docs/reference/nx-daemon

Metadata and documentation checked 2026-10-06. Runtime success and final fixture commands remain unexecuted until approval; documentation reflects current primary docs and must be confirmed against these pinned binaries after download.
