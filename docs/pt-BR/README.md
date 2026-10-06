[English](../../README.md) · Português (Brasil)

# Mudarro

Detecta a stack de um projeto e gera um menu Bash com ASCII art, submenus e scripts de operação. A detecção e a geração funcionam sem rede. Downloads só acontecem em ações de instalação ou nas ferramentas acionadas pelo usuário.


## Telas gravadas

Exemplo atual: projeto `mudarro`, definido por `Config.Name`, com capturas reais atualizadas em MUD-055. Esta revisão publica o exemplo atualizado a partir de `246048a`; a etapa anterior Aurora foi preservada como histórico.

Execução real do CLI em PTY, com entrada controlada. PNGs são frames renderizados da gravação; não são screenshots do desktop. O banner e o rodapé permanecem fixos enquanto o centro mostra menus, prévia e execução.

**Revisão gravada:** estas telas reais acompanham o shell permanente de MUD-032. As gravações anteriores foram preservadas; a validação abaixo foi executada localmente, independente da CI.

### Claro · português

![Frame da saída real longa, com banner e rodapé preservados.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-light-ptbr.png)

Frame da saída real longa, com banner e rodapé preservados.

![GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-light-ptbr.gif)

GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.

### Escuro · inglês

![Frame real da gravação PTY na paleta escura, com banner e rodapé fixos.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-dark-en.png)

Frame real da gravação PTY na paleta escura, com banner e rodapé fixos.

![GIF animado da sessão real com tema escuro em inglês.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-dark-en.gif)

GIF animado da sessão real com tema escuro em inglês. O fundo do renderizador ilustra a paleta; o Mudarro não altera configurações de fundo do terminal.

[Ver ambos os temas, detalhes e reprodução](ui-shell.md#telas-gravadas).

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.


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
name: mudarro
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

Checkpoint anterior da UI configurável: **74 testes,76,38% cobertura de statements**, race/vet/build Linux/cross-build Darwin passaram. [Escopo e limites](coverage.md).

- [Shell de terminal permanente](ui-shell.md)
