# Reproduce JavaScript workspace checks

Copy a saved workspace sample to a fresh owned temporary directory. Use the exact already-approved manager/interpreter and subprocess-scoped private caches/config listed in runtime reports; do not execute archived helpers without replacing their original exclusive /tmp paths. Samples omit installed node_modules, stores, tempdirs, process state and caches.

Build Mudarro with an existing Go SDK, put its binary directory on subprocess PATH, run scan --json, init only if no config exists, and generate twice. Run app-javascript:install at the workspace root; run packages-a-javascript:test and packages-b-javascript:test to confirm their cwd markers. A member has no automatically generated package-manager install action. Explicit custom/manual commands remain usable. Existing config must not be overwritten.

Core fixtures require zero third-party application dependencies. Node runtime probes write only OWN_CWD_MARKER inside the copied fixture. pnpm11 uses pnpm_config_store_dir; Bun uses scoped TMPDIR/BUN_INSTALL_CACHE_DIR. Yarn private wrappers/caches are required, without changing global configuration. Prefer npm sample for a minimal replay with an already existing npm.

```bash
go build -o /tmp/mudarro-js-demo ./cmd/mudarro
fixture_dir=$(mktemp -d /tmp/mudarro-js-demo.XXXXXX)
cp -R docs/samples/menu-auto/evidence/javascript-workspaces/samples/npm/. "$fixture_dir/"
/tmp/mudarro-js-demo scan --root "$fixture_dir" --json
/tmp/mudarro-js-demo run packages-a-javascript:test --root "$fixture_dir"
/tmp/mudarro-js-demo menu --root "$fixture_dir"
```

The zero-dependency child test uses node directly through npm and does not require installation. Actual root install was tested separately with private npm configuration/caches; replaying it should also scope NPM_CONFIG_USERCONFIG/GLOBALCONFIG/CACHE and disable audit/fund, not modify user-global settings.
