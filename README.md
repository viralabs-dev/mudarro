English · [Português (Brasil)](docs/pt-BR/README.md)

# Mudarro

Detect a project's stack and generate a Bash menu with ASCII lettering, submenus and operational scripts. **No AI:** scanning and generation work offline. Downloads occur only through installation actions or tools explicitly invoked by the user.


## Recorded screens

Real CLI execution in a PTY with controlled input. PNGs are rendered recording frames, not desktop screenshots. The banner and footer remain fixed while the center shows menus, preview and execution.

**Recorded revision:** these real screens accompany the persistent shell added by MUD-032. Previous recordings are preserved; validation below was performed locally, independently of CI.

### Dark · English

![Frame of real long output with the banner and footer preserved.](docs/samples/menu-auto/evidence/ui-shell/aurora-shell-dark-en.png)

Frame of real long output with the banner and footer preserved.

![Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.](docs/samples/menu-auto/evidence/ui-shell/aurora-shell-dark-en.gif)

Animated GIF of the real session: preview, execution, error, input, cancellation and return to the menu.

### Light · Portuguese

![Real PTY recording frame in the light palette, with fixed banner and footer.](docs/samples/menu-auto/evidence/ui-shell/aurora-shell-light-ptbr.png)

Real PTY recording frame in the light palette, with fixed banner and footer.

![Animated real light-theme session in Portuguese.](docs/samples/menu-auto/evidence/ui-shell/aurora-shell-light-ptbr.gif)

Animated real light-theme session in Portuguese. Renderer background illustrates the palette; Mudarro does not change terminal background settings.

[See both themes, details and reproduction](docs/ui-shell.md#recorded-screens).

## Current CI correction — MUD-030

Verified [push CI run37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: completed **failure**. Ubuntu, database and container jobs succeeded; macOS executed and failed in EOF-loop/race-cleanup handling, including panic with index n=-1. Therefore macOS is no longer globally “not executed”: local Darwin cross-compilation passed, but the actual CI runtime failed. Local Podman remains unavailable; container CI passed, distinct from local execution.

Authenticated Actions permissions read returned enabled=true/allowed_actions=all; this does not identify who changed settings. No workflow/config settings, enable action, dispatch or rerun was performed. This corrective commit uses the official [skip ci] marker for the authorized main-only push, honoring the requested no-Actions execution without changing settings. Earlier raw gates/evidence are preserved as historical and superseded for current CI status. Local Linux74-test/76.38% results remain valid for their recorded checkpoint. The portability fix passed 75 local Linux tests with race, vet and Linux/Darwin arm64 compilation. Corrected macOS runtime remains unvalidated; skipping CI does not mean it passed.


## Installation

Linux, macOS and WSL; amd64 and arm64. Users do not need Go installed.

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh | bash
export PATH="$HOME/.local/bin:$PATH"
```

The installer verifies SHA-256 before installing. Pin a version:

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh |
  MUDARRO_VERSION=v0.1.0 bash
```

Repeat installation to update. Uninstall with `rm "$HOME/.local/bin/mudarro"`; project files remain. See [validation](docs/validation.md) for measured availability and platform limits.

## Usage

From the application directory:

```bash
mudarro scan
mudarro init --interactive       # select existing scripts
# Review mudarro.yaml and resolve pending choices.
mudarro generate --dry-run
mudarro generate
./menu.sh
```

```bash
mudarro doctor
mudarro run app:up
mudarro run app:logs
mudarro run app:down
```

Actual service IDs appear in scan output and `mudarro.yaml`. Generated menus require `mudarro` on PATH.

```yaml
version: 1
name: My application
services:
  - id: app
    dir: .
    language: javascript
    manager: npm
    infrastructure:
      kind: local
    commands:
      start:
        args: [npm, run, dev]
        group: aplicacao
```

Automatic language families: JavaScript/TypeScript, Python and Go, including multiple services. Infrastructure: local, Docker/Podman Compose and Dockerfile, existing Kubernetes and custom commands. Databases: PostgreSQL/SQLite with Prisma, Django, Alembic and Goose. Menu groups retain their configuration identifiers (`aplicacao`, `infraestrutura`, `banco`, `qualidade`, `dependencias`, `scripts`).

Ambiguous infrastructure requires explicit configuration. Generation preserves existing files and detects manual changes to generated files. Migrations and seeds are separate; Mudarro does not invent business models.

- [Configuration](docs/configuration.md)
- [Architecture and SOLID](docs/architecture.md)
- [Support matrix](docs/support-matrix.md) and [language roadmap](docs/language-roadmap.md)
- [Operations](docs/operations.md)
- [Validation](docs/validation.md), [coverage](docs/coverage.md) and [real menu sample](docs/samples/menu-auto/README.md)
- [Terminal visuals](docs/terminal-visual.md) and [Go decision](docs/language-decision.md)
- [Documentation languages and maintenance](docs/documentation.md)

## Development

```bash
go test -race ./...
go vet ./...
go build -o bin/mudarro ./cmd/mudarro
bash scripts/smoke.sh "$PWD/bin/mudarro"
```

- [Approved configurable UI implementation plan](docs/ui-configuration-plan.md)

- [UI configuration and read-only previews](docs/ui-configuration.md)

Earlier configurable-UI checkpoint: **74 tests,76.38% statement coverage**, race/vet/Linux build/Darwin cross-build passed. [Scope and limits](docs/coverage.md).

- [Persistent terminal shell](docs/ui-shell.md)
