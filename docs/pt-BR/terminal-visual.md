[English](../terminal-visual.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Visual do terminal

Cores semânticas na saída TTY: ciano para marca/títulos, laranja em terminais 256 cores (amarelo básico) para escolhas/prompts, cinza para instruções secundárias, vermelho para erro. Texto e números carregam o significado mesmo sem cor. Fundo e texto normal usam padrões do terminal. A marca continua original, com lettering bitmap 5x7 original, profundidade e barra lateral em TTY colorido; fallback ASCII 3x5 embarcado, sem figlet/rede/dependência visual obrigatória.

A camada `terminal/terminal_style.go` só participa do menu; scan JSON, logs não interativos e arquivos gerados permanecem plain. `terminal/terminal_unix.go` consulta o stream de saída por ioctl: pipes e `/dev/null` não são TTYs. Cores desativadas com TERM vazio/dumb ou NO_COLOR não vazio (inclusive 0), conforme [NO_COLOR](https://no-color.org/). Não altera preferências globais. A largura do ASCII usa a janela ou COLUMNS explícito; o layout mantém ASCII como fallback. Títulos/caminhos longos são truncados à largura; caracteres Unicode de largura dupla ainda podem ultrapassar a largura estimada por runes.

## Referências públicas

Contagens aproximadas consultadas em 2026-10-06, variam com o tempo. Inspiração em hierarquia, espaços, escolhas e cores; nenhum código, logo, fonte ou marca copiado.

| Projeto | Estrelas aproximadas | Licença | Referência |
|---|---:|---|---|
| [Lip Gloss](https://github.com/charmbracelet/lipgloss) | 11,9 mil | MIT | Hierarquia e paleta semântica Go |
| [Gum](https://github.com/charmbracelet/gum) | 24,5 mil | MIT | Prompts e escolhas shell |
| [Fastfetch](https://github.com/fastfetch-cli/fastfetch) | 24,9 mil | MIT | Marca ASCII compacta |
| [ascii-image-converter](https://github.com/TheZoraiz/ascii-image-converter) | 3,5 mil | Apache-2.0 | Assets ASCII; não adotados glyphs Braille |
| [VHS](https://github.com/charmbracelet/vhs) | 21,1 mil | MIT | Demos de terminal; não é dependência |

O [agg oficial v1.9.0](https://github.com/asciinema/agg/releases/tag/v1.9.0) foi autorizado e instalado exclusivamente em `/tmp/mudarro-media-tools/agg`, sem sudo. Renderiza a gravação real PTY; o GIF não é screenshot de tela gráfica nem imagem gerada por IA.

## Validação visual

`menu-real-v2.gif`: renderização real conferida por decodificação GIF e inspeção de frames compostos; contraste/layout em fundo escuro. A gravação separa scan, geração, serviços/submenus, ação real test e proteção de arquivo editado. Sem marcas de terceiros ou conteúdo de outras aplicações. NO_COLOR/TERM=dumb, pipes e largura estreita têm testes específicos. O tema claro pode ser renderizado do mesmo cast para inspeção; preferências do usuário permanecem intactas.

Screenshot real permanece pendente: o portal anterior foi cancelado. Execuções visíveis futuras somente na área de trabalho 4 da tela integrada do notebook, na sessão já confirmada pelo usuário; não reabrir seletor antes de pedido. PTY headless não abre janela nem captura desktop.

## Revisão v2

Referência visual fornecida pelo usuário foi materializada e inspecionada. Lettering Mudarro original, barra laranja, contexto, menu e rodapé responsivos; `q` volta/sai. Em TTY colorido a tela redesenha e ações aguardam Enter; saída plain permanece sequencial. Layout estreito, políticas de cor e EOF têm regressões. `frame-menu-v2-recording.png` é frame composto fiel do GIF real. V1 preservada em `samples/menu-auto/evidence/v1/`.

Após extração em pacotes, nova sessão genuína em `samples/menu-auto/evidence/packages-visual/`; checkpoints anteriores preservados. Apresentação mantém contrato visual; aceite do usuário ainda pendente.

## Identidade do projeto consumidor — MUD-022

O título do menu usa `Config.Name` do projeto consumidor, não o nome da ferramenta Mudarro. O sample controlado agora usa Aurora (package/config); demo e gravação definem Aurora explicitamente por padrão, com override `MUDARRO_SAMPLE_NAME`. A CLI já usava Config.Name; as entradas de demonstração foram corrigidas. Aceite visual do usuário recebido; somente screenshot gráfica genuína continua pendente.

A sanitização remove controles e caracteres Unicode Cf de formatação. Estimativa conservadora de células considera CJK/fullwidth/emoji duas células e marcas combinantes zero; não garante tratamento perfeito de grafemas/largura em todos os terminais. Glyphs bitmap não suportados ou wordmark longo demais usam fallback textual legível, preservando identidade correta.

A etapa atual tem 43 funções Test e 24 casos nome/modo rich/narrow/plain/NO_COLOR: AtlasAPI, Aurora, Café, 漢字😀, nomes longos e injeção de controles, além de integração CLI com nome consumidor explícito. Race passou; cobertura final e gates de build passaram, registrados abaixo. Medições anteriores permanecem históricas. Mídias genuínas atuais ficam em `evidence/project-name-visual/` no sample canônico; `evidence/packages-visual/` permanece histórico.

## Validação atual do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Implementação atual: [configuração UI](ui-configuration.md). Config/i18n/tema/preview implementados; fonte/aceite mouse em emulador real continuam pendentes separados. Gates finais:74 testes,76,38% statements agregados; ver [cobertura](coverage.md).
