[English](../../README.md) · Português (Brasil)

# Mudarro

Detecta a stack de um projeto e gera um menu Bash com ASCII art, submenus e scripts de operação. **Sem IA**: a detecção e a geração funcionam sem rede. Downloads só acontecem em ações de instalação ou nas ferramentas acionadas pelo usuário.

## Instalação

Linux, macOS e WSL; amd64 e arm64. Não exige Go na máquina do usuário.

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh | bash
export PATH="$HOME/.local/bin:$PATH"
```

O instalador baixa a release e verifica SHA-256 antes de instalar. Para fixar uma versão:

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh |
  MUDARRO_VERSION=v0.1.0 bash
```

Atualize repetindo a instalação. Desinstale com `rm "$HOME/.local/bin/mudarro"`; os arquivos dos projetos são preservados. A disponibilidade e os testes da versão estão em [Validação](validation.md).

## Uso

Na pasta da aplicação:

```bash
mudarro scan
mudarro init --interactive       # seleciona scripts existentes para os menus
# Revise mudarro.yaml e resolva as escolhas pendentes.
mudarro generate --dry-run
mudarro generate
./menu.sh
```

Também é possível executar ações diretamente:

```bash
mudarro doctor
mudarro run app:up
mudarro run app:logs
mudarro run app:down
```

Os IDs reais aparecem em `scan` e em `mudarro.yaml`. Os menus gerados dependem do binário `mudarro` no `PATH`.

## Configuração mínima

```yaml
version: 1
name: Minha aplicação
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

| Área | Suporte inicial |
|---|---|
| Linguagens | JavaScript/TypeScript, Python e Go; projetos com múltiplos serviços |
| Infraestrutura | Local, Docker/Compose, Podman/Compose, Kubernetes existente e comandos personalizados |
| Banco | PostgreSQL e SQLite; Prisma, Django, Alembic e Goose |
| Menus | Aplicação, infraestrutura, banco, qualidade, dependências e scripts personalizados |

Infraestrutura ambígua exige configuração explícita. A geração preserva arquivos existentes e detecta mudanças manuais nos arquivos que criou. Migrations e seeds são ações separadas; o Mudarro não inventa o modelo de dados.

- [Configuração e exemplos](configuration.md)
- [Arquitetura, SOLID e pacotes internos](architecture.md)
- [Matriz de suporte atual](support-matrix.md) e [planejamento das próximas linguagens](language-roadmap.md)
- [Operação e limitações](operations.md)
- [Matriz de validação](validation.md) e [sample com menu/GIF reais](samples/menu-auto/README.md)
- [Visual e cores opcionais](terminal-visual.md) · [avaliação da escolha Go](language-decision.md)

## Desenvolvimento

```bash
go test -race ./...
go vet ./...
go build -o bin/mudarro ./cmd/mudarro
bash scripts/smoke.sh "$PWD/bin/mudarro"
```

- [Plano aprovado de implementação da UI configurável](ui-configuration-plan.md)

- [Configuração UI e previews somente leitura](ui-configuration.md)

Gate atual UI configurável: **74 testes,76,38% cobertura de statements**, race/vet/build Linux/cross-build Darwin passaram. [Escopo e limites](coverage.md).
