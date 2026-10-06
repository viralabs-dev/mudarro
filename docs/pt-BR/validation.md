[English](../validation.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

## Gate atual da UI configurável

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

O repositório é público. A release [v0.1.0](https://github.com/viralabs-dev/mudarro/releases/tag/v0.1.0) foi publicada em 2026-10-06, com binários Linux/macOS para amd64/arm64 e `checksums.txt`. A instalação pública real em diretório temporário verificou o checksum e retornou `v0.1.0`. Consulte a [página de releases](https://github.com/viralabs-dev/mudarro/releases) e a CI para versões futuras.

## Limites conhecidos

- A detecção é baseada em evidências estáticas; não compreende regras de negócio nem garante identificar todo framework possível.
- Não há shell interativo nativo do Windows; usar WSL.
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

## Validação atual do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
