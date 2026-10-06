English · [Português (Brasil)](pt-BR/language-roadmap.md)

## Current authorized runtime checkpoint

The approved private toolchains in `/tmp/mudarro-approved-runtimes.6AsmDd` now have genuine local validation: **Rust1.99.0: 46 checks** (offline build/test, selected binary42, controlled failure7, startup/restart/down); **JDK25.0.4.1+1/Maven3.10.0: 18 checks** (real Java standard-library self-test42/failure7; Maven version only); **.NET10.0.401: 31 checks** (console self-test42/failure7, build failure/recovery, zero NuGet packages; HOME preserved in final run). [Runtime evidence and reproduction](samples/menu-auto/evidence/next-languages/runtime/README.md).

Maven lifecycle/plugins and Gradle remain blocked pending separate approval; no Maven compile/test runtime is claimed. .NET console validation does not claim dotnet test/framework execution. PHP/Composer and Ruby/Bundler remain unexecuted. Framework metadata and Node/Python standard-library checks do not prove Flask/FastAPI/Express runtime. The seven production gates passed again with **235 groups/83.06%**, without a production fix or new publication. Rust and .NET completion applies only to their approved subsets; broad language/framework support is not implied.

## Historical static checkpoint — MUD-039/040/041/042/043/005

Static adapters and framework metadata checks are implemented locally. Final gates: **235 top-level groups PASS, zero test failures/skips**, four packages without tests; canonical coverage **3050/3672 = 83.06%**. Race, vet, Linux build, Darwin arm64 cross-build, smoke, menu and resize exited0. [Final evidence](samples/menu-auto/evidence/next-languages/README.md). Cross-build does not prove macOS runtime. No Rust/Cargo, JDK/Maven, PHP/Composer, .NET or Ruby/Bundler runtime was executed. Existing Node/Python standard-library checks do not prove Flask/FastAPI/Express runtime support. Installation plans are research only, without authorization to install. [Preparation evidence](samples/menu-auto/evidence/next-adapter-preparation/README.md) · [Reviewed installation plans](samples/menu-auto/evidence/next-adapter-preparation/queue/install-plans/README.md).

| Capability | Implemented static contract | Remaining limits |
| --- | --- | --- |
| Rust/Cargo | Bounded TOML, explicit contained workspace/member/target diagnostics; build/test opt-in | No dependency/build-script evaluation, glob graph resolution or invented startup; runtime pending |
| Java/Maven; Gradle/Kotlin | Bounded Maven XML/module checks; compile/test suggestions for eligible projects; Gradle DSL evidence only | No plugin/profile/property or Gradle evaluation; no inferred main/start; JVM runtime pending |
| PHP/Composer | Bounded Composer metadata and explicitly selected script suggestions | No plugins/hooks execution during scan or inferred framework startup; runtime pending |
| C#/.NET | Bounded project XML, explicit per-project build/test suggestions; solution evidence | No MSBuild evaluation or resolved solution graph/startup; runtime pending |
| Ruby/Bundler | Bounded Gemfile/gemspec/lock evidence | No Ruby DSL evaluation, lock semantics, inferred Rake/Rails/start or automatic commands; runtime pending |
| Framework metadata | Declared framework/test dependency evidence and explicit entrypoint contract | No importing application modules during scan or framework runtime claim |

MUD-044/MUD-045 were published at `c400f99` with personal `daneiel` attribution, 89 files and `[skip ci]`; two authenticated checks found zero Actions runs for that publication. Its 199-test/81.73% checkpoint and actual Mix/native samples remain historical evidence, not measurements of this new static batch. [Published Mix/native evidence](samples/menu-auto/evidence/mix-native/README.md).

Earlier checkpoints below retain their original results and publication/runtime limits at the time recorded.

## Historical checkpoint — MUD-044/MUD-045

199racePASS/zero failures-testskips/81.73%; static Mix opt-in compile/test and explicit umbrella, C/C++custom Make user targets validated in internal samples.26Mixchecks+8rechecks/63nativecommands+rechecks; no installation/inventedstart/publication of this batch. MUD037 published6c868c95[daneiel/skipci/0Actions]. [Evidence and limits](samples/menu-auto/evidence/mix-native/README.md).


