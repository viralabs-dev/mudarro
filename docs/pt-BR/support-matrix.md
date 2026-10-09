[English](../support-matrix.md) · Português (Brasil)

## Windows nativo — MUD-068 (substitui a exclusão da ADR-MUD-007)

O Mudarro agora compila e roda nativamente no Windows 10 1809+/Windows 11 (amd64; arm64 por cross-build). Nenhuma dependência de módulo nova: o código Windows usa só o pacote padrão `syscall`, pelo wrapper interno `winapi` sobre a `kernel32.dll`.

| Área | Comportamento no Windows | Limite |
|---|---|---|
| Build | `GOOS=windows go build ./...` e `go vet ./...` passam; código só-Unix foi para arquivos `*_unix.go` (`!windows`) ou `linux \|\| darwin` | Antes da MUD-068 o build Windows falhava em `syscall.Flock`, `syscall.Kill`, `Setsid` e `O_NOFOLLOW` |
| Ações do menu (`ContextOS`) | Cada comando nasce suspenso, em grupo de processos novo, entra num Job Object próprio e só então é retomado. O cancelamento (Ctrl+C) envia `CTRL_BREAK_EVENT` ao grupo e encerra o job após 300 ms, matando todos os descendentes | `CTRL_BREAK_EVENT` só alcança processos no mesmo console do Mudarro; o resto é encerrado pelo job, sem parada graciosa. Se o host negar jobs aninhados, só o filho direto é encerrado |
| Supervisor local (`local up/down/restart/status/logs`) | O supervisor roda destacado (`DETACHED_PROCESS`, grupo novo); o serviço roda num Job Object kill-on-close, então o `down` (que encerra o supervisor) derruba a árvore inteira. O lock usa `LockFileEx` | Não há SIGTERM: o `down` é parada forçada, e o serviço precisa tolerá-la. A identidade é a imagem do executável mais um mutex nomeado atrelado ao argv completo e ao token do supervisor (`Local\mudarro-supervisor-<sha256>`), e não o argv lido de `/proc`/`ps`; vale dentro da mesma sessão de logon |
| Menu interativo e moldura | Console cru via `SetConsoleMode` (`ENABLE_VIRTUAL_TERMINAL_INPUT`, sem eco/linha/entrada processada, Ctrl+C lido como tecla), saída VT via `ENABLE_VIRTUAL_TERMINAL_PROCESSING`, tamanho via `GetConsoleScreenBufferInfo`, espera de tecla via `WaitForSingleObject` + `PeekConsoleInput` | Exige console com VT (Windows 10+). Sem VT, a saída é tratada como não interativa e vale o menu numerado simples. O modo VT de saída fica ligado no console depois que o Mudarro sai |
| Cor e `TERM` | `NO_COLOR` e `TERM=dumb` continuam desligando cor e menu interativo. Sem `TERM`, o Windows Terminal (`WT_SESSION`) é tratado como `xterm-256color`, e o console clássico como `windows-console` | Mouse ligado por padrão só no Windows Terminal; nos demais, use `ui.preview.mouse: on` |
| Comandos | Comandos `args` rodam nativamente. Um prefixo `env NOME=valor` (Go com `go_workspace: off`) vira variável de ambiente quando não existe programa `env` | Comandos `shell:`, seeds/scripts `bash` e o `menu.sh` gerado ainda exigem `bash` no PATH (Git for Windows). O `C:\Windows\System32\bash.exe` é o WSL e roda os comandos dentro do Linux |
| Caminhos | Globs de exclusão sobre caminhos com `/` usam `path.Match` (mesmo resultado em todo SO) | A prévia de scripts recusa symlinks/reparse points via `Lstat` em vez de `O_NOFOLLOW`; FIFO não bloqueante é só Unix |

