# Yarn Classic and Modern isolated runtime evidence

Versions: Yarn Classic 1.22.22 and Yarn Modern 4.18.1. See install-summary.json, classic-installed-file-integrity.json and node-runtime.json for official sources, verified integrity and actual Node runtimes.

The Python reproduce-install.py and reproduce-fixtures.py document installation and initial fixture creation; use only in a NEW explicitly authorized private root, replacing their base path. They are not safe repeat-init commands for an existing fixture. Installation requires official registry access. No global installation or Corepack enable is used.

For existing preserved fixtures, reproduce-e2e.py repeats actual CLI/wrapper/supervisor and preservation checks; reproduce-pty.py records genuine PTY menu actions. Both use scoped environment JSON and /tmp/mudarro-timeout-checkpoint/mudarro. E2E deliberately changes then restores only its owned menu.sh to test overwrite protection. It shuts down owned fixture services. No dependency install in the personal project occurs.

Final e2e-summary.json reports actual pass results and expected failure exit codes; initial scan/init/generate/test/custom/failure logs add discovery evidence. Final menu casts contain real Yarn execution. GIFs derive from casts, not synthetic screenshots. PTY injected keyboard/wheel events are not a desktop emulator test.

Fixtures have zero dependencies and real generated lockfiles. Modern uses node-modules, not PnP; no workspace, external-dependency, Windows or macOS runtime coverage. Classic initial install uses offline mode; generated install is not network-disabled, although zero dependencies require no package fetch. Modern fixture network is disabled. Modern --frozen-lockfile is accepted with YN0050 deprecation; prefer --immutable.

Initial Classic global-folder fallback emitted a warning mentioning /usr/local and ~/.yarn. Final reruns use an explicit private prefix/global/cache and emit no warning. The specific ~/.yarn path is absent now; no baseline proves all historical writes absent. No personal home data was scanned or cleaned.

No Git, Vault, global settings or repository source changes were made by this worker. Source casts are frozen for parent media rendering.
