## Current configurable-UI gate

**74 Test functions passed with race; 1,552/2,032 executed statements = 76.38% (Go displays76.4%)**, 480 unexecuted. Terminal package:379/422=89.81%. Vet, Linux build and whole-binary Darwin arm64 cross-build passed; compilation does not validate macOS runtime. Genuine UI E2E passed custom-config wrapper/supervisor propagation, localized scan, read-only preview and secret redaction, plus smoke/menu execution. Final Python/Go language and pnpm real E2Es passed; Dockerfile and eight database combinations with existing dependencies also passed. Previous stage counts/coverage below remain historical with different denominators. Font integration and actual-emulator mouse acceptance remain pending.

Evidence: [tests-race.txt](samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](samples/menu-auto/evidence/ui-configuration/ui-real.txt).

Implementation status: core config/i18n/theme/action-preview implemented; font integration and actual-emulator mouse acceptance remain pending. This plan retains original design requirements as history; implementation guidance: [UI configuration](ui-configuration.md). Final consolidated gates passed; current measurement is recorded below.

English · [Português (Brasil)](pt-BR/ui-configuration-plan.md)

# Configurable UI plan — MUD-023..028

Plan approved; the user authorized implementation, commit and push. Core features are implemented as summarized above; remaining requirements are identified separately. The functional baseline is checkpoint v5: consumer-name menu, 43 tests, 73.26% statement coverage and validated Aurora recordings. Existing YAML/CLI and generated hash guards remain the compatibility baseline.

## Architecture and order

Keep Bash wrappers thin: they invoke the Go CLI; Bash does not implement a second UI. Keep configuration/validation and action planning in orchestration, translations and rendering under terminal, OS input/terminal lifecycle isolated from the deterministic UI state. Add a read-only preview builder over existing Action/Command values. Only explicit execution invokes Runner; focus, wheel, scrollbar and expand/collapse never do.

1. MUD-025: unified configuration reader and equivalent YAML/JSON contracts.
2. MUD-023/024: message catalog and semantic palettes, preserving line-mode behavior.
3. MUD-026: pure focus/preview/viewport state, then isolated Linux/Darwin raw-input adapter and mouse support.
4. MUD-027: terminal-specific font decision, independent of CLI operation.

Do not introduce a UI framework/dependency yet. First estimate and validate raw-mode cleanup, event parsing and viewport complexity against a small dependency-free PTY spike. If maintainability requires a library, present its dependency/licensing/platform tradeoff before adding it. Execution proceeds in staged gates; approval does not imply implementation or validation.

## Configuration contract

One internal Config model with yaml/json tags and common validation; no duplicated parsers with different semantics. Proposed optional `ui` fields extend version 1; existing configs remain readable by the new binary. Older strict binaries reject new fields: document rollback by removing `ui`, never silently rewriting user files. A future incompatible change needs version 2 and explicit conversion.

Proposed defaults: locale=en, theme=auto (terminal foreground/background preserved), density=comfortable, lettering=auto, preview.enabled=true, preview.mouse=auto. Locales en/pt-BR; themes auto/light/dark; density comfortable/compact; lettering auto/ascii/text. Light/dark initially select contrasting semantic foreground palettes; changing an emulator background is separate and never implicit. Font-family is not accepted as an effective generic setting until a supported integration exists.

Selection: explicit new `--config` project-relative guarded path; otherwise exactly one of mudarro.yaml or mudarro.json. Neither: existing missing-config error. Both: ambiguity error requesting explicit selection, no merge. `init` keeps YAML default; explicit format selects JSON. Generated wrappers must carry a non-default config choice in structured, safely quoted arguments. Existing `--json` output mode remains distinct from config format.

Precedence for UI options: explicit CLI override > selected project file > built-in defaults. No new global user config or environment locale guessing. NO_COLOR and TERM/TTY capability constrain color/interactivity independently. Validate unknown fields, duplicate keys, multiple documents/trailing JSON, types/enums and size/path/symlink limits identically. JSON duplicate-key rejection needs deliberate validation: encoding/json alone accepts duplicates. Blank optional fields receive defaults; invalid nonblank values produce actionable errors. Saving preserves user-owned files and the generator's hash protections.

Equivalent examples accepted by the implementation:

```yaml
version: 1
name: Aurora
ui:
  locale: en
  theme: dark
  density: comfortable
  lettering: auto
  preview:
    enabled: true
    mouse: auto
services:
  - id: app
    dir: .
    infrastructure: {kind: local}
    commands:
      test: {args: [npm, test], group: qualidade}
```

