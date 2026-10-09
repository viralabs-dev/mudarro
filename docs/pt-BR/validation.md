[English](../validation.md) · Português (Brasil)

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


# Validação

Estado inicial medido em 2026-10-05. Os resultados abaixo distinguem execução real, testes automatizados e plataformas ainda não exercitadas.

| Verificação | Evidência / estado |
|---|---|
| Testes unitários e concorrência | `go test -race ./...` passou em Linux; detecção, ambiguidade, exclusões, YAML inválido, argumentos, geração, edição manual, menu e confirmação destrutiva |
| Análise estática | `go vet ./...` e sintaxe Bash |
| Processo local | Smoke real: up repetido, status, logs, restart e down, em diretório com espaço |
| Docker Compose | Stack temporária Alpine: up, status, logs, restart e down passaram |
| Kubernetes | Cluster kind exclusivo: up, status, logs, restart e down passaram; objeto PVC existente após down; fixture inicial sem prova de dados; cluster removido |
| PostgreSQL e SQLite | As oito combinações com Goose, Alembic, Django e Prisma passaram em ambientes temporários |
| Seeds | Django/Alembic: inserção e repetição idempotente; Goose/Prisma: chamada do hook de seed definido pelo projeto |
| Instalador | Download simulado, instalação, repetição/atualização, checksum corrompido e remoção passaram |
| Sem rede | Script dedicado usa Docker `--network none` para scan/init/generate |
| Podman, macOS | Jobs próprios na CI; resultado remoto deve ser consultado antes de afirmar validação |
| WSL | Usa binário Linux; smoke real em WSL ainda não executado |
| PostgreSQL nativo | Ainda não executado: initdb/pg_ctl ausentes no host |
| Helm, Dockerfile isolado | Passaram na rodada ampliada, scripts próprios |

Ferramentas exercitadas localmente: Go 1.27.1, Node 26.7.0, Python 3.14.7, Prisma 7.10.0, Alembic 1.20.0, Django 6.1.1, Goose 3.24.1 e PostgreSQL 17. A CI utiliza versões declaradas nos workflows; não se presume equivalência com esta máquina.

## Reprodução

```bash
go test -race ./...
go vet ./...
go build -o bin/mudarro ./cmd/mudarro
bash scripts/smoke.sh "$PWD/bin/mudarro"
bash scripts/test-installer.sh
bash scripts/integration-containers.sh "$PWD/bin/mudarro" docker
bash scripts/integration-containers.sh "$PWD/bin/mudarro" podman
bash scripts/integration-offline.sh "$PWD/bin/mudarro"
bash scripts/validate-kind.sh "$PWD/bin/mudarro"
go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
bash scripts/validate-databases.sh "$PWD/bin/mudarro"
```

Os testes de integração baixam imagens/pacotes e criam ambientes temporários. Os scripts `validate-*` limpam apenas os recursos exclusivos que criam. `integration-k8s.sh` exige um contexto informado e usa um namespace próprio.

## Distribuição

