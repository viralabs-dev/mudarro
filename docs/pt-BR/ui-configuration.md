[English](../ui-configuration.md) · Português (Brasil)

[Ir para as telas gravadas](#telas-gravadas).

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

## Gate atual da UI configurável

**74 funções Test passaram com race; 1.552/2.032 statements executados = 76,38% (Go exibe76,4%)**, 480 não executados. Terminal:379/422=89,81%. Vet, build Linux e compilação cruzada do binário inteiro Darwin arm64 passaram; compilação não valida runtime macOS. E2E UI genuíno passou propagação de config personalizada nos wrappers/supervisor, scan traduzido, preview somente leitura/redação de segredos e smoke/menu. Rodadas reais Python/Go e pnpm passaram; Dockerfile e oito combinações de bancos com dependências existentes também passaram. Números anteriores abaixo são históricos com denominadores distintos. Fonte/aceite mouse em emulador real pendentes.

Evidências: [tests-race.txt](../samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](../samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](../samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](../samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](../samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](../samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](../samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](../samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Configuração da UI

## Telas gravadas

Exemplo atual: projeto `mudarro`, definido por `Config.Name`, com capturas reais atualizadas em MUD-055. Esta revisão publica o exemplo atualizado a partir de `246048a`; a etapa anterior Aurora foi preservada como histórico.

Execução real do CLI em PTY, com entrada controlada. PNGs são frames renderizados da gravação; não são screenshots do desktop. O banner e o rodapé permanecem fixos enquanto o centro mostra menus, prévia e execução.

**Revisão gravada:** estas telas reais acompanham o shell permanente de MUD-032. As gravações anteriores foram preservadas; a validação abaixo foi executada localmente, independente da CI.

### Escuro · inglês

![Frame da saída real longa, com banner e rodapé preservados.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-dark-en.png)

Frame da saída real longa, com banner e rodapé preservados.

![GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-dark-en.gif)

GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.

### Claro · português

![Frame da saída real longa, com banner e rodapé preservados.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-light-ptbr.png)

Frame da saída real longa, com banner e rodapé preservados.

![GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.](../samples/menu-auto/evidence/ui-shell/mudarro-shell-light-ptbr.gif)

GIF animado da sessão real: prévia, execução, erro, entrada, cancelamento e retorno ao menu.



Implementados: Config interno único/validação comum YAML/JSON, ui opcional em versão1, mensagens próprias traduzidas, temas semânticos e preview de ações somente leitura. Integração de fonte não implementada; aceite mouse em emulador real não verificado. Evidências históricas preservadas; cobertura consolidada final registrada abaixo.

## Seleção e defaults

--config recebe caminho relativo guardado, inclusive nomes personalizados YAML/JSON. Sem opção, selecionar exatamente um mudarro.yaml/mudarro.json; ambos geram ambiguidade explícita, nenhum erro de arquivo ausente. Sem merge. Init usa YAML default; --format json seleciona JSON. --json existente é saída, não formato de configuração. Wrappers preservam caminho selecionado em argumentos com quoting seguro; supervisor local propaga seleção.

Overrides CLI prevalecem sobre arquivo/defaults. Sem config global ou adivinhação locale pelo ambiente. Defaults: locale=en, theme=auto, density=comfortable, lettering=auto, preview.enabled=true, preview.mouse=auto. Locales en/pt-BR; temas auto/light/dark; densidade comfortable/compact; lettering auto/ascii/text. Validação rejeita desconhecidos/duplicatas/trailing/documentos múltiplos/tipos/enums inválidos, mantendo guardas tamanho/caminho/symlink. Binários antigos estritos rejeitam ui; rollback exige remoção explícita, nunca reescrita silenciosa.

```yaml
version: 1
name: mudarro
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

IDs, chaves/grupos, argv e tokens destrutivos estáveis. Saída de processos externos não traduzida. Temas selecionam foreground semântico contrastante sem alterar fundo do emulador; NO_COLOR remove cores, mantendo teclado quando interativo suportado. CLI não seleciona fonte universalmente; não implementada instalação de fontes/mudança global de terminal.

## Preview e entrada

Raw restrito à visão de ações; seleção serviço/grupo permanece modo linha compatível. Números/Tab/Shift-Tab focam; Enter executa; Down expande e Up recolhe. PageUp/PageDown/Home/End/roda scrollam; scrollbar trata clique/drag de eventos reportados. Foco/expansão/scroll nunca chamam Runner. Confirmação destrutiva obrigatória. Fallback linha usa pN para preview e número comum para execução.

Preview mostra cwd, argv/shell/origem/sequências especiais planejadas rotuladas. Nunca executa, inspeciona processos/containers, source scripts, carrega .env ou expande valores env. Referências env permanecem placeholders; flags secretas/credenciais URL reconhecidas redigidas. Segredos arbitrários em shell não são inferíveis; gravações sem segredos. Leitura guardada/limitada.

Raw exige TTY entrada/saída suportados; restaura saída/EOF/erro/interrupção e antes de passar stdin ao filho, retomando depois. SIGKILL não permite cleanup. Pipes/TERM dumb plain sem modos mouse/cursor. Estimativa conservadora de células Unicode não garante grafemas perfeitos. Eventos mouse exercitados em PTYs isolados, sem aceite em emulador real; screenshot gráfica pendente. Ver [plano/limites](ui-configuration-plan.md), [configuração](configuration.md), [visual](terminal-visual.md).

## Prévia gravada anteriormente

Checkpoint anterior à moldura permanente; preservado como histórico.

![Prévia gravada anteriormente — PNG](../samples/menu-auto/evidence/ui-configuration/frame-dark-en.png)

![Prévia gravada anteriormente — GIF](../samples/menu-auto/evidence/ui-configuration/aurora-dark-en.gif)
