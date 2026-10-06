English · [Português (Brasil)](pt-BR/documentation.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Documentation languages and maintenance

English is the repository default: README.md and docs/**/*.md outside docs/pt-BR. Brazilian Portuguese mirrors live in docs/pt-BR; docs/pt-BR/README.md is the localized project entry point. Vault notes remain Portuguese and are maintained separately.

Assets, generated samples and evidence have one canonical copy outside docs/pt-BR. Localized sample documentation links back to that copy. Do not translate command names, flags, configuration keys/group identifiers, destructive confirmation tokens, log contents, IDs, checksums, fixture paths, licenses or user-supplied values. English prose may retain verbatim Portuguese comments inside existing executable examples.

For behavior changes, update the English document and matching pt-BR document in the same change. Reconcile tables, status, limits and links rather than copying an obsolete historical result. Both versions link to each other. Keep older evidence clearly labeled by stage; new test layout/builds require new actual results. Translation never turns a contract test into runtime validation or a composed recording frame into a screenshot.

If a translation cannot be updated, mark its top with `Translation pending: <source path and changed section/date>` and record the exact follow-up in the activity; do not silently present it as synchronized. Reviewers check matching figures/commands/constraints, language navigation and relative links. No automatic translation service or publishing step is required.

The initial localization preserves complete original Portuguese documents and supplies equivalent English operational guidance. Portuguese historical narratives may be more detailed; shared runtime limits and evidence remain authoritative. Further updates should keep semantic parity, not require identical sentence counts.

- [Approved UI configuration plan](ui-configuration-plan.md)

Current implementation: [UI configuration](ui-configuration.md). Core configuration/i18n/theme/action-preview implemented; fonts and actual-emulator mouse acceptance remain separate pending work. Final gates:74 tests,76.38% aggregate statements; see [coverage](coverage.md).