O repositório é público. A release [v0.1.0](https://github.com/viralabs-dev/mudarro/releases/tag/v0.1.0) foi publicada em 2026-10-06, com binários Linux/macOS para amd64/arm64 e `checksums.txt`. A instalação pública real em diretório temporário verificou o checksum e retornou `v0.1.0`. A release [v0.2.0](https://github.com/viralabs-dev/mudarro/releases/tag/v0.2.0) (2026-10-09, tag em `41c29c3`, CI verde em Linux/macOS/Windows) acrescentou os zips Windows amd64/arm64; a instalação pública foi verificada em Linux, macOS 14 (arm64), macOS 15 (Intel) e Windows com PowerShell 7 e 5.1, todos retornando `v0.2.0`. Consulte a [página de releases](https://github.com/viralabs-dev/mudarro/releases) e a CI para versões futuras.

## Limites conhecidos

- A detecção é baseada em evidências estáticas; não compreende regras de negócio nem garante identificar todo framework possível.
- Windows nativo desde a v0.2.0, com os limites da [matriz de suporte](support-matrix.md) (comandos `shell:` e scripts `bash` continuam exigindo `bash` no PATH, como o do Git for Windows).
- Migrations/seed em containers exigem comandos explícitos quando os runtimes não estão no host.
- O Prisma 7 testado exige inicializar o arquivo SQLite ausente antes da criação da migration; `db-init` cobre esse caso.
- O `doctor` verifica executáveis e Compose, não é um health check da aplicação.
- Os templates de banco não criam modelos de domínio. Ajustes em configurações existentes permanecem explícitos.

## Rodada local de 2026-10-06: menu automático

Checkout pessoal `main @ 049a405`; nova cobertura e [sample completo](samples/menu-auto/README.md) com configuração/menu realmente gerados, gravação PTY e comandos de reprodução. Go/race passou (62,8% statements no pacote interno), vet, instalador simulado, launchers/submenu, processo local, Docker Compose, geração sem rede e Kubernetes temporário passaram. Na fixture inicial, o objeto PVC permaneceu existente, Pending, sem prova de dados. Na fixture ampliada montada, Bound/UID/conteúdo após down/up/restart foram comprovados em `kubernetes-persistence.txt`.

A matriz de bancos foi repetida com dependências existentes e passou 8/8; logs no sample. PostgreSQL nativo, Podman real, macOS real e WSL não executados por ferramentas/plataformas indisponíveis. Dockerfile isolado e Helm agora passaram com scripts próprios. GIFs reais renderizados e frame PNG foram validados; screenshot gráfica permanece pendente de captura na sessão já confirmada na área 4 da tela do notebook. Consulte a matriz no sample para as fronteiras exatas. O CI remoto de referência foi confirmado success em `f62ce9d`; as novas alterações permanecem locais.

Visual v2 final: `go-tests-v2.txt`, `go-vet-v2.txt`, `menu-v2.txt` e `media-v2-validation.txt` registram race/coverage, vet, CLI real e decodificação da mídia. Título longo estreito falhou inicialmente e foi corrigido antes do checkpoint.

Rodada ampliada manifests e Helm: volume montado Bound, mesmo UID e marcador lido após down/up/restart (`kubernetes-persistence.txt`, `helm-persistence.txt`). Isso comprova persistência durante o ciclo da fixture, não backup/disaster recovery nem retenção após excluir o cluster. Revisões paralelas somente leitura integradas, recursos próprios removidos, sem instalações adicionais.

## Pacotes e expansão de testes — 2026-10-06

Executor, model/adapters, projectfs e terminal são pacotes reais; interfaces e orquestração permanecem na raiz. Gates por etapa em `packages-executor.txt`, `packages-adapters.txt`, `packages-terminal.txt`. Gate final race/vet/build Linux/Darwin/E2E passou: **72,0% agregados** dos pacotes internos via coverpkg/perfil. Não comparar diretamente ao antigo62,8% do pacote raiz. Revisão independente confirmou grafo acíclico/guardas preservadas.

Testes adicionais cobrem managers7famílias, workspace/ambiguidades, banco/infra/k8s3modos, scan/init/generate/menu/argv/hashes/symlink/manifests. Colisão de nomes JS demonstrou falha antes e aprovação após correção. Runtime real Node/npm, pnpm existente, Python e Go; Compose/Dockerfile, kustomize/Helm e bancos8/8 revalidados. Sem yarn/uv/poetry/Podman real ou WSL/macOS real. [Matriz detalhada](support-matrix.md) separa contratos de execução e [roadmap](language-roadmap.md) planeja novos adapters, não os implementa.

Snapshot final v4: `/tmp/mudarro-checkpoint-v4/bin/mudarro`, SHA256 `76ccc665970c48a7f7c95f9a9052599cd23aee4580060b11a265a5b1bd91d942`. Nenhum commit/push/publicação. GIF final em `samples/menu-auto/evidence/packages-visual/`; PNG gráfica pendente; aceite visual recebido.

[Como interpretar cobertura, totais e lacunas](coverage.md):1.028/1.427 statements =72,04%; E2Es externos ficam separados do perfil.

## Gate final dos testes — 2026-10-06

As 41 funções Test passaram com race e instrumentação explícita coverpkg após migração. Cobertura deduplicada: **1.035/1.427 statements = 72,53% (Go exibe 72,5%)**; anterior 1.028/1.427=72,04% permanece histórica. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e diff check passaram. Compilação Darwin não representa runtime macOS. Árvore de testes contém 16 arquivos: nove testes da orquestração e um helper; dois testes terminal e quatro helpers PTY. Sem hooks de teste ou overlays em produção.

Evidências: [test-tree-tests.txt](../samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](../samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](../samples/menu-auto/evidence/test-tree-gates.txt).

go test sem instrumentação também passou: [test-tree-plain-tests.txt](../samples/menu-auto/evidence/test-tree-plain-tests.txt).

## Validação histórica do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
