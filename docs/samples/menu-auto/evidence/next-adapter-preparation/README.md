# Next adapter preparation — MUD-044/MUD-045

Preparation only; no new adapter implemented. Existing runtimes exercised internal zero-dependency samples: Mix root/umbrella10checks; Make C/C++40commands. Sources/configs/logs/EN/PT plans saved here; compiled binaries, _build, deps, cache/state are excluded. Existing Mudarro custom-command/manual integration proves that route only.

[Mix contract and reproduction](elixir/PLAN-REPRO.md), [Mix results](elixir/summary.json); [Make contract](native/plan.md), [Make results](native/summary.json), [Make PT](native/README.pt-BR.md). [Remaining nine cards](queue/README.md).

Reviewable recommendation before production: Mix file-only offline detection with compile/test opt-in, umbrella ownership only when declared explicitly, no startup/manifest evaluation/deps downloads. Make C/C++ stays existing custom service with explicit selected targets; mixed source trees remain custom, sources without Make are evidence/pending only, no guessed compiler flags/start. CMake/Meson remain outside this phase. These choices are proposals, not delivered adapter support. No new tool installation or publication.
