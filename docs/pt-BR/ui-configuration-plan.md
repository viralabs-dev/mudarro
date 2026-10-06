## Gate atual da UI configurável

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

**74 funções Test passaram com race; 1.552/2.032 statements executados = 76,38% (Go exibe76,4%)**, 480 não executados. Terminal:379/422=89,81%. Vet, build Linux e compilação cruzada do binário inteiro Darwin arm64 passaram; compilação não valida runtime macOS. E2E UI genuíno passou propagação de config personalizada nos wrappers/supervisor, scan traduzido, preview somente leitura/redação de segredos e smoke/menu. Rodadas reais Python/Go e pnpm passaram; Dockerfile e oito combinações de bancos com dependências existentes também passaram. Números anteriores abaixo são históricos com denominadores distintos. Fonte/aceite mouse em emulador real pendentes.

Evidências: [tests-race.txt](../samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](../samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](../samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](../samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](../samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](../samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](../samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](../samples/menu-auto/evidence/ui-configuration/ui-real.txt).

Estado: config/i18n/tema/preview implementados; fonte e aceite mouse em emulador real pendentes. Requisitos originais abaixo são histórico de desenho; orientação atual: [configuração UI](ui-configuration.md). Gates consolidados passaram; medição atual registrada abaixo.

[English](../ui-configuration-plan.md) · Português (Brasil)

# Plano de UI configurável — MUD-023..028

Plano aprovado, execução autorizada pelo usuário, incluindo implementação, commit e push. Funcionalidades centrais implementadas conforme resumo acima; pendências separadas. Baseline funcional: checkpoint v5, menu com nome do consumidor, 43 testes, cobertura de statements 73,26% e gravações Aurora validadas. CLI/YAML existentes e guardas de hashes gerados são a referência de compatibilidade. Instalações e alterações de segurança/terminal não recebem autorização implícita.

## Arquitetura e ordem

Manter wrappers Bash finos: chamam a CLI Go; Bash não implementa segunda UI. Configuração/validação e planejamento permanecem na orquestração; traduções/renderização no terminal; entrada SO/ciclo de terminal isolados do estado determinístico da UI. Acrescentar builder de preview somente leitura sobre Action/Command existentes. Só execução explícita chama Runner; foco, roda, scrollbar e expandir/recolher nunca executam.

1. MUD-025: leitor unificado e contratos equivalentes YAML/JSON.
2. MUD-023/024: catálogo de mensagens e paletas semânticas, preservando modo linha.
3. MUD-026: estado puro de foco/preview/viewport, depois adapter raw-input Linux/Darwin isolado e mouse.
4. MUD-027: decisão de fonte específica do terminal, independente da operação CLI.

Não introduzir framework/dependência UI ainda. Primeiro estimar e validar cleanup raw-mode, parsing de eventos e viewport com spike PTY pequeno sem dependências. Se manutenção exigir biblioteca, apresentar dependências/licença/plataformas antes de adicioná-la. Execução seguirá etapas e gates; aprovação do plano não prova implementação nem validação.

## Contrato de configuração

Um Config interno com tags yaml/json e validação comum; evitar parsers com semânticas diferentes. Campos opcionais ui estendem versão1; configs existentes seguem legíveis pelo novo binário. Binários antigos estritos rejeitam campos novos: documentar rollback removendo ui, nunca reescrever arquivo do usuário silenciosamente. Mudança incompatível futura exige versão2/conversão explícita.

Defaults propostos: locale=en, theme=auto (preserva foreground/background do terminal), density=comfortable, lettering=auto, preview.enabled=true, preview.mouse=auto. Locales en/pt-BR; temas auto/light/dark; densidade comfortable/compact; lettering auto/ascii/text. Light/dark inicialmente selecionam paletas foreground contrastantes; alterar fundo do emulador é separado e nunca implícito. Font-family não é configuração genérica efetiva até existir integração suportada.

Seleção: novo --config explícito com caminho relativo guardado; caso contrário exatamente um mudarro.yaml ou mudarro.json. Nenhum: erro existente de config ausente. Ambos: erro de ambiguidade solicitando seleção, sem merge. Init mantém YAML default; formato explícito seleciona JSON. Wrappers gerados carregam escolha não-default em argumentos estruturados/quoting seguro. --json de saída permanece distinto de formato da configuração.

Precedência UI: override CLI explícito > arquivo selecionado > defaults. Sem novo config global nem adivinhação de locale por ambiente. NO_COLOR e capacidade TERM/TTY limitam cores/interatividade independentemente. Validar campos desconhecidos, chaves duplicadas, documentos múltiplos/trailing JSON, tipos/enums e limites tamanho/caminho/symlink identicamente. Encoding/json sozinho aceita duplicatas: rejeição exige validação deliberada. Campos opcionais vazios recebem defaults; valores inválidos não vazios geram erros acionáveis. Salvamento preserva arquivos do usuário/guardas de hashes.