Validação registrada na MUD-068: `go test ./...` no Linux PASSOU; `go vet` PASSOU para linux, darwin e windows; builds `GOOS=windows`/`darwin`/`freebsd` PASSARAM; `gofmt -l` vazio. Os testes só-Windows (`*_windows_test.go`: E/S do executor, código de saída, cancelamento de árvore que ignora CTRL_BREAK, prefixo `env`, supervisor up/status/down e recusa de estado forjado) rodam apenas no job `windows-latest` do CI; como fumaça extra (não prova), os binários de teste Windows também passaram no Wine 10, exceto fixtures que o Wine não representa (symlinks invisíveis ao `Lstat`, nomes de diretório temporário não ASCII). Os testes só-Unix mantêm build tag `linux || darwin` (PTY, termios, FIFO, fixtures com `sh`). O aceite físico em console (conhost e Windows Terminal) continua pendente.

## Checkpoint vigente de runtimes autorizados

Toolchains privados aprovados em `/tmp/mudarro-approved-runtimes.6AsmDd` agora possuem validação local genuína: **Rust1.99.0: 46 checks** (build/test offline, binário escolhido42, falha controlada7, startup/restart/down); **JDK25.0.4.1+1/Maven3.10.0: 18 checks** (self-test Java real com biblioteca padrão42/falha7; Maven somente versão); **.NET10.0.401: 31 checks** (console self-test42/falha7, falha/recovery de build, zero pacotes NuGet; HOME preservado na execução final). [Evidências e reprodução runtime](../samples/menu-auto/evidence/next-languages/runtime/README.md).

Lifecycle/plugins Maven e Gradle continuam bloqueados por aprovação separada; compile/test Maven runtime não afirmados. Validação console .NET não comprova dotnet test/framework. PHP/Composer e Ruby/Bundler continuam não executados. Metadados framework e checks stdlib Node/Python não comprovam runtime Flask/FastAPI/Express. Os sete gates de produção passaram novamente com **235 grupos/83,06%**, sem fix de produção ou nova publicação. Conclusão Rust/.NET vale somente para subsets aprovados; não implica suporte completo de linguagem/framework.

## Checkpoint estático histórico — MUD-039/040/041/042/043/005

Adapters estáticos e verificações de metadados de frameworks implementados localmente. Gates finais: **235 grupos principais PASS, zero falhas/skips de testes**, quatro pacotes sem testes; cobertura canônica **3050/3672 = 83,06%**. Race, vet, build Linux, cross-build Darwin arm64, smoke, menu e resize terminaram com exit0. [Evidências finais](../samples/menu-auto/evidence/next-languages/README.md). Cross-build não comprova runtime macOS. Nenhum runtime Rust/Cargo, JDK/Maven, PHP/Composer, .NET ou Ruby/Bundler foi executado. Checks com bibliotecas padrão Node/Python existentes não comprovam runtime Flask/FastAPI/Express. Planos de instalação são somente pesquisa, sem autorização para instalar. [Evidências da preparação](../samples/menu-auto/evidence/next-adapter-preparation/README.md) · [Planos de instalação revisados](../samples/menu-auto/evidence/next-adapter-preparation/queue/install-plans/README.md).

| Capacidade | Contrato estático implementado | Limites restantes |
| --- | --- | --- |
| Rust/Cargo | TOML limitado, diagnósticos de workspace/membros/targets explícitos e contidos; build/test opt-in | Sem avaliar dependências/build scripts, resolver grafo/globs ou inventar startup; runtime pendente |
| Java/Maven; Gradle/Kotlin | XML/módulos Maven limitados; sugestões compile/test para projetos elegíveis; Gradle DSL somente evidência | Sem avaliar plugins/profiles/propriedades ou Gradle; main/start não inferidos; runtime JVM pendente |
| PHP/Composer | Metadados Composer limitados e sugestões de scripts explicitamente selecionados | Sem executar plugins/hooks no scan nem inferir startup de framework; runtime pendente |
| C#/.NET | XML limitado, sugestões build/test por projeto explícito; solution como evidência | Sem avaliar MSBuild, resolver grafo de solution ou startup; runtime pendente |
| Ruby/Bundler | Evidências limitadas Gemfile/gemspec/lock | Sem avaliar DSL Ruby/semântica lock, inferir Rake/Rails/start ou comandos automáticos; runtime pendente |
| Metadados frameworks | Evidência de dependências framework/test declaradas e contrato de entrada explícita | Sem importar módulos da aplicação no scan ou afirmar runtime de framework |

