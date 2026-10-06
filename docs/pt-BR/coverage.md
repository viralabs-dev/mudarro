[English](../coverage.md) · Português (Brasil)

A base publicada agora é `b887f64`, após reescrita autorizada da autoria pessoal de `56dadba`, com árvore idêntica. [Mapa de SHAs e provas de preservação](../samples/menu-auto/evidence/git-identity-repair/README.md). As novas mudanças de robustez/Unicode e esta referência permanecem locais.

## Checkpoint integrado mais recente — parser oficial Go workspace

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


# Como interpretar a cobertura medida

Medição atual após migração de testes, 2026-10-06: **1.035/1.427 statements executados = 72,53%**, Go exibe 72,5%;392 statements não executados. Os 41 testes passaram.

Medição histórica anterior à migração de testes, 2026-10-06: **1.028/1.427 statements executados = 72,04%**, exibidos pelo Go como 72,0%;399 statements não executados no perfil. O checkpoint anterior71,8% foi atualizado após teste/fix de colisões JS. Não é percentual de linhas, branches, requisitos ou plataformas validadas.

```bash
GOCACHE=/tmp/mudarro-go-cache /home/danielsouza/sdk/go1.27.1/bin/go test -race \
  -coverpkg=./internal/mudarro/... \
  -coverprofile=docs/samples/menu-auto/evidence/unicode-resize/coverage.out -v ./...
/home/danielsouza/sdk/go1.27.1/bin/go tool cover \
  -func=docs/samples/menu-auto/evidence/unicode-resize/coverage.out
```

Escopo instrumentado: pacotes internos sob internal/mudarro. O comando executa testes de todos os pacotes; cmd/mudarro não entra no coverpkg. Model contém só tipos, sem statements executáveis. Perfil atomic agrega contadores dos blocos; totais por pacote calculados somando NumStmt uma vez por bloco, coberto quando Count>0. Não se faz média simples dos percentuais de pacotes/test binaries.

Totais históricos da etapa MUD-017:

| Pacote | Cobertos / total | Percentual |
|---|---:|---:|
| `/` | 589/935 | 62.99% |
| `/adapters` | 279/304 | 91.78% |
| `/executor` | 10/13 | 76.92% |
| `/projectfs` | 19/23 | 82.61% |
| `/terminal` | 138/152 | 90.79% |
| `model` | Sem statements executáveis | Não aplicável |

## O que entra e o que fica separado

Entram os testes Go de detecção, contratos, simulações do Executor, geração/preservação/configuração, CLI Main chamada no processo e terminal. Race detecta corridas nas execuções exercitadas; não mede todos os interleavings.

Scripts E2E executam binários comuns, sem instrumentação coverage: Node/npm/pnpm/Python/Go, supervisor real, Docker/Compose/kind/Helm e bancos não incrementam esse perfil. São evidências independentes. OSExecutor Run executado em teste contabiliza seu código pai, não cobertura do subprocesso/aplicação. Portanto0% local_runner no perfil NÃO significa ausência de execução real do supervisor; seus E2Es passaram. Da mesma forma, runtime aprovado não transforma automaticamente branch unitário em coberto.

A medição antiga 62,8% era do pacote raiz anterior à extração. A medição atual reúne raiz+adapters+executor+projectfs+terminal e testes novos; escopo/denominador mudaram. Não interpretar diferença como ganho equivalente nem somar a cobertura de cada test binary.

## Lacunas e prioridades por risco

1. **Supervisor/identidade/sinais**: running/readState/local/supervise0% neste perfil, embora E2E real passou. Priorizar estado stale/corrompido, recusa de PID/token/projeto estranho, falha de startup, árvore de processos e concorrência. São riscos de controle de processo/preservação; não apenas percentual.
2. **Doctor e dependências**: doctor0%; exercitado E2E somente em ambientes disponíveis. Acrescentar negativos executável ausente/não executável, configuração bloqueada, symlink e retorno correto.
3. **Erros de IO/execução**: atomicWrite53,3%, OSExecutor.Output40%, SafePath85,7%/ReadManifest77,8%. Priorizar falha de criação/rename/leitura, stdout+erro e manutenção de arquivos/exitstatus; sem injetar abstrações apenas para cobrir linhas.
4. **TTY e apresentação**: terminalSize final81,8%, Choose97,3%, HasColor100%; PTYs reais isolados exercitam ramo rico. Redimensionamento/Unicodewide/cancelamento ainda merecem aceites próprios. Helpers legados ascii/terminalColumns sem consumidor permanecem0%. Não equivale a risco operacional proporcional.
5. **Templates e combinações**: djangoDatabase33,3%, scaffoldFiles74,8%, Config.Validate71,0%. Verificar variantes não exercitadas e configurações inválidas, preservando defaults e modelos não inventados.

Não há meta100% automática. Priorizar comportamento, falhas e ambientes relevantes; cobertura de statements não comprova segurança, ausência de bugs ou execução de todos os gerenciadores. Yarn/uv/poetry/Podman e macOS/WSL runtime continuam limites explícitos.

Evidências: [perfil](../samples/menu-auto/evidence/test-tree-coverage.out), [funções](../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [totais](../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [matriz real](support-matrix.md).

## Gate final dos testes — 2026-10-06

As 41 funções Test passaram com race e instrumentação explícita coverpkg após migração. Cobertura deduplicada: **1.035/1.427 statements = 72,53% (Go exibe 72,5%)**; anterior 1.028/1.427=72,04% permanece histórica. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e diff check passaram. Compilação Darwin não representa runtime macOS. Árvore de testes contém 16 arquivos: nove testes da orquestração e um helper; dois testes terminal e quatro helpers PTY. Sem hooks de teste ou overlays em produção.

Evidências: [test-tree-tests.txt](../samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](../samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](../samples/menu-auto/evidence/test-tree-gates.txt).

## Validação histórica do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Terminal atual: 164/176 statements (93,18%). Funções: cleanText100%, fitText94,1%, runeColumns80%, wordmark95,7%, HasColor100%, terminalSize81,8%. Limites de largura/grafemas/redimensionamento permanecem.
