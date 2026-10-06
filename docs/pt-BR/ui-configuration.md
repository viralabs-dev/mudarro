[English](../ui-configuration.md) · Português (Brasil)

## Gate atual da UI configurável

**74 funções Test passaram com race; 1.552/2.032 statements executados = 76,38% (Go exibe76,4%)**, 480 não executados. Terminal:379/422=89,81%. Vet, build Linux e compilação cruzada do binário inteiro Darwin arm64 passaram; compilação não valida runtime macOS. E2E UI genuíno passou propagação de config personalizada nos wrappers/supervisor, scan traduzido, preview somente leitura/redação de segredos e smoke/menu. Rodadas reais Python/Go e pnpm passaram; Dockerfile e oito combinações de bancos com dependências existentes também passaram. Números anteriores abaixo são históricos com denominadores distintos. Fonte/aceite mouse em emulador real pendentes.

Evidências: [tests-race.txt](../samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](../samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](../samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](../samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](../samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](../samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](../samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](../samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Configuração da UI

Implementados: Config interno único/validação comum YAML/JSON, ui opcional em versão1, mensagens próprias traduzidas, temas semânticos e preview de ações somente leitura. Integração de fonte não implementada; aceite mouse em emulador real não verificado. Evidências históricas preservadas; cobertura consolidada final registrada abaixo.

## Seleção e defaults

--config recebe caminho relativo guardado, inclusive nomes personalizados YAML/JSON. Sem opção, selecionar exatamente um mudarro.yaml/mudarro.json; ambos geram ambiguidade explícita, nenhum erro de arquivo ausente. Sem merge. Init usa YAML default; --format json seleciona JSON. --json existente é saída, não formato de configuração. Wrappers preservam caminho selecionado em argumentos com quoting seguro; supervisor local propaga seleção.

Overrides CLI prevalecem sobre arquivo/defaults. Sem config global ou adivinhação locale pelo ambiente. Defaults: locale=en, theme=auto, density=comfortable, lettering=auto, preview.enabled=true, preview.mouse=auto. Locales en/pt-BR; temas auto/light/dark; densidade comfortable/compact; lettering auto/ascii/text. Validação rejeita desconhecidos/duplicatas/trailing/documentos múltiplos/tipos/enums inválidos, mantendo guardas tamanho/caminho/symlink. Binários antigos estritos rejeitam ui; rollback exige remoção explícita, nunca reescrita silenciosa.

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

IDs, chaves/grupos, argv e tokens destrutivos estáveis. Saída de processos externos não traduzida. Temas selecionam foreground semântico contrastante sem alterar fundo do emulador; NO_COLOR remove cores, mantendo teclado quando interativo suportado. CLI não seleciona fonte universalmente; não implementada instalação de fontes/mudança global de terminal.

## Preview e entrada

Raw restrito à visão de ações; seleção serviço/grupo permanece modo linha compatível. Números/Tab/Shift-Tab focam; Enter executa; Down expande e Up recolhe. PageUp/PageDown/Home/End/roda scrollam; scrollbar trata clique/drag de eventos reportados. Foco/expansão/scroll nunca chamam Runner. Confirmação destrutiva obrigatória. Fallback linha usa pN para preview e número comum para execução.

Preview mostra cwd, argv/shell/origem/sequências especiais planejadas rotuladas. Nunca executa, inspeciona processos/containers, source scripts, carrega .env ou expande valores env. Referências env permanecem placeholders; flags secretas/credenciais URL reconhecidas redigidas. Segredos arbitrários em shell não são inferíveis; gravações sem segredos. Leitura guardada/limitada.

Raw exige TTY entrada/saída suportados; restaura saída/EOF/erro/interrupção e antes de passar stdin ao filho, retomando depois. SIGKILL não permite cleanup. Pipes/TERM dumb plain sem modos mouse/cursor. Estimativa conservadora de células Unicode não garante grafemas perfeitos. Eventos mouse exercitados em PTYs isolados, sem aceite em emulador real; screenshot gráfica pendente. Ver [plano/limites](ui-configuration-plan.md), [configuração](configuration.md), [visual](terminal-visual.md).