MUD-044/MUD-045 publicados em `c400f99`, autoria pessoal `daneiel`, 89 arquivos e `[skip ci]`; duas consultas autenticadas encontraram zero runs Actions dessa publicação. Seu checkpoint de199 testes/81,73% e samples reais Mix/nativos permanecem evidências históricas, não métricas deste lote estático. [Evidências publicadas Mix/nativos](../samples/menu-auto/evidence/mix-native/README.pt-BR.md).

Checkpoints anteriores abaixo preservam resultados e limites de publicação/runtime da época registrada.

## Checkpoint histórico — MUD-044/MUD-045

199racePASS/zero falhas-pulados/81,73%; Mixestático compile-test opt-in/umbrella declarado e C-C++customMake targetsusuário validados em samples internos.26Mixchecks+8rechecks/63comandosnativos+rechecks; nenhuma instalação/startinventado/publicação deste lote. MUD037 publicado6c868c95[daneiel/skipci/0Actions]. [Evidências e limites](../samples/menu-auto/evidence/mix-native/README.pt-BR.md).


A base publicada agora é `b887f64`, após reescrita autorizada da autoria pessoal de `56dadba`, com árvore idêntica. [Mapa de SHAs e provas de preservação](../samples/menu-auto/evidence/git-identity-repair/README.md). As novas mudanças de robustez/Unicode e esta referência permanecem locais.

## Checkpoint integrado mais recente — workspaces JavaScript

176 PASS race (173 grupos substantivos + três helpers),0 falhas/skips; **2381/2945 = 80,85%**. Herança sómanager do workspace declarado, install raiz, scripts/cwd próprios, globs/nested/exclusões/offline/conflitos e preservação manual passaram. Cinco managers existentes passaram85checksCLIreais, negativos e rechecks finais. MídiasPTY/samples salvos. Mudanças locais sobre44c4e04publicado, sem nova publicação. [Evidence / reprodução](../samples/menu-auto/evidence/javascript-workspaces/README.md).

## Checkpoint histórico165 — parser oficial Go workspace

165 PASS race (162 grupos substantivos + três helpers), 0 falhas/skips; **2199/2746 = 80,08%**. Parser oficial go.work/go.mod, diagnósticos de membros/conflitos/exclusões, contexto herdado e isolamento explícito passaram. Gates finais e rechecks reais Go/cinco runtimes instalados passaram. Quatro archives cross-build e fixtures do instalador preservam licenças. macOS nativo e aceite gráfico físico não executados. [Evidence / reprodução](../samples/menu-auto/evidence/go-work-parser/README.md).

## Checkpoint histórico156 — decisões managers/Go

156 PASS race (154 grupos substantivos + dois helpers), 0 falhas/skips; **79,47%**. SemVer exato, seleção start Go e isolamento workspace off passaram; preview/dry-run/doctor exibem/verificam isolamento. Gates locais e scans/qualidade dos cinco runtimes com binário atual passaram. Parser completo de membros go.work pendente; macOS nativo não executado. [Evidence / reprodução](../samples/menu-auto/evidence/manager-go-decisions/README.md).

## Validação runtime local mais recente — managers isolados

Pipenv2026.8.0, uv0.12.23, Poetry2.5.1 e YarnClassic1.22.22/Modern4.18.1 aprovados agora possuem validação genuína em fixtures sem dependências. Locks reais, menus/wrappers, idempotência, negativos de preservação e lifecycle passaram; start Python é explícito. Provas de instalação e mídias arquivadas. Código de produção não mudou;143 testes race/79,08% continuam o último checkpoint medido. [Evidence / reprodução](../samples/menu-auto/evidence/isolated-runtimes/README.md).

## Checkpoint local mais recente — timeout Compose no doctor

143 PASS race (141 grupos substantivos + dois helpers), 0 falhas/skips; **79,08%**. Sondagem Compose real expirou em 10,05s; configuração preservada e processo direto encerrado. Gates locais passaram; macOS nativo continua não executado. [Evidence / reprodução](../samples/menu-auto/evidence/doctor-compose-timeout/README.md).

## Checkpoint local anterior — doctor e arquivos especiais