Exemplos equivalentes aceitos pela implementação atual:

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

Português usa locale pt-BR no arquivo selecionado. IDs, grupos, tokens de confirmação e argv estáveis; traduzir labels/help/erros do Mudarro. Saída externa preservada literalmente. Inglês vira default UI só no checkpoint de implementação, não retroativamente nas gravações atuais.

## Preview somente leitura e entrada

Mostrar cwd, argv estruturado (argumentos separados), shell explícito e wrapper/origem quando aplicável. Preservar quoting real, sem avaliação shell. Ações especiais de lifecycle/banco mostram sequência planejada rotulada com placeholders runtime; não inventar comando final nem inspecionar processos/containers para preencher preview. Ler arquivos guardados/limitados apenas quando associados explicitamente à ação. Nunca source, expandir valores env ou carregar .env no preview. Mostrar ${ENV_NAME}; redigir flags de segredo/credenciais URL reconhecidas e campos sensíveis configurados. Segredos arbitrários em shell livre não são inferíveis com confiabilidade: documentar limite, usar fixtures sem segredos nas gravações.

Na visão interativa proposta, números/Tab/Shift-Tab mudam foco; Down expande preview e Up recolhe. PageUp/PageDown/Home/End/roda scrollam preview; scrollbar desenhada permite clique/drag quando há reports mouse. Enter executa explicitamente; confirmação destrutiva continua obrigatória. Foco/toggle/drag nunca executam. Mudanças de entrada exigem ajuda visível/testes; modo numérico linha legado continua compatível.

Raw mode exige TTY entrada/saída e capacidades suportadas, independente de NO_COLOR. Preservar/restaurar termios/cursor/mouse na saída normal, EOF, erro, interrupção e antes de passar stdin ao comando; SIGKILL não permite cleanup. Retomar depois do comando. Desativar tracking fora do preview e restaurar estado anterior quando disponível. TERM sozinho não prova mouse. Parsing limitado e opt-out explícito evitam travar consultas sem suporte. Pipes/TERM dumb usam comandos preview modo linha/texto, sem escapes cursor/mouse. Janelas estreitas usam texto/viewport compacto; NO_COLOR mantém teclado em monocromático.

Mouse segue protocolos do terminal, não scrollbar de navegador. Xterm documenta coordenadas/press-release SGR e modos tracking; roda/clique/drag exigem PTYs reais e pelo menos um emulador real antes de alegar aceite mouse real. [Sequências Xterm](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html).

## Propriedade da fonte e decisão

CLI genérica não escolhe universalmente fonte instalada: emulador controla fonte. Nenhum executável kitty/wezterm encontrado no PATH desta máquina; fc-list existe para inventário futuro somente leitura. Isso não identifica emulador da sessão do usuário. Primeiro oferecer densidade/ascii/text, sem instalar fontes/alterar settings.

Integração futura opcional pode lançar emulador escolhido explicitamente com opções locais do projeto, sem config global. WezTerm oferece config-file/config/font-family; Kitty configuração de fontes, enquanto remote font-size altera tamanho nas subjanelas da janela SO. Contratos específicos, não troca universal de família. Não habilitar permissões remote-control nem abrir janela GUI sem confirmar integração/alvo. Fontes: [WezTerm CLI](https://wezterm.org/cli/general.html), [fonte WezTerm](https://wezterm.org/config/lua/wezterm/font.html), [config Kitty](https://sw.kovidgoyal.net/kitty/conf/), [remote-control Kitty](https://sw.kovidgoyal.net/kitty/remote-control/).

Decisão necessária antes de integração fonte: emulador do usuário e aceitação de janela separada lançada pelo projeto. Config/i18n/tema/preview podem prosseguir independentemente sob autorização já recebida.

## Aceite e coordenação

Testar paridade/defaults/duplicatas/desconhecidos/tamanho/symlink/ambiguidade; YAML antigo/hash; mensagens traduzidas/IDs estáveis; paleta plain/NO_COLOR/estreita; estado puro de preview sem chamadas Executor; PTY real teclas/mouse/resize/restauração/stdin filho; fixtures segredo/injeção e previews limitados. Mouse/drag em emulador e aplicação de fonte separados de parsing sintético/PTY; relatar distintamente. Cross-compile não é runtime plataforma.

Integrador possui produção compartilhada/Vault. Revisões somente leitura paralelas; edições em arquivos disjuntos atribuídos. Preservar v5, novo checkpoint por etapa funcional, registrar aprovado/falho/não executado; atualizar EN/pt-BR e Vault somente português. MUD-028 entrega plano aprovado; MUD-023..027 começam conforme ordem/gates e ownership registrados, sem alegar funcionalidades prontas. Usuário autorizou implementação, commit e push; integrador verifica gates/escopo antes da publicação. Instalações, mudanças de segurança e integração fonte continuam sem autorização implícita.
