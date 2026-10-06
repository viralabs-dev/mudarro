# packageManager contract triage — MUD-035 (analysis only)

Fourteen genuine offline CLI scan fixtures, [raw results](matrix.json), on local binary from MUD-053; no manager executed or dependency installed. This is characterization, not approval of malformed values or new Bun support. No production change in035. [Proposed matrix](readonly-audit.md), [workspace dependency analysis](workspaces-readonly.md).

Observed: absent/empty→npm; npm/pnpm/yarn versions select names; unversionedpnpm accepted. Leading whitespace, empty version, unknown names and path-like names accepted as managers, producing invalid/unvalidated argv; empty name yields pending start. Conflictinglocks without override are pending; explicitpnpm wins bothlocks. Bun literal is simply accepted by generic splitting, without proving adapter/lockfile/runtime support.

Remaining contract must reconcile known-manager/version validation and existing unversioned/override compatibility with dependencyMUD033. No Bun ambiguity policy is decided by this triage. Explicit commands in consumer configuration are separate from inferred packageManager validation. Runtime stage remains separate, using existing tools only.

To reproduce any case, create an exclusive directory containing package.json with the matrix's packageManager and scripts.start="printf fixture-only", optional named lockfiles, then run `mudarro scan --root <fixture> --json`. Retain raw output; do not run inferred actions of invalid cases.


Rechecked after Bun033 on binary SHA2566214994c4e1dfcc06acfbd7beae14e9b88e592b7e25049fbb06a55e3a8c12033: [14 real CLI cases](matrix-after-bun.json), [comparison](recheck-summary.json). All exited0; contract output unchanged after normalizing only temporary fixture name. Bun lock/runtime support is now separately validated in [Bun033](../bun/README.md), while malformed packageManager values remain an actual open issue. No manager execution/installation or035production change.