141 PASS race (140 grupos substantivos + helper),0 falhas/skips; **79,03%**. Doctor rejeita arquivos especiais relativos com bits executáveis; CLI real antes/depois preservou fixtures. Vet/builds/smoke/menu/config/resize passaram; macOS nativo não executado, mudanças novas locais. [Evidence / reprodução](../samples/menu-auto/evidence/doctor-special/README.md).


## Checkpoint local anterior — nomes packageManager / IO fontes Go

140 PASS race (139 grupos substantivos + helper subprocesso), 0 falhas/skips; **2.058/2.604 = 79,03%**. Validação de nomes e leitura regular/limitada Go passaram; vet/builds/smoke/menu/config/resize passaram. Gramática versões, seleção/contexto Go aguardam decisão; runtime UV/Poetry ausente. Mudanças locais. [Evidence / reprodução](../samples/menu-auto/evidence/manager-go-source/README.md).


## Checkpoint local anterior — MUD-034/038 parciais

**131 testes principais race PASS, 0 falhas/skips; 2.041/2.588 = 78,86%.** Detecção/scripts/argv de banco Pipenv e regressões de entrada Go passaram; vet/build Linux/cross-build Darwin/smoke/menu/config/resize passaram. Runtime Pipenv ausente; membros/contexto/seleção targets go.work ainda não implementados, com aviso explícito. [Evidências e reprodução](../samples/menu-auto/evidence/pipenv-go/README.md). Mudanças locais.


## Checkpoint local anterior — MUD-033

121 testes principais com race PASS, zero falhas/skips; 1.971/2.517 statements internos = **78,31%**. Runtime genuíno Bun 1.4.0, menu/wrappers/supervisor e regressões bun.lock/bun.lockb/conflitos passaram sem instalação. Vet, build Linux, cross-build Darwin e gates smoke/menu/config/resize PTY passaram. [Evidências, sample salvo e reprodução](../samples/menu-auto/evidence/bun/README.md). macOS nativo e aceite gráfico/físico continuam pendentes. Mudanças novas locais sobre b887f64; MUD-035 somente análise.


## Último checkpoint local — MUD-053

117 testes principais racePASS/zero falhas/skips; cobertura deduplicada1.965/2.515=78,13%. U+3000 agora ocupa2células; clipping respeita1–3colunas semexpandir orçamento. 17variantes ePTYresize/NO_COLOR/restauração passaram, assim como vet/buildLinux/crosscompileDarwin eE2Essmoke/menu/config/resize. [Evidências, reprodução e limites](../samples/menu-auto/evidence/unicode-resize/README.md). Política conservadora por rune,sem shapingperfeito/grafemas indivisíveis;aceite físico026/027separado,macOSnativo031pendente. Mudanças locais sobre56dadba,sem novo commit/push/CI;histórico056aguardadecisão. [Triagem035somenteanálise](../samples/menu-auto/evidence/package-manager-contracts/README.md).

## Anterior checkpoint local de robustez — MUD-050/051/052

112 testes principais passaram com race, zero falhas/skips; nova cobertura interna 1.943/2.510 = **77,41%**. Supervisor Linux verifica inode real e argv exato; estado JSON inválido é rejeitado. Geração propaga erros de saída antes de escrever, e leitura de manifests rejeita arquivos especiais estáticos e limita a 4 MiB. Negativos do doctor passaram sem executar ações da aplicação. Vet/build Linux, crosscompile Darwin, smoke/menu/idempotência/preservação e E2E de configuração selecionada passaram. [Reprodução, logs e limites](../samples/menu-auto/evidence/robustness/README.md).

Trabalho local sobre base publicada56dadba, sem novo commit/push ou validação CI. macOS nativo segue não validado e seu fallback ps tem garantias de identidade menores. IO é atômico por arquivo, sem rollback coletivo ou proteção atômica contra troca concorrente de caminhos. Medições históricas abaixo mantêm seus denominadores originais.

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

## Gate histórico da UI configurável

