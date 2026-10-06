English · [Português (Brasil)](pt-BR/ui-configuration.md)

## Current configurable-UI gate

**74 Test functions passed with race; 1,552/2,032 executed statements = 76.38% (Go displays76.4%)**, 480 unexecuted. Terminal package:379/422=89.81%. Vet, Linux build and whole-binary Darwin arm64 cross-build passed; compilation does not validate macOS runtime. Genuine UI E2E passed custom-config wrapper/supervisor propagation, localized scan, read-only preview and secret redaction, plus smoke/menu execution. Final Python/Go language and pnpm real E2Es passed; Dockerfile and eight database combinations with existing dependencies also passed. Previous stage counts/coverage below remain historical with different denominators. Font integration and actual-emulator mouse acceptance remain pending.

Evidence: [tests-race.txt](samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# UI configuration

Implemented: one internal Config model and common validation for YAML/JSON, optional version1 ui settings, translated Mudarro-owned messages, semantic themes and read-only action previews. Font integration and real-emulator mouse acceptance remain unimplemented/unverified respectively. Historical evidence is retained; final consolidated coverage is recorded below.

## Selection and defaults

Use --config with a guarded project-relative path, including custom YAML/JSON names. Without it, exactly one mudarro.yaml or mudarro.json is selected; both produce an explicit ambiguity error, neither the missing-config error. No merging. Init defaults YAML; --format json selects JSON. Existing --json controls output, not configuration format. Generated wrappers preserve selected config path in quoted structured arguments; the local supervisor propagates that selection.

CLI UI overrides take precedence over the selected file and built-in defaults. No global config or automatic environment locale guessing. Defaults: locale=en, theme=auto, density=comfortable, lettering=auto, preview.enabled=true, preview.mouse=auto. Locales en/pt-BR; themes auto/light/dark; density comfortable/compact; lettering auto/ascii/text. Unknown fields, duplicate keys, trailing/multiple documents, invalid types/enums and guarded path/size/symlink rules are validated. Older strict binaries reject ui fields; rollback requires explicit removal, not silent rewriting.

```yaml
version: 1
name: Aurora
ui:
  locale: pt-BR
  theme: dark
  density: compact
  lettering: auto
  preview:
    enabled: true
    mouse: auto
services:
  - id: app
    dir: .
    infrastructure: {kind: custom}
    commands:
      test: {args: [npm, test], group: qualidade}
```

```bash
mudarro init --format json
mudarro menu --config project.json --locale en --theme light
```

Identifiers, configuration/group keys, argv and destructive-confirmation tokens remain stable. External process output is not translated. Themes select contrasting semantic foreground colors, preserving emulator background; NO_COLOR disables color but not keyboard access where interactive capabilities exist. A CLI cannot universally select font family; no font installation/global terminal change is implemented.

## Action preview and input

Raw interaction is restricted to action views, leaving service/group selection in compatible line mode. Numbers/Tab/Shift-Tab focus actions; Enter executes explicitly; Down expands preview, Up collapses. PageUp/PageDown/Home/End and mouse wheel scroll; scrollbar click/drag handles reported mouse events. Merely focusing, expanding or scrolling never invokes Runner. Destructive confirmation stays mandatory. Fallback line mode supports pN for a read-only preview and ordinary numeric execution.

Preview shows cwd, argv/shell/source and labeled planned special-action sequences. It never executes commands, inspects processes/containers, sources scripts, loads .env or expands environment values. Environment references remain placeholders; recognized secret flags/URL credentials are redacted. Arbitrary embedded free-form shell secrets cannot be inferred reliably; recordings use nonsecret fixtures. File previews are guarded and bounded.

Raw mode requires supported input/output TTYs; restoration covers normal exit/EOF/errors/interruption and occurs before handing stdin to a child, then reacquires after completion. SIGKILL cannot run cleanup. Pipes/dumb terminals remain plain without mouse/cursor modes. Conservative Unicode cell estimation is not perfect grapheme handling. Mouse protocol events were exercised through isolated PTYs, not accepted in an actual emulator; graphical screenshot remains pending. See [plan and limits](ui-configuration-plan.md), [configuration](configuration.md), [visuals](terminal-visual.md).
