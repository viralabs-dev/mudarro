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

Atualize repetindo a instalação. Desinstale com `rm "$HOME/.local/bin/mudarro"`; os arquivos dos projetos são preservados. A disponibilidade e os testes da versão estão em [Validação](docs/validation.md).

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

- [Configuração e exemplos](docs/configuration.md)
- [Arquitetura, SOLID e novos adaptadores](docs/architecture.md)
- [Operação e limitações](docs/operations.md)
- [Matriz de validação](docs/validation.md)

## Desenvolvimento

```bash
go test -race ./...
go vet ./...
go build -o bin/mudarro ./cmd/mudarro
bash scripts/smoke.sh "$PWD/bin/mudarro"
```
