[English](../language-decision.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# ADR-MUD-014 — Reavaliar Go para o Mudarro

Status: avaliação concluída; manter Go+Bash. Data: 2026-10-06. Não migra linguagem nem reescreve ADR-MUD-001; a escolha original foi recuperada do chat “Criar CLI universal de projetos” e das Decisões do Vault.

Go é uma escolha adequada aos requisitos atuais: CLI distribuída por binário, scanner/template offline, argumentos estruturados, YAML/TOML, supervisão e testes. A evidência favorece manter a implementação existente; não demonstra que Go é universalmente superior ou que outro protótipo teria desempenho pior. Não foram implementados/benchmarkados protótipos comparativos.

## Requisitos e alternativas

| Critério do projeto | Go atual | Rust | Python / Node | Bash puro |
|---|---|---|---|---|
| Instalação sem runtime do usuário | Binário Linux/macOS, parsers incorporados | Binário também atende | Interpretadores usuais ou empacotamento específico | Bash existente, parsers YAML/TOML viram dependências externas |
| Scanner/template sem rede | Implementado e testado offline | Viável, não prototipado | Viável; rede não é propriedade da linguagem | Viável, mas parsing robusto/manutenção ficam mais trabalhosos |
| Subprocessos/identidade/sinais | os/exec + supervisor Unix exercitado | Controle preciso possível, maior trabalho de reescrita atual | APIs de subprocesso disponíveis; empacotamento/runtime adicional | Fácil orquestração simples; quoting/estado/erros e testes difíceis nesta escala |
| Portabilidade real | Linux/macOS; Bash/ps/flock/Unix; WSL não testado | Linguagem não remove dependências Unix atuais | Linguagens não removem dependências Unix atuais | UNIX/WSL; Bash não oferece Windows nativo |
| Testes/manutenção | go test/race/vet, interfaces pequenas, testes de integração | Segurança de memória sem GC é vantagem em cenários específicos | Ergonomia/ecossistema úteis se equipe/produto já dependerem deles | Custo crescente com modelos/tipos/adaptadores e manifest SHA-256 |
| TUI e distribuição | Camada ANSI/ASCII pequena hoje; bibliotecas Go opcionais | Ecossistema TUI possível, migração não necessária para cores | Ecossistemas TUI/web úteis se produto mudar | Prompts simples suficientes; TUI complexa limita manutenção |

As avaliações de custo/manutenção são inferências do código/requisitos atuais. Não são medições universais de produtividade.

## Onde Go pode deixar de ser a melhor escolha

- Binário/runtime e memória: Go incorpora runtime/GC; se limites rígidos de tamanho, memória ou latência aparecerem, medir Go e protótipo Rust/C antes de decidir. Hoje o caminho crítico inclui filesystem e ferramentas externas, não foi demonstrado gargalo de CPU/GC.
- Windows nativo: o bloqueio é também o desenho Bash/Unix, não a linguagem Go. Exige novo launcher/supervisor e testes, mesmo mantendo Go ou migrando.
- Produto predominantemente web ou plugins do ecossistema Node/Python: pode justificar incorporar outro runtime ou separar componente; não precisa reescrever o scanner inteiro.
- Equipe exclusivamente especializada em outra linguagem: medir custo de manutenção e suporte antes de reimplementar. Competência e requisitos operacionais pesam mais que popularidade.
- Controle low-level/sem GC como requisito mensurável: Rust é alternativa plausível. Hoje isso não é requisito do menu/scanner.

## Decisão e limites

Manter Go pelo encaixe comprovado e custo de migração sem benefício demonstrado. Preservar Go+Bash e interfaces por consumidor; não adotar framework pesado apenas por estilo. O shell/interpreter das aplicações acionadas (Node, Python, banco, Docker, kubectl etc.) continua requisito da ação, não do scanner. Compilação cruzada Darwin arm64 passou nesta rodada; não é execução real macOS. WSL/Windows nativo continuam não validados.

Reconsiderar com novo requisito registrado no Kanban, baseline de tamanho/memória/tempo de scan/startup, protótipo mínimo da alternativa e teste operacional em sistemas alvo. Nenhuma migração ou publicação foi realizada.

## Fontes primárias

- [Go FAQ](https://go.dev/doc/faq): objetivos, compilação e runtime/GC (Go tem runtime; binário não significa ausência de runtime).
- [os/exec](https://pkg.go.dev/os/exec): execução de processos e CommandContext; o supervisor atual conserva protocolo Unix/sinais para sobreviver à CLI.
- [Rust Book — Hello World](https://doc.rust-lang.org/book/ch01-02-hello-world.html): compilação para executável distribuível, sem exigir Rust no destinatário.
- [Python zipapp](https://docs.python.org/3/library/zipapp.html): arquivo executável Python depende de interpretador disponível; módulos C têm limitações específicas.
- [Node single executable applications](https://nodejs.org/api/single-executable-applications.html): Node também oferece empacotamento em executável; não afirmar que sempre precisa instalar Node separadamente.

Evidências locais e limitações operacionais: [sample/matriz](samples/menu-auto/README.md), [arquitetura](architecture.md) e [visual](terminal-visual.md).
