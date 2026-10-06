# Genuine Turbo runtime evidence

Tool installation audited by parent: `/tmp/mudarro-orchestrators.c6vPF6/turbo/node_modules/.bin/turbo`, version 2.11.7. Existing Node24.15.0/npm11.12.1 only. `env.json` contains the complete allowlist used; it includes private cache/XDG/tmp/npm configuration paths and no inherited credentials/HOME. Tool installation and npm integrity audit are separate parent evidence.

Final fixture: `fixture/`. Archive package.json, package-lock.json, turbo.json, scripts/, packages/*/package.json and packages/lib/input.txt; omit node_modules, .turbo, dist, installed binaries and caches from reusable sample. `source-hashes.json` records the final fixture and generated configuration/wrappers.

`runtime.py` records scan/init/select/generate twice/preview/run/generated wrappers and negatives. `cache-proof.py` records explicit input cache/order proof. `final-check.py` rechecks scan/run/wrappers with the final CLI binary. These are historical single-run harnesses, not safe commands to rerun over the preserved fixture: reproductions must use a new private base/fixture/cache and adjust the hardcoded base. No cleanup/delete of existing evidence is needed.

Use existing authorized private Turbo and a fresh copy of the source fixture. Generate its zero-dependency workspace lock/install with exact npm and explicit empty user/global configuration, private npm cache, `--offline --ignore-scripts --no-audit --no-fund`. Scope environment from env.json to the new directory. Execute:

```
<private-turbo> run build --cache=local:rw --cache-dir <new-private-cache> --summarize
<private-turbo> run test --cache=local:rw --cache-dir <new-private-cache> --summarize
<private-turbo> run build --cache=local:rw --cache-dir <new-private-cache> --summarize
# Change packages/lib/input.txt, then repeat build.
<private-turbo> run fail --cache=local:rw --cache-dir <new-private-cache> --summarize
<private-turbo> run absent-task --cache=local:rw --cache-dir <new-private-cache> --summarize
```

Acceptance: lib/app values21/42; lib endTime <= app startTime; cold MISS both, repeat local HIT both with remote false; changed input MISS both and changed hashes. Direct fail exits7; missing task exits1. Mudarro child failure reports exit status7 but CLI/generated wrapper exits1. The initial `--force`+`--cache` incompatibility and initial default-input MISS are preserved separately; neither is a production defect. Final sample declares task inputs package.json/input.txt and globalDependencies scripts/** because no Git default inputs included output/log files.

For menu recording use final binary `/tmp/mudarro-mud037-root/mudarro-final`, the allowlisted env, and action `app-javascript:turbo-test`. Parent owns genuine PTY recording. All subprocesses were awaited; no persistent app or personal process was started/signaled.
