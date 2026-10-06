[English](../support-matrix.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

## Gate atual da UI configurável

**74 funções Test passaram com race; 1.552/2.032 statements executados = 76,38% (Go exibe76,4%)**, 480 não executados. Terminal:379/422=89,81%. Vet, build Linux e compilação cruzada do binário inteiro Darwin arm64 passaram; compilação não valida runtime macOS. E2E UI genuíno passou propagação de config personalizada nos wrappers/supervisor, scan traduzido, preview somente leitura/redação de segredos e smoke/menu. Rodadas reais Python/Go e pnpm passaram; Dockerfile e oito combinações de bancos com dependências existentes também passaram. Números anteriores abaixo são históricos com denominadores distintos. Fonte/aceite mouse em emulador real pendentes.

Evidências: [tests-race.txt](../samples/menu-auto/evidence/ui-configuration/tests-race.txt), [coverage.out](../samples/menu-auto/evidence/ui-configuration/coverage.out), [coverage-functions.txt](../samples/menu-auto/evidence/ui-configuration/coverage-functions.txt), [coverage-summary.tsv](../samples/menu-auto/evidence/ui-configuration/coverage-summary.tsv), [vet.txt](../samples/menu-auto/evidence/ui-configuration/vet.txt), [smoke.txt](../samples/menu-auto/evidence/ui-configuration/smoke.txt), [menu-real.txt](../samples/menu-auto/evidence/ui-configuration/menu-real.txt), [ui-real.txt](../samples/menu-auto/evidence/ui-configuration/ui-real.txt).


# Inventário Mudarro — leitura de código, 2026-10-06

Escopo: checkout pessoal, inventário conciliado após extração de pacotes. Nenhuma execução nova, instalação ou adapter implementado nesta frente. Suporte detectado não equivale a execução integral de todo ecossistema.

## Linguagens e gerenciadores

| Família | Evidência de detecção | Ações inferidas | Limites atuais |
|---|---|---|---|
| JavaScript / TypeScript | package.json; TS por tsconfig.json ou dependência typescript | install; scripts classificados start/dev/build aplicação, test/lint qualidade; demais selecionáveis | npm default; npm/pnpm/yarn por lock único; packageManager tem prioridade, mas nome arbitrário não validado; Bun não reconhecido por lock; frameworks next/@nestjs/core/vite/express são rótulos |
| Python | pyproject.toml / requirements.txt / manage.py / Pipfile | pip venv/install, uv sync, poetry install; entrypoints project/poetry sugeridos; Django runserver/test | uv/poetry só por locks; ambos geram pendência; Pipfile identifica Python mas NÃO pipenv: cai em pip; Python não-Django exige start explícito; não detecta Flask/FastAPI/pytest como ações |
| Go | go.mod; parser de package main em arquivos não-test | build ./..., test ./..., mod download; go run da única pasta main | múltiplas/nenhuma entrada bloqueiam inferência start; go.work é evidência, não resolve módulos; build tags não avaliadas como go list |
| Outras / custom | fallback sem manifest reconhecido | comandos declarados; sugestões Makefile/*.sh | configuração manual pode orquestrar qualquer runtime, mas não é adapter automático |

Manager é string de configuração, não registro fechado validado. Entry scripts Python vêm de maps e ordenação final ocorre por serviço/nome. Scripts JS normalizam : e . para -. Colisões agora geram pendência e não oferecem seleção ambígua; regressão demonstrou falha antes e aprovação após correção.

## Monorepos

Walk recursivo detecta serviços por diretório com manifests. Ordenação de diretórios determinística; exclusões, ignorados e symlinks respeitados. JS workspace root e filhos podem virar serviços separados; não há resolução de package.json.workspaces/pnpm-workspace.yaml/yarn workspace/Turborepo/Nx. Não herda manager/infra do root para filho. Go ignora módulos aninhados ao buscar main e detecta cada go.mod separadamente; go.work não produz plano de execução agregado. Sugestões Makefile/shell associadas ao serviço ancestral mais próximo, com fallback índice 0 quando fora de qualquer serviço.

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

Logs packages-* em samples/menu-auto/evidence. Referência anterior 62,8% era do pacote raiz antes da divisão; usar perfil agregado novo. Não instalar gerenciadores nem executar db-install/install automaticamente. Novas linguagens propostas em [roadmap](language-roadmap.md), sem adapters novos implementados.

Pnpm runtime existente também passou (`packages-pnpm-real.txt`), sem install. Cobertura final agregada 72,0% com `-coverpkg=./internal/mudarro/...`; cobertura de cada teste individual não representa o total do produto. Integrações banco/infra usaram snapshot da extração `b21c91f...`; depois, correção só de colisão no scan JS e literais nomeados model foi validada por suíte final e E2E Node/Python/Go/pnpm.

## Gate final dos testes — 2026-10-06

As 41 funções Test passaram com race e instrumentação explícita coverpkg após migração. Cobertura deduplicada: **1.035/1.427 statements = 72,53% (Go exibe 72,5%)**; anterior 1.028/1.427=72,04% permanece histórica. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e diff check passaram. Compilação Darwin não representa runtime macOS. Árvore de testes contém 16 arquivos: nove testes da orquestração e um helper; dois testes terminal e quatro helpers PTY. Sem hooks de teste ou overlays em produção.

Evidências: [test-tree-tests.txt](../samples/menu-auto/evidence/test-tree-tests.txt), [test-tree-coverage.out](../samples/menu-auto/evidence/test-tree-coverage.out), [test-tree-coverage-functions.txt](../samples/menu-auto/evidence/test-tree-coverage-functions.txt), [test-tree-coverage-summary.tsv](../samples/menu-auto/evidence/test-tree-coverage-summary.tsv), [test-tree-gates.txt](../samples/menu-auto/evidence/test-tree-gates.txt).

## Validação atual do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).
