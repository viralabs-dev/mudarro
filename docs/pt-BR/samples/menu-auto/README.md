[English](../../../samples/menu-auto/README.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Sample de criação automática do menu

Medição em 2026-10-06 no Pop!_OS 24.04, checkout `main @ 049a405a3c19369e92bebfad468e7f36d5ae0db4`. CLI compilada deste checkout; conta pessoal GitHub verificada: `daneiel`. O chat local “Criar CLI universal de projetos” foi consultado e estava concluído.

O sample Node não possui dependências externas: `app.js` imprime `sample ready` e permanece ativo até a parada. `mudarro.yaml`, `menu.sh`, wrappers e manifest nesta pasta são saídas reais de `init --select app-javascript:hello` e `generate`. Nenhuma operação de instalação de dependências é acionada pelo teste.

## Reprodução sem modificar o sample

Na raiz do repositório:

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go build -o /tmp/mudarro-validation-bin/mudarro ./cmd/mudarro
bash scripts/validate-menu.sh /tmp/mudarro-validation-bin/mudarro
```

Requisitos: Go para compilar, Bash, Node, npm, `realpath`, `sha256sum`, `find`, `cmp` e `ps`. O script copia os dois arquivos de entrada para diretório temporário com espaço no nome, valida os launchers em outro cwd e remove exclusivamente os recursos que criou. Ações do menu: aplicação, dependências, infraestrutura, qualidade e scripts. `install` está disponível no menu, mas não é executado pelo teste.

## Matriz efetivamente executada

| Fluxo / risco | Verificação | Resultado |
|---|---|---|
| Detecção automática e seleção | scan JSON → init seleciona hello | Passou, CLI real |
| Init repetido | recusa e hash da configuração preservado | Passou |
| Dry-run | menu inexistente após planejamento | Passou |
| Geração repetida | SHA-256 de todos os arquivos sem diferenças | Passou |
| Launcher gerado | Bash, cwd /tmp, caminho com espaço | Passou |
| Menu | entrada inválida, serviços, grupos, voltar/sair | Passou |
| Ação via submenu | qualidade → test → saída sample quality OK | Passou |
| Wrappers gerados | test e hello invocam comandos reais npm | Passou |
| Supervisor local | up duas vezes, status, logs, restart, down | Passou |
| Ação desconhecida | código de saída diferente de zero | Passou |
| Arquivo editado | geração recusa; launcher e manifest preservados | Passou |
| Preflight | launcher editado / JSON inválido não criam ação parcial | Passou, regressão Go |
| Seleção inválida / EOF | init não salva configuração parcial | Passou, regressão Go |
| JS/TS, Python, Go; ambiguidade, symlinks, configuração, argumentos e reset | suíte Go existente + regressões, race | Passou; 62,8% statements em internal/mudarro |
| Análise estática e shell | go vet; bash -n nos launchers e script | Passou |
| Instalador | download simulado, repetição, checksum inválido | Passou |
| Docker Compose | up/status/logs/restart/down, Alpine | Passou nesta rodada |
| Sem rede | Docker --network none: scan/init/generate repetido | Passou nesta rodada |
| Kubernetes | kind exclusivo, ciclo e volume montado | Passou; PVC Bound, UID estável e conteúdo após down/up/restart |
| Podman/macOS/WSL | execução nesta máquina | Não executado; Podman ausente, macOS/WSL indisponíveis |
| PostgreSQL/SQLite × ferramentas | nova rodada da matriz de bancos | Passou 8/8 reutilizando dependências existentes; migrations/seed/repetição reais, sem instalação |
| Dockerfile isolado | build, up repetido, stop/resume do mesmo ID, logs, restart, proprietário externo | Passou, scripts/integration-dockerfile.sh |
| Helm | kind exclusivo, upgrade/status/logs/restart/down | Passou; release mantida, PVC Bound/mesmo UID/conteúdo após down/up/restart (`helm-persistence.txt`) |
| PostgreSQL nativo | initdb/pg_ctl/postgres no host | Não executado: ferramentas ausentes; não instaladas |
| GIF / frame PNG | sessão PTY real renderizada, decode e frames compostos | Passou: menu-real.gif e menu-real-light.gif, 983×739, 7 frames, 15,05 s |
| Screenshot gráfica | janela real do menu | Pendente: portal cancelado; sessão área4/tela notebook já confirmada |

A suíte original não encontrou falhas. As novas regressões SOLID demonstraram inspect fora da interface Executor e panic em args vazios antes da correção; ambas foram corrigidas e a suíte final passou. Rejeições de init repetido, ação inválida e edição manual são os resultados esperados. O primeiro comando Go bloqueou por cache read-only; `GOCACHE=/tmp/mudarro-go-cache` resolveu. Docker/API GitHub foram acessados com execução local autorizada fora do sandbox.

## Evidências reais

- `evidence/menu-session.txt`: saída da sessão PTY genuína, incluindo comandos, menu e erros esperados; saída final 0.
- `evidence/menu-session.timing`: tempos/bytes da captura `script`.
- `evidence/go-tests.txt`: execução verbose de Go/race/coverage.
- `evidence/go-vet.txt`: saída vazia com retorno 0.
- `evidence/generated-files.txt`: caminhos efetivamente gerados no sample.

Reproduzir a gravação:

```bash
scriptreplay -T docs/samples/menu-auto/evidence/menu-session.timing \
  -O docs/samples/menu-auto/evidence/menu-session.txt
```

A gravação foi aberta com `scriptreplay` e conferida até `Automatic menu end-to-end: OK`. `visual-session.*` contém sessão PTY distinta, pausada para leitura, com ASCII/cores e execução real de test via submenu. `menu-real.gif` e `menu-real-light.gif` foram renderizados com agg oficial 1.9.0 autorizado em /tmp (sem sudo), decodificados e inspecionados nos temas escuro/claro. `frame-menu-recording.png` é um frame composto fiel da gravação, não screenshot gráfica. `media-validation.txt` contém dimensões/frames/duração.

O executor CUA retornou apps/browsers vazios. COSMIC Wayland/cosmic-screenshot existem; a captura anterior foi cancelada no portal e nenhum PNG de tela foi produzido. Execuções visíveis futuras somente na área de trabalho 4 da tela integrada do notebook, com confirmação do usuário de janela/área/tela. Não reabrir seletor antes de novo pedido. Gravações headless não abrem janelas.

Regravar/renderizar (requer agg já aprovado/disponível):

```bash
env -u NO_COLOR script -q -e -O docs/samples/menu-auto/evidence/visual-session.txt \
  -T docs/samples/menu-auto/evidence/visual-session.timing \
  -c 'bash scripts/record-menu.sh /tmp/mudarro-validation-bin/mudarro'
python3 scripts/typescript-to-cast.py
/tmp/mudarro-media-tools/agg --font-family 'DejaVu Sans Mono' --font-size 16 \
  --theme github-dark --fps-cap 10 --last-frame-duration 3 \
  docs/samples/menu-auto/evidence/visual-session.cast \
  docs/samples/menu-auto/evidence/menu-real.gif
```

NO_COLOR é removido apenas no processo de demonstração para mostrar as cores solicitadas; a preferência global é preservada. [Visual/fontes](../../../terminal-visual.md), [SOLID/decisão](../../../architecture.md) e [avaliação Go](../../../language-decision.md) descrevem tradeoffs.

Logs finais de integração: `smoke.txt`, `installer.txt`, `docker-compose.txt`, `offline.txt`, `kubernetes.txt`, `dockerfile.txt` e `helm.txt`. Scripts iniciais finalizaram com retorno 0; recursos exclusivos foram removidos. PVCs iniciais permaneceram Pending: comprovavam objeto existente, não dados. A rodada ampliada manifests em `kubernetes-persistence.txt` comprova Bound, mesmo UID e conteúdo após down/up/restart. A primeira fixture ampliada falhou YAML antes da execução; erro corrigido, registro preservado em `kubernetes-persistence-fixture-failure.txt`.

CI de referência confirmado: https://github.com/viralabs-dev/mudarro/actions/runs/37405606285 (`success`, SHA `f62ce9d7cb1e3101dfb38bc2e0729bec1f00812c`). Esse CI precede as alterações locais desta rodada e não comprova sua validação remota. Nenhum push, merge ou deploy foi feito.

Binário Linux histórico v1 validado (SHA-256): `33acac5d061ba7944990e133a75de9b98df75534080042617a077ef39c2a62d8`; estaticamente ligado. Testes finais e logs medidos em 2026-10-06; alterações locais não commitadas.

## Checkpoint visual final

Race/vet e E2E repetidos sobre o código final: `go-tests-v2.txt` (62,8%), `go-vet-v2.txt` (retorno 0), `menu-v2.txt` (OK). GIF real `menu-real-v2.gif`, frame composto `frame-menu-v2-recording.png`, validação `media-v2-validation.txt`; V1 em `evidence/v1/`. SHA256 do binário checkpoint final v3: `e9a25a2f584a76dd1a145e7bcbc9f4ed958793b407f20d0e14e6cd253012b188`.

Na sessão já confirmada área 4/notebook:

```bash
env -u NO_COLOR TERM=xterm-256color bash /tmp/mudarro-checkpoint-v3/scripts/demo-menu.sh /tmp/mudarro-checkpoint-v3/bin/mudarro
```

Fixture exclusiva, removida na saída; não disputa containers/clusters. Alternativa durável: construir binário e executar `scripts/demo-menu.sh /caminho/absoluto/mudarro`. Sequência: 1 serviço, 4 qualidade, 1 test, Enter, 0 voltar. `q` volta/sai. Screenshot gráfica continua pendente; posição já foi confirmada, sem novo pedido de posição.

## Integrações ampliadas e reprodução sem instalações

Bancos são relevantes aos itens do grupo banco e wrappers gerados; executar migrations não é requisito para scan/init/generate. A rodada ampliada verificou 8/8 PostgreSQL/SQLite × Goose/Alembic/Django/Prisma, após descobrir dependências existentes em /tmp fora do PATH. Logs: `databases-sqlite.txt`, `databases-postgres.txt`, `databases-sqlite-data.txt`, `databases-postgres-data.txt`. Alembic/Django: uma linha real após seed repetido; Goose/Prisma: hook de seed, sem inserção afirmada. `databases-reuse-runner.txt` valida modo de reprodução sem download. Primeira conexão PostgreSQL foi bloqueada pelo sandbox; retry escalado autorizado passou, container próprio removido.

```bash
PATH="/tmp/mudarro-test-tools:$PATH" \
MUDARRO_PYTHON_ENV=/tmp/mudarro-prisma-debug/python-env \
MUDARRO_JS_DEPS=/tmp/mudarro-prisma-debug/js-deps \
bash scripts/validate-databases.sh /tmp/mudarro-checkpoint-v2/bin/mudarro
```

Fornecer dependências equivalentes existentes se /tmp não persistir. Sem ambas variáveis, o runner original instala pacotes e exige autorização correspondente. PostgreSQL17 roda em container exclusivo com porta localhost efêmera, não substitui PostgreSQL nativo host.

Podman real: binário ausente; instalar requer autorização, portanto não executado. `podman-mock.txt` valida escolhas owned/foreign/absent sem alegar runtime. WSL/macOS não disponíveis nesta máquina. Unicode CJK/emoji pode exceder largura estimada por runes.

```bash
PATH="$HOME/.local/bin:$PATH" bash scripts/validate-kind.sh /tmp/mudarro-checkpoint-v2/bin/mudarro
PATH="$HOME/.local/bin:$PATH" bash scripts/integration-helm.sh /tmp/mudarro-checkpoint-v2/bin/mudarro
```

Frentes coordenadas: bancos com recursos próprios, revisão de código e auditoria somente leitura; integrador único para repo/Vault. Plano na hierarquia `06-backlog/Plano de execução paralela (Mudarro).md`. Aceite visual do usuário recebido; PNG gráfica permanece pendente.

Snapshot v2 permanece preservado; bancos/PVC usaram seu binário, pois correção posterior afeta somente truncamento de erro no menu. Snapshot final v3 teve race/vet/E2E e gravação refeitos; sem alterações visuais concorrentes.

Fechamento ampliado: `helm-persistence.txt` também passou Bound/UID/conteúdo. GIF final regenerado após correção estreita: 983×739, 7 frames, 15,04s (`media-v2-validation.txt`). `parallel-review.txt` registra revisão/auditoria reais. Snapshot v3 para demo; snapshot v2 preservado.

## Fechamento de pacotes e suporte

Extrações por etapa passaram; final `packages-final-tests.txt`, `packages-vet.txt`, `packages-coverage.out`, `packages-coverage-functions.txt`: race/vet e **72,0% agregados**, builds Linux/Darwin (cross-build), E2E Node/Python/Go/pnpm. Logs `packages-menu-real.txt`, `packages-languages-real.txt`, `packages-pnpm-real.txt`, `packages-containers.txt`, `packages-dockerfile.txt`, `packages-kustomize.txt`, `packages-helm.txt`, `packages-databases.txt`. Kustomize/Helm passaram Bound/mesmoUID/conteúdo;8/8bancos passaram. Recursos próprios removidos. Nenhuma instalação adicional.

Colisão `check:fast`/`check.fast` agora exige comando explícito: `packages-collision-before.txt` falha esperada no baseline; `packages-collision-after.txt` passou. [Arquitetura](../../../architecture.md), [matriz](../../../support-matrix.md), [roadmap](../../../language-roadmap.md) e plano Vault MUD014..016 conciliam alcance/limites. Yarn/uv/poetry/Podman runtime ausentes; testes de contrato não equivalem à execução dessas ferramentas.

Snapshot final v4, anteriores preservados:

```bash
env -u NO_COLOR TERM=xterm-256color bash /tmp/mudarro-checkpoint-v4/scripts/demo-menu.sh /tmp/mudarro-checkpoint-v4/bin/mudarro
```

SHA256 `76ccc665970c48a7f7c95f9a9052599cd23aee4580060b11a265a5b1bd91d942`. Gravação real refeita após código estável em `evidence/packages-visual/`, frame composto identificado como tal; screenshot gráfica ainda não produzida. Build não representa runtime macOS. Bancos/infra testaram snapshot b21c91f da extração; última alteração afetou somente colisão JS/nomes de campos, com E2E/contratos finais repetidos.

## Gate final dos testes — 2026-10-06

As 41 funções Test passaram com race e instrumentação explícita coverpkg após migração. Cobertura deduplicada: **1.035/1.427 statements = 72,53% (Go exibe 72,5%)**; anterior 1.028/1.427=72,04% permanece histórica. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e diff check passaram. Compilação Darwin não representa runtime macOS. Árvore de testes contém 16 arquivos: nove testes da orquestração e um helper; dois testes terminal e quatro helpers PTY. Sem hooks de teste ou overlays em produção.

Evidências: [test-tree-tests.txt](../../../samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](../../../samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](../../../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](../../../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](../../../samples/menu-auto/evidence/test-tree-gates.txt).

## Identidade do projeto consumidor — MUD-022

O título do menu usa `Config.Name` do projeto consumidor, não o nome da ferramenta Mudarro. O sample controlado agora usa Aurora (package/config); demo e gravação definem Aurora explicitamente por padrão, com override `MUDARRO_SAMPLE_NAME`. A CLI já usava Config.Name; as entradas de demonstração foram corrigidas. Aceite visual do usuário recebido; somente screenshot gráfica genuína continua pendente.

A sanitização remove controles e caracteres Unicode Cf de formatação. Estimativa conservadora de células considera CJK/fullwidth/emoji duas células e marcas combinantes zero; não garante tratamento perfeito de grafemas/largura em todos os terminais. Glyphs bitmap não suportados ou wordmark longo demais usam fallback textual legível, preservando identidade correta.

A etapa atual tem 43 funções Test e 24 casos nome/modo rich/narrow/plain/NO_COLOR: AtlasAPI, Aurora, Café, 漢字😀, nomes longos e injeção de controles, além de integração CLI com nome consumidor explícito. Race passou; cobertura final e gates de build passaram, registrados abaixo. Medições anteriores permanecem históricas. Mídias genuínas atuais ficam em `evidence/project-name-visual/` no sample canônico; `evidence/packages-visual/` permanece histórico.

## Validação atual do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../../../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../../../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../../../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../../../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../../../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../../../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Checkpoint atual com nome consumidor: `/tmp/mudarro-checkpoint-v5`; anteriores preservados.

```bash
env -u NO_COLOR TERM=xterm-256color bash /tmp/mudarro-checkpoint-v5/scripts/demo-menu.sh /tmp/mudarro-checkpoint-v5/bin/mudarro
```

## Gravações atuais da UI configurável

Checkpoint final:74 testes;1.552/2.032=76,38% statements agregados (Go76,4%),480 não executados; terminal379/422=89,81%. Race/vet/build Linux/cross-build Darwin e gates reais UI/linguagens/pnpm/Dockerfile/bancos passaram. SHA-256 binário: `260fb0f6703fda29ef02040a0e4e691e5c3a3233f908f5f1bd89fd27d67e4053`. Etapas históricas preservadas.

Preview genuíno: escuro inglês983×739/17frames/12,90s; claro pt-BR983×739/16frames/12,90s. Cast9,917s de saída; renderização acrescenta duração final. Labels dos frames compostos inspecionados. Eventos teclado/mouse injetados no PTY exercitam parsing, sem afirmar aceite mouse em emulador ou screenshot gráfica. Fixtures de scripts longos não executados: aurora-dark-en/aurora-light-ptbr no sample canônico. Preview não os executa; fixtures verificam SHOULD_NOT_EXECUTE ausente.

Reproduzir em diretório temporário isolado, com binário real compilado:

```bash
python3 scripts/record-ui-preview.py /absolute/path/to/mudarro /tmp/mudarro-ui-preview-reproduction
/tmp/mudarro-media-tools/agg --theme asciinema \
  /tmp/mudarro-ui-preview-reproduction/aurora-dark-en.cast \
  /tmp/mudarro-ui-preview-reproduction/aurora-dark-en.gif
/tmp/mudarro-media-tools/agg --theme solarized-light \
  /tmp/mudarro-ui-preview-reproduction/aurora-light-ptbr.cast \
  /tmp/mudarro-ui-preview-reproduction/aurora-light-ptbr.gif
```

Executar da raiz do repositório. agg é renderizador temporário já autorizado; comandos pressupõem disponibilidade. Fundos vêm explicitamente dos temas de renderização asciinema/solarized-light para demonstrar contraste das paletas foreground. Aplicativo não altera fundo do emulador/settings globais.

[aurora-dark-en.gif](../../../samples/menu-auto/evidence/ui-configuration/aurora-dark-en.gif), [aurora-light-ptbr.gif](../../../samples/menu-auto/evidence/ui-configuration/aurora-light-ptbr.gif), [frame-dark-en.png](../../../samples/menu-auto/evidence/ui-configuration/frame-dark-en.png), [frame-light-ptbr.png](../../../samples/menu-auto/evidence/ui-configuration/frame-light-ptbr.png), [languages-real.txt](../../../samples/menu-auto/evidence/ui-configuration/languages-real.txt), [pnpm-real.txt](../../../samples/menu-auto/evidence/ui-configuration/pnpm-real.txt), [databases-real.txt](../../../samples/menu-auto/evidence/ui-configuration/databases-real.txt).
