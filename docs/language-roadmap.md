English · [Português (Brasil)](pt-BR/language-roadmap.md)

## Current local implementation status — 2026-10-06

Bun locks and Bun1.4.0 real runtime are validated locally. Pipfile/Pipenv detection, scripts and database argv are implemented and tested, and Pipenv2026.8.0 runtime now passes a zero-dependency fixture. Known packageManager names are validated; supplied versions now require exact SemVer syntax (metadata syntactic only). Go host entrypoints and bounded source IO are corrected; explicit target selection and scoped isolation pass; full workspace membership parser remains open. Yarn1.22.22/4.18.1, uv0.12.23 and Poetry2.5.1 now pass isolated zero-dependency runtime validation. This narrows the remaining roadmap below; no new-language demand or adapters are inferred. [Latest evidence](samples/menu-auto/evidence/isolated-runtimes/README.md). Changes uncommitted on published base b887f64.


## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Next language capabilities

No new-language demand was established in the reviewed material. This order is a hypothesis about value/cost for Daniel's actual repositories, not a popularity ranking or comparative benchmark.

1. **Close manager/monorepo gaps in current families.** remaining packageManager provisioning, real dependency/framework integrations, JS workspace ownership/inheritance and Go module graph/framework integrations (entry selection and official go.work membership parsing delivered). Low/medium cost; define inheritance and selection before implementation. Require realistic fixtures, explicit ambiguity and generation/menu tests.
2. **Rust/Cargo only with a real target repository.** Existing TOML parsing may help, but workspace/default-members, library vs multiple binaries and targets need explicit rules. Medium complexity; no automatic start merely from Cargo.toml. Require real authorized cargo execution before claiming runtime support.
3. **PHP/Composer only with personal web demand.** JSON/scripts fit existing mechanisms; PHP/Composer runtime and framework startup differ. Medium complexity; start with explicit scripts, evidence-based Laravel/Symfony separately.
4. **Java/Kotlin or .NET only with concrete demand.** Maven/Gradle wrappers/multimodule or csproj/solution selection increase complexity. Separate build/test from startup; do not execute wrappers during detection.
5. **Ruby/Elixir/C/C++ after a confirmed target.** Managers/frameworks and Make/CMake/Meson targets add ambiguity. Custom commands already cover some workflows; measure actual friction first.

Prioritize an authorized real repository, repetitive manual action, adequate static evidence, offline detection test and safe minimal E2E. Record demand/value/complexity/dependencies in Kanban; one integrator writes shared files, independent reviewers use isolated snapshots. Avoid a nominal catalog that recognizes files without reliable commands.

States: inventoried → demand confirmed → contract defined → adapter implemented → generation tested → runtime exercised → limits documented. Recognizing a manifest or compiling Mudarro is insufficient to claim support.

MUD-015 plans capabilities; MUD-016 expands current tests. JS script collisions were fixed with regression evidence. Other gaps require explicit prioritization. No Rust/PHP/Java/.NET adapters were implemented in this round. See [matrix](support-matrix.md).


Current Go workspace parser, explicit entry selection and inherited/off policy are delivered and validated locally; replacement/toolchain/dependency graph resolution and workspace service inheritance remain outside this implementation. [Evidence](samples/menu-auto/evidence/go-work-parser/README.md).
