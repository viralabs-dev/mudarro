English · [Português (Brasil)](pt-BR/language-roadmap.md)

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Next language capabilities

No new-language demand was established in the reviewed material. This order is a hypothesis about value/cost for Daniel's actual repositories, not a popularity ranking or comparative benchmark.

1. **Close manager/monorepo gaps in current families.** Bun and packageManager validation, Pipfile without false pip, JS workspace ownership/inheritance and Go multiple binaries/go.work. Low/medium cost; define inheritance and selection before implementation. Require realistic fixtures, explicit ambiguity and generation/menu tests.
2. **Rust/Cargo only with a real target repository.** Existing TOML parsing may help, but workspace/default-members, library vs multiple binaries and targets need explicit rules. Medium complexity; no automatic start merely from Cargo.toml. Require real authorized cargo execution before claiming runtime support.
3. **PHP/Composer only with personal web demand.** JSON/scripts fit existing mechanisms; PHP/Composer runtime and framework startup differ. Medium complexity; start with explicit scripts, evidence-based Laravel/Symfony separately.
4. **Java/Kotlin or .NET only with concrete demand.** Maven/Gradle wrappers/multimodule or csproj/solution selection increase complexity. Separate build/test from startup; do not execute wrappers during detection.
5. **Ruby/Elixir/C/C++ after a confirmed target.** Managers/frameworks and Make/CMake/Meson targets add ambiguity. Custom commands already cover some workflows; measure actual friction first.

Prioritize an authorized real repository, repetitive manual action, adequate static evidence, offline detection test and safe minimal E2E. Record demand/value/complexity/dependencies in Kanban; one integrator writes shared files, independent reviewers use isolated snapshots. Avoid a nominal catalog that recognizes files without reliable commands.

States: inventoried → demand confirmed → contract defined → adapter implemented → generation tested → runtime exercised → limits documented. Recognizing a manifest or compiling Mudarro is insufficient to claim support.

MUD-015 plans capabilities; MUD-016 expands current tests. JS script collisions were fixed with regression evidence. Other gaps require explicit prioritization. No Rust/PHP/Java/.NET adapters were implemented in this round. See [matrix](support-matrix.md).