```json
{"version":1,"name":"Aurora","ui":{"locale":"en","theme":"dark","density":"comfortable","lettering":"auto","preview":{"enabled":true,"mouse":"auto"}},"services":[{"id":"app","dir":".","infrastructure":{"kind":"local"},"commands":{"test":{"args":["npm","test"],"group":"qualidade"}}}]}
```

To use Portuguese, set locale to pt-BR in the same selected file. IDs, group keys, confirmation tokens and argv remain stable; labels/help/errors owned by Mudarro are translated. External tool output is preserved verbatim. English becomes the UI default only in the implementation checkpoint, not retroactively in current recordings.

## Read-only preview and input

Show working directory, structured argv (each argument separately), explicit shell text and generated wrapper/source where applicable. Preserve actual quoting rather than substituting shell evaluation. For special lifecycle/database actions, show a labeled planned sequence with runtime placeholders; do not invent a fixed final command or inspect processes/containers to populate preview. Read bounded guarded files only when explicitly associated with the action. Never source scripts, expand env values or load .env during preview. Show ${ENV_NAME} references; redact recognized secret flags/URL credentials and configured sensitive fields. Arbitrary secrets embedded in free-form shell cannot be reliably inferred: document that limit, use nonsecret fixtures for recordings.

In the proposed interactive action view, numbers/Tab/Shift-Tab change focus; Down expands its preview and Up collapses. PageUp/PageDown/Home/End and mouse wheel scroll the preview; a drawn scrollbar supports click/drag when mouse reports exist. Enter is the explicitly labeled execution action; destructive confirmation remains mandatory. Changing focus, toggling or dragging cannot execute. Input-mode changes need clear on-screen help and tests; legacy numeric line mode remains compatible.

Raw mode requires both input and output TTYs and supported capabilities, independently of NO_COLOR. Preserve and restore termios/cursor/mouse modes on ordinary exit, EOF, error, interruption and before handing stdin to a command; SIGKILL cannot run cleanup. Reacquire after command completion. Disable mouse tracking outside preview interaction and restore prior state where available. TERM alone is not proof of mouse capability. Bounded escape/event parsing and explicit opt-out avoid hanging on unsupported queries. Pipes/dumb terminals use line-mode preview commands and ordinary text, no ANSI cursor/mouse escapes. Narrow windows switch to compact text/viewport; NO_COLOR retains keyboard access with monochrome rendering.

Mouse events follow terminal protocols, not a browser scrollbar. Xterm documents SGR mouse coordinates/press-release and tracking modes; wheel/click/drag must be tested through real PTYs plus at least one actual emulator before claiming real mouse acceptance. [Xterm control sequences](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html).

## Font ownership and decision

A generic CLI cannot universally choose an installed font: the terminal emulator owns font selection. No kitty/wezterm executable was found in this machine's PATH; fc-list exists for a future read-only inventory. This does not establish which emulator currently owns the user's session. First provide density/ascii/text controls, requiring no font installation/settings change.

Optional future integration could launch an explicitly selected emulator with project-local options, without editing global config. WezTerm exposes config-file/config options and font-family selection; Kitty exposes font configuration, while its remote font-size command concerns size and affects the OS window's subwindows. These are emulator-specific contracts, not universal family switching. Do not enable remote-control permissions or open a new GUI window without confirming that integration/session target. Sources: [WezTerm CLI](https://wezterm.org/cli/general.html), [WezTerm font](https://wezterm.org/config/lua/wezterm/font.html), [Kitty font configuration](https://sw.kovidgoyal.net/kitty/conf/), [Kitty remote control](https://sw.kovidgoyal.net/kitty/remote-control/).

Decision needed before font integration: user's terminal/emulator and whether a separate project-launched window is acceptable. Main config/i18n/theme/preview work can proceed independently under the implementation authorization already received.

## Acceptance and coordination

Config parity/defaults/duplicate/unknown/size/symlink/ambiguity tests; old YAML loading and hash preservation; translated app messages with stable identifiers; palette plain/NO_COLOR/narrow cases; pure preview state proving no Executor calls; actual PTY key/mouse/resize/restoration/child-input tests; secrets/injection fixtures and capped previews. Actual emulator mouse/drag and font application are distinct from synthetic event/PTY parsing and must be reported separately. Cross-compilation is not platform runtime.

Integrator owns shared production and Vault. Read-only reviews can run in parallel; edits only in assigned disjoint files. Keep v5 stable, create a new checkpoint after each functional stage, record passed/failed/not-run gates and update EN/pt-BR plus Portuguese-only Vault. MUD-028 delivers this plan; MUD-023..027 start in the recorded order with explicit ownership and gates; no feature completion is inferred. The user authorized implementation, commit and push; the integrator checks gates and scope before publication. Installations, security changes and font integration have no implicit authorization.
