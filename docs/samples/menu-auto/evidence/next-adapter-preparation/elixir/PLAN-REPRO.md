# MUD-044 Mix preparation / Preparação Mix

Existing runtime: `/home/linuxbrew/.linuxbrew/opt/elixir/bin/mix`, Mix 1.20.2, Erlang/OTP 29. Ten real checks passed: compile/test/repetition for root and umbrella; deliberate ExUnit failure and invalid Mix syntax rejected. No Hex setup, dependency download, third-party dependencies or persistent application was used. `summary.json` and per-fixture logs contain actual status/output. Local TCP filesystem locks required authorized sandbox escalation. Initial umbrella application name did not match its directory; corrected only this owned sample to `:member` before passing.

Runtime existente: dez verificações reais aprovadas para root/umbrella, repetição e falhas esperadas. Nenhum download de dependências ou configuração Hex. A escalada autorizada permitiu apenas os sockets locais de locks exigidos pelo Mix. Não há validação de adapter Mudarro nesta etapa.

## Proposed offline contract / Contrato offline proposto

Detect contained regular `mix.exs` as evidence for language `elixir` and manager `mix`; bounded reads must reject FIFO, symlink escape and oversized sources without executing code. Do not evaluate Mix.Project, shell, aliases, dependencies or user-defined functions during scan. Scan must disclose that arbitrary project/umbrella semantics are not validated.

Conventional suggestions may explicitly offer `mix compile`, `mix test` and `mix deps.get` in the owner's cwd; selecting the latter authorizes that later action, not execution during discovery. Keep startup unresolved: no invented `mix run --no-halt`, Phoenix server or release action. Native `mix.exs` can declare anything; aliases cannot be safely inferred from this evidence-only contract.

Umbrella: keep root and discovered member evidence visible, but do not infer ownership solely from directory name `apps/`. `apps_path` is executable code. Either a declarative config explicitly supplies umbrella ownership, or remain independent service contexts with a warning. Avoid duplicate automatic compile/test execution or inheritance of infrastructure/database commands. Preserve manually selected paths and config bytes; repeated init rejects existing configuration.

Negativos futuros: scan with executable trap Mix on PATH never invokes it; malformed Mix code remains disclosed evidence without pretending valid project; special/oversized files produce diagnostics; missing start remains pending; explicit custom cwd/config survives; umbrella path/ownership declared manually is verified for containment. Root/member compile/test should be proven by real runtime, separately from scan assertions.

## Reproduce / Reproduzir

`prepare.py` creates ONLY owned sources in `/tmp/mudarro-mud044-preparation`, then invokes absolute existing Mix. It rewrites its own fixture files; do not aim it at a personal project. Environments use private MIX_HOME/HEX_HOME/MIX_BUILD_PATH/MIX_DEPS_PATH/TMPDIR and preserve HOME. No command runs `mix deps.get`, installs Hex, or evaluates a Mix file for scanning.

`sample-list.json` allows source-only copies of root, umbrella, failure and invalid samples. Exclude `_build`, deps, state, logs, caches and binaries. Future adapter execution and Mudarro integration remain pending parent implementation; these preparation results are not adapter support claims.
