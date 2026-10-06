# MUD-044/MUD-045 — Mix and explicit native Make workflows

**199 race groups PASS, zero failures/test skips;2634/3223 production statements=81.73%.** Includes optional real PTY recording. [Summary](summary.json), [events](test-results.json), [coverage](verified-coverage.out). Vet, Linux build, Darwin arm64 cross-build, smoke, real menu/config and PTY resize PASS. Binary SHA256 `118b9c03628481fd3bf662c1590a036e0ec812d744d59055d23983bdd3a5aef4`. Sources local on published6c868c95; no publication of these adapters.

## Implemented contracts

Mix recognizes contained regular bounded4MiB mix.exs as executable-source evidence only. It neither evaluates nor parses Elixir project code, aliases, dependencies or apps_path. Invalid source syntax/UTF8 is still evidence, not proof of a valid project; the real Mix runtime validates execution. Compile/test are opt-in suggestions with exact argv/cwd. No inferred start/install/deps download. mix_umbrella:true is an explicit configuration declaration requiring language:elixir, manager:mix and readable regular mix.exs; no application membership/graph is inferred or inherited. Each nested manifest stays independent unless user config chooses services/actions. Mix handles declared umbrella compile/test at runtime. deps/_build ignored ONLY adjacent to regular nonexcluded mix.exs, preserving unrelated JS/Go/Python projects.

Native .c/.cpp/.cc/.cxx sources and .h/.hpp headers are filename/regular-file evidence; source contents are not parsed/read. Sources group under nearest readable bounded Makefile, otherwise by source directory. Existing language ownership stays intact; deeper nested Make creates independent custom service. Headers alone do not create a native service. Custom actions/start/compiler flags are never invented. Make targets are existing static opt-in suggestions or user commands; arbitrary dynamic targets require explicit config. No full Make parser, CMake/Meson, dependency graph or automatic compiler support claim. Orphan Make/shell files are diagnosed and skipped instead of falling back to another service.

| Matrix | Result |
| --- | --- |
| Mix no-evaluation traps/invalid evidence/bounds/symlink/excluded/nested/explicit umbrella config | 5 new race groups PASS |
| Native regular/socket/excluded/symlink/Make boundaries/owners/mixed/source-only/collisions/orphans | 7 new race groups PASS |
| Global Mix ignore hiding a JS workspace | Reproduced before FAIL; context-scoped fix + regression PASS |
| Mix root and umbrella scan/select/generate twice/preview/run/wrappers/manual config/failure restore | 26 actual checks PASS; [results](mix/runtime-summary.json) +8 [final rechecks](mix/final-recheck.json) |
| C/C++/mixed/source-only automatic scan, explicit Make selection, generation/menu/wrappers/repeat/fail/missing | 63 commands; [summary](native/runtime-summary.json), [commands](native/runtime-results.json), [final recheck](native/final-recheck.json) |
| Genuine menu+selected test in fresh source copies | Mix compiles/runs ExUnit; mixed C/C++ compiles both and verifies21/42; exit0/termios restored |

Existing Mix1.20.2/OTP29, Make4.3 and GCC/G++13.3 only. No downloads/Hex/install/start/persistent apps. Mix locks use owned local TCP sockets; private Unix-socket negative test also required authorized escalation after ordinary sandbox EPERM. These environment failures are not production passes. No security/global config change. [Reproduction](reproduction.md), [Mix source/config snapshot](mix/samples/umbrella/mudarro.json), [mixed Make snapshot](native/samples/mixed/Makefile). Samples exclude builds/deps/cache/state/binaries.

![Actual Mix umbrella test](media/mix.png)

![Recorded Mix test](media/mix.gif)

![Actual mixed C/C++ build and test](media/native.png)

![Recorded native test](media/native.gif)

Casts accompany media; raw PTY bytes remain local; injected q exits the actual menu, then CLI runs the selected test. AllGIFframes decoded Mix6/native4 at944×647; PNGframes5/3 visually inspected. No browser, desktop screenshot or synthetic frames. Framework/runtime dependency coverage, native macOS/WSL/Podman, physical font/mouse and graphical capture remain unexecuted/separate.

Publication selection: [actual command excerpts](command-evidence.json), [file hashes](SHA256SUMS). Full raw logs, test event stream, helper scripts and intermediate outputs remain local and are excluded. No commit/push performed.