## Earlier implementation inventory — before this static batch

Bun locks and Bun1.4.0 real runtime are validated locally. Pipfile/Pipenv detection, scripts and database argv are implemented and tested, and Pipenv2026.8.0 runtime now passes a zero-dependency fixture. Known packageManager names are validated; supplied versions now require exact SemVer syntax (metadata syntactic only). Go host entrypoints and bounded source IO are corrected; explicit target selection and scoped isolation pass; full workspace membership parser remains open. Yarn1.22.22/4.18.1, uv0.12.23 and Poetry2.5.1 now pass isolated zero-dependency runtime validation. This narrows the remaining roadmap below; no new-language demand or adapters are inferred. [Latest evidence](samples/menu-auto/evidence/isolated-runtimes/README.md). Changes uncommitted on published base b887f64.


## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.

# Earlier prioritization proposal (superseded by the approved static contracts)

No new-language demand was established in the reviewed material. This order is a hypothesis about value/cost for Daniel's actual repositories, not a popularity ranking or comparative benchmark.

1. **Close manager/monorepo gaps in current families.** remaining packageManager provisioning, real dependency/framework integrations, remaining JS dependency graph/Turbo/Nx (declared ownership/manager inheritance delivered) and Go module graph/framework integrations (entry selection and official go.work membership parsing delivered). Low/medium cost; define inheritance and selection before implementation. Require realistic fixtures, explicit ambiguity and generation/menu tests.
2. **Rust/Cargo only with a real target repository.** Existing TOML parsing may help, but workspace/default-members, library vs multiple binaries and targets need explicit rules. Medium complexity; no automatic start merely from Cargo.toml. Require real authorized cargo execution before claiming runtime support.
3. **PHP/Composer only with personal web demand.** JSON/scripts fit existing mechanisms; PHP/Composer runtime and framework startup differ. Medium complexity; start with explicit scripts, evidence-based Laravel/Symfony separately.
4. **Java/Kotlin or .NET only with concrete demand.** Maven/Gradle wrappers/multimodule or csproj/solution selection increase complexity. Separate build/test from startup; do not execute wrappers during detection.
5. **Ruby/Elixir/C/C++ after a confirmed target.** Managers/frameworks and Make/CMake/Meson targets add ambiguity. Custom commands already cover some workflows; measure actual friction first.

Prioritize an authorized real repository, repetitive manual action, adequate static evidence, offline detection test and safe minimal E2E. Record demand/value/complexity/dependencies in Kanban; one integrator writes shared files, independent reviewers use isolated snapshots. Avoid a nominal catalog that recognizes files without reliable commands.

States: inventoried → demand confirmed → contract defined → adapter implemented → generation tested → runtime exercised → limits documented. Recognizing a manifest or compiling Mudarro is insufficient to claim support.

MUD-015 plans capabilities; MUD-016 expands current tests. JS script collisions were fixed with regression evidence. Other gaps require explicit prioritization. No Rust/PHP/Java/.NET adapters were implemented in this round. See [matrix](support-matrix.md).


Current Go workspace parser, explicit entry selection and inherited/off policy are delivered and validated locally; replacement/toolchain/dependency graph resolution and workspace service inheritance remain outside this implementation. [Evidence](samples/menu-auto/evidence/go-work-parser/README.md).


Declared JavaScript workspace manager-only inheritance and root installation are delivered locally; dependency graph/Turbo/Nx and full glob-language equivalence remain separate. [Evidence](samples/menu-auto/evidence/javascript-workspaces/README.md).


## MUD-037 — local checkpoint

Explicit opt-in Turbo/Nx tasks and root ownership are implemented locally. Real pinned Turbo2.11.7/Nx23.2.1 build/test/dependency ordering/local cache/repeat/input-change/failure and Mudarro wrappers passed using internal samples.186racePASS/81.38%, no test skips. Strict JSON only; no scanner graph/plugin/default inference. Native macOS/WSL/Podman and physical acceptance remain pending. New batch not published. [Evidence and limits](samples/menu-auto/evidence/orchestration/README.md).