**74 funções Test passaram com race; 1.552/2.032 statements executados = 76,38% (Go exibe76,4%)**, 480 não executados. Terminal:379/422=89,81%. Vet, build Linux e compilação cruzada do binário inteiro Darwin arm64 passaram; compilação não valida runtime macOS. E2E UI genuíno passou propagação de config personalizada nos wrappers/supervisor, scan traduzido, preview somente leitura/redação de segredos e smoke/menu. Rodadas reais Python/Go e pnpm passaram; Dockerfile e oito combinações de bancos com dependências existentes também passaram. Números anteriores abaixo são históricos com denominadores distintos. Fonte/aceite mouse em emulador real pendentes.

Evidências: [tests-race.txt](../samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](../samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](../samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](../samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](../samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](../samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](../samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](../samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Inventário Mudarro — leitura de código, 2026-10-06

Escopo: checkout pessoal, inventário conciliado após extração de pacotes. Nenhuma execução nova, instalação ou adapter implementado nesta frente. Suporte detectado não equivale a execução integral de todo ecossistema.

## Linguagens e gerenciadores

| Família | Evidência de detecção | Ações inferidas | Limites atuais |
|---|---|---|---|
| JavaScript / TypeScript | package.json; TS por tsconfig.json ou dependência typescript | install; start/dev/build aplicação, test/lint qualidade; demais scripts selecionáveis | npm default; npm/pnpm/yarn/bun por locks únicos ou packageManager conhecido explícito; bun.lock/bun.lockb deduplicados; nomes inválidos diagnosticados, versões informadas exigem SemVer exato; metadata apenas sintática, sem verificar hash; frameworks são rótulos |
| Python | pyproject.toml / requirements.txt / manage.py / Pipfile | pip venv/install, uv sync, poetry install; Pipenv install/sync e sugestões de scripts opt-in; Django runserver/test | Pipfile seleciona pipenv, conflitos pendentes; TOML/JSON inválidos diagnosticados, sem schema/hash completo; runtime Pipenv2026.8.0 sem dependências validado; não-Django exige start explícito; sem ações automáticas Flask/FastAPI/pytest |
| Go | go.mod; constraints do host e parser de func main ativo | build ./..., test ./..., mod download; start para única entrada ativa | fontes regulares limitadas4MiB; múltiplas/nenhuma entrada ficam pendentes; não infere tags custom/imports/compilabilidade; parser oficial go.work/go.mod valida membros bounded dentro da raiz; conflitos/exclusões diagnosticados; seleção promove start, contexto herdado por padrão e off explícito disponível |
| Elixir/Mix | mix.exs regular/limitado sem avaliar código | compile/test opt-in; mix_umbrella declarado explicitamente | sem inferir aliases/membros/grafo/deps/start;26checksroot/umbrella+rechecksreais |
| C/C++ / custom | evidências fontes/headers regulares, Make ancestral/owner explícito | targetsMake opt-in ou comandos usuário; source-only pending | sem parserfontes/flags/start, sem CMake/Meson;63comandosreais+rechecks |

Manager é string de configuração, não registro fechado validado. Entry scripts Python vêm de maps e ordenação final ocorre por serviço/nome. Scripts JS normalizam : e . para -. Colisões agora geram pendência e não oferecem seleção ambígua; regressão demonstrou falha antes e aprovação após correção.

## Monorepos

Walk recursivo detecta serviços por diretório com manifests. Ordenação de diretórios determinística; exclusões, ignorados e symlinks respeitados. Workspaces JS declarados resolvem ownership e herdam somente manager; instalação na raiz, scripts/cwd próprios dos membros, conflitos explícitos. Sem grafo Turbo/Nx ou herança de infra/banco/comandos. Go ignora módulos aninhados ao buscar main e detecta cada go.mod separadamente; go.work não produz plano de execução agregado. Sugestões Makefile/shell associadas ao serviço ancestral mais próximo, com fallback índice 0 quando fora de qualquer serviço.

## Infraestrutura

local supervisor; Docker e Podman Compose/Dockerfile; Kubernetes manifests/kustomize/Helm. Compose exige escolha docker/podman, múltiplos candidatos geram pendência. Compose prevalece sobre Dockerfile no mesmo diretório; candidatos k8s podem tornar ambíguo. Context/namespace Kubernetes explícitos. Infra só pesquisada junto ao serviço, não descoberta global de Helm chart aninhado sem serviço. Podman Compose usa podman-compose explicitamente. Kubernetes down escala workloads, não apaga PVC; preservar objeto não prova dados persistidos.

## Bancos

Kinds operacionais PostgreSQL/SQLite. Ferramentas Prisma, Django, Alembic, Goose; detecção estática schema.prisma/manage.py/alembic.ini/migrations SQL com -- +goose. Kind inferido apenas por strings Prisma; Django/Alembic/Goose exigem kind explícito. Múltiplas ferramentas geram pendência. Migration, criação de migration e seed separados; reset marcado destrutivo quando existe. Alembic não oferece reset automático. Templates não inventam modelos. PostgreSQL nativo, containers/Kubernetes e SQLite dependem das ferramentas/ambiente da ação; runtime de linguagem não é instalado pelo scanner.

Fontes locais: scan.go, adapters/language_adapters.go, actions.go, adapters/infrastructure_adapters.go, adapters/database_adapters.go, database_setup.go, config.go, test/internal/mudarro/core_test.go e scripts de integração.

# Matriz de cobertura e critérios

Leitura estática desta frente; resultados executados pertencem aos logs do integrador/agente de bancos. Race/vet e cobertura global não demonstram cada combinação abaixo.

| Área | Cobertura observada no código de testes / evidência existente | Lacuna útil / aceite proposto |
|---|---|---|
| JS/TS/pnpm | fixture TypeScript pnpm, scripts; sample npm real menu/supervisor | npm/yarn/packageManager versionado, conflitos locks com override, Bun, nomes normalizados colidentes; sem executar código no scan |
| Python/uv | fixture pyproject scripts uv; Django e bancos em integração própria | poetry/pip/requirements/Pipfile incoerente, uv+poetry conflito, entrypoint inexistente; não inferir módulo web |
| Go | fixture módulo main único; CLI real Go compilada | zero/múltiplas entradas, root main ./., módulos aninhados, go.work, build tags; nenhuma execução no scan |
| Monorepo | fixture três linguagens em pastas, exclusões/symlink | workspace JS root/filhos, manager herdado explícito, scripts fora de serviços, IDs semelhantes; ordenação determinística e escopo da sugestão |
| Infra | planos Compose sem --volumes; contexto k8s bloqueado; local/Compose/Dockerfile/kind/Helm reais na rodada | Podman runtime ausente; kustomize operacional específico; persistência PVC exige escrita/leitura e UID/recurso, não só Pending |
| Banco | quatro ferramentas planos SQLite; scaffolds PostgreSQL; nome migration/env/reset; runner oito combinações separado | atualizar status somente após nova matriz; dados seed idempotente por ferramenta, segredo ausente no dry-run, falha de migration sem alterar origem |
| Geração/menu | idempotência/hash, preflight inválido, symlinks, seleção EOF, overrides, argv/exit, menu/rich/plain | repetir contratos para cada adapter novo, evidência originada no manifest, comparação hashes em diretório com espaço |

Contrato mínimo por novo adapter: manifest válido e inválido; ausência/ambiguidade; nenhum subprocesso/rede no scan; root/subdiretório/monorepo; exclusion/symlink; comandos argv previsíveis; start não inventado; nomes válidos/colisões explícitas; init seleção e generate twice sem diffs; launcher real mínimo sem instalar deps; overrides manuais preservados; failure/exit propagado; README suporte e limites.

Integração runtime requer disponibilidade e autorização; isolar fixtures/cache/portas/containers e não testar reset contra DB do usuário. Plataformas macOS/WSL não disponíveis permanecem limites, compilação não conta como runtime.

## Matriz executada nesta rodada de pacotes

| Requisito/cenário | Fixture / teste | Nível | Estado / limite |
|---|---|---|---|
| JS/TS npm/pnpm/yarn e packageManager | test/internal/mudarro/support_matrix_test.go | Detecção/contrato | Passou; start/test/script selection, override e locks |
| Python pip/uv/poetry e ambiguidades | test/internal/mudarro/support_matrix_test.go | Detecção/contrato | Passou; nenhum módulo start inventado |
| Go workspace/módulos/múltiplos mains | test/internal/mudarro/support_matrix_test.go | Detecção/contrato | Passou; build tags continuam limitação |
| Todas7 combinações manager/família | test/internal/mudarro/support_generation_test.go | Main scan/init/generate/menu/dry-run | Passou; hashes/2ªgeração/argv, sem executar gerenciadores |
| Docker/Podman Compose, Dockerfile, k8s3modos | test/internal/mudarro/infrastructure_matrix_test.go | Plano/executor simulado | Passou; provider/contexto/labels/workloads/PVC/DaemonSet/JSON inválido |
| Manifest inválido/grande/symlink/exclusões | test/internal/mudarro/core_test.go + test/internal/mudarro/support_matrix_test.go + test/internal/mudarro/generation_regression_test.go | Negativos/fixtures | Passou; irmão válido preservado; preflight sem escrita parcial |
| Script normalizado colidente | test/internal/mudarro/script_collision_test.go | Antes/depois + init | Falhou antes; passou após fix; seleção ambígua não salva config |
| Node/npm real | scripts/validate-menu.sh | E2E | Passou; wrappers/menu/supervisor/idempotência/proteção |
| Python/Go real | scripts/integration-languages.sh | E2E | Passou; explícito start Python, wrappers test e supervisor; sem instalar dependências |
| Bancos8combinações | scripts/validate-databases.sh modo existing deps | Runtime real | Passou; migrations/seed/repetição. Hook não equivale a inserção Goose/Prisma |
| Docker Compose/Dockerfile | integration-containers.sh/integration-dockerfile.sh | Runtime real | Passou; recursos próprios/ownership |
| Kustomize/Helm | validate-kind.sh modo kustomize/integration-helm.sh | Runtime real | Passou, retorno0; PVC montado Bound/UID/conteúdo |
| Yarn/uv/poetry/Podman runtime | ferramentas host | Não executado | Ausentes; contratos acima não equivalem a runtime. Nenhuma instalação autorizada |
| macOS/WSL runtime | plataforma externa | Não executado | Darwin cross-build apenas; WSL não disponível |
| Windows nativo (MUD-068) | testes `*_windows_test.go` | Cross-build/vet local; runtime no CI `windows-latest` | Ver "Windows nativo — MUD-068"; aceite físico em console pendente |

Logs packages-* em samples/menu-auto/evidence. Referência anterior 62,8% era do pacote raiz antes da divisão; usar perfil agregado novo. Não instalar gerenciadores nem executar db-install/install automaticamente. Novas linguagens propostas em [roadmap](language-roadmap.md), sem adapters novos implementados.

Pnpm runtime existente também passou (`packages-pnpm-real.txt`), sem install. Cobertura final agregada 72,0% com `-coverpkg=./internal/mudarro/...`; cobertura de cada teste individual não representa o total do produto. Integrações banco/infra usaram snapshot da extração `b21c91f...`; depois, correção só de colisão no scan JS e literais nomeados model foi validada por suíte final e E2E Node/Python/Go/pnpm.

## Gate final dos testes — 2026-10-06

As 41 funções Test passaram com race e instrumentação explícita coverpkg após migração. Cobertura deduplicada: **1.035/1.427 statements = 72,53% (Go exibe 72,5%)**; anterior 1.028/1.427=72,04% permanece histórica. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e diff check passaram. Compilação Darwin não representa runtime macOS. Árvore de testes contém 16 arquivos: nove testes da orquestração e um helper; dois testes terminal e quatro helpers PTY. Sem hooks de teste ou overlays em produção.

Evidências: [test-tree-tests.txt](../samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](../samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](../samples/menu-auto/evidence/test-tree-gates.txt).

## Validação histórica do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).


## MUD-037 — checkpoint local

Tarefas explícitas selecionáveis Turbo/Nx e ownership da raiz implementados localmente. Build/test/ordenação/cache local/repetição/mudança de entrada/falhas e wrappers reais passaram com Turbo2.11.7/Nx23.2.1 e samples internos.186racePASS/81,38%, nenhum teste pulado. JSON estrito; scan não infere grafos/plugins/defaults. macOS/WSL/Podman nativos e aceite físico continuam pendentes. Novo lote não publicado. [Evidências e limites](../samples/menu-auto/evidence/orchestration/README.md).
