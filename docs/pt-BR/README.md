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

Não exige Go na máquina do usuário. Todo instalador baixa uma release do GitHub por HTTPS, confere o SHA-256 contra o `checksums.txt` e não instala nada se houver divergência.

| Plataforma | Arquiteturas | Instalador | Pacote da release | Situação |
|---|---|---|---|---|
| Linux | amd64, arm64 | `install.sh` | `mudarro_linux_<arch>.tar.gz` | Suportado |
| macOS | amd64 (Intel), arm64 (Apple silicon) | `install.sh` | `mudarro_darwin_<arch>.tar.gz` | Suportado |
| Windows 10/11 | amd64, arm64 | `install.ps1` (PowerShell 5.1 ou 7+) | `mudarro_windows_<arch>.zip` | Novo: `mudarro.exe` nativo a partir da primeira release que publicar o zip; limites na [matriz de suporte](support-matrix.md) |
| WSL | amd64, arm64 | `install.sh` dentro da distribuição | pacote Linux | Suportado (comporta-se como Linux) |

### Linux

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh | bash
export PATH="$HOME/.local/bin:$PATH"
```

Para fixar uma versão ou mudar o destino:

```bash
curl -fsSL https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.sh |
  MUDARRO_VERSION=v0.1.0 MUDARRO_INSTALL_DIR="$HOME/bin" bash
```

Atualize repetindo a instalação. Desinstale com `rm "$HOME/.local/bin/mudarro"` (e `rm -r "$HOME/.local/bin/mudarro-licenses"`); os arquivos dos projetos são preservados.

### macOS

O mesmo comando `install.sh` do Linux, para Intel e Apple silicon. O binário não é notarizado; como é baixado pelo `curl`, o macOS não marca o arquivo com quarentena e o Gatekeeper não o bloqueia. Se baixar o `.tar.gz` pelo navegador, retire a quarentena uma vez: `xattr -d com.apple.quarantine ~/.local/bin/mudarro`.

### Windows

No PowerShell (Windows PowerShell 5.1 ou PowerShell 7+), como usuário comum, sem administrador e sem mudar a política de execução:

```powershell
irm https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.ps1 | iex
```

Ou baixe o `install.ps1`, revise e execute como arquivo (o bypass vale só para esse processo):

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\install.ps1
```

- Instala o `mudarro.exe` em `%LOCALAPPDATA%\Programs\mudarro` e as licenças em `mudarro-licenses\`, ao lado.
- Acrescenta esse diretório ao `PATH` **do usuário** uma única vez (sem duplicar); abra um terminal novo depois.
- Executa `mudarro.exe version` antes de trocar uma instalação existente; falha de download, de checksum ou dessa verificação mantém a versão anterior.
- Variáveis: `MUDARRO_VERSION` (`latest` ou `vX.Y.Z`), `MUDARRO_INSTALL_DIR` (caminho absoluto) e `MUDARRO_REPOSITORY` (`dono/nome`). Exemplo: `$env:MUDARRO_VERSION = 'vX.Y.Z'; irm https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.ps1 | iex`.
- **Atualizar:** repita o mesmo comando.
- **Desinstalar:** `& ([scriptblock]::Create((irm https://raw.githubusercontent.com/viralabs-dev/mudarro/main/install.ps1))) -Uninstall` (ou `.\install.ps1 -Uninstall`). Remove o `mudarro.exe`, o `mudarro-licenses\` e a entrada do `PATH`; os arquivos dos projetos são preservados.
- **SmartScreen e antivírus:** o `mudarro.exe` não é assinado. O SmartScreen não interfere em arquivos baixados pelo PowerShell, mas pode avisar se o zip for baixado pelo navegador (*Mais informações → Executar assim mesmo*, depois de conferir o SHA-256 no `checksums.txt`). Alguns antivírus sinalizam binários Go não assinados; se a troca falhar, feche processos `mudarro` em execução ou verifique a quarentena.

A disponibilidade e os testes da versão estão em [Validação](validation.md).

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

## Licença

MIT — veja [`LICENSE`](../../LICENSE). As dependências embutidas mantêm as próprias licenças; veja [`THIRD_PARTY_NOTICES.md`](../../THIRD_PARTY_NOTICES.md) e [`LICENSES/`](../../LICENSES/).
