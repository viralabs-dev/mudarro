[English](../configuration.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Configuração

`mudarro init` cria `mudarro.yaml` sem sobrescrever configuração existente. `--interactive` seleciona sugestões; `--select serviço:script,...` permite seleção não interativa, e `--all-scripts` inclui todas. Scripts padrão de start/dev/build/test/lint detectados entram como operações conhecidas; scripts adicionais são sugestões. O scanner nunca os executa.

## Contrato v1

| Campo | Uso |
|---|---|
| `version` | Obrigatoriamente `1` |
| `name` | Nome exibido no banner; init usa manifest raiz ou nome da pasta |
| `exclude` | Lista de caminhos relativos ou padrões `filepath.Match`; `**` não é recursivo |
| `services[].id` | Identificador único: letras, números, `_`, `-` |
| `dir` | Diretório relativo dentro do projeto |
| `language`, `manager`, `framework` | Tecnologia, gerenciador e framework; escolhas explícitas prevalecem |
| `infrastructure` | `kind`, `mode`, `file`, `context`, `namespace`, `image`, `port`, `generate` |
| `database` | `kind`, `tool`, `url_env`, `path`, `generate` |
| `commands` | Mapa de ações: `args` **ou** `shell`, `group`, `requires`, `destructive` |
| `pending` | Observações produzidas pelo scan; histórico editável, sem bloquear outras operações |

`file` e `database.path` são relativos ao diretório do serviço. Ações declaradas em `commands` substituem ações de mesmo nome dos adaptadores. Nomes de scripts npm com `:` ou `.` são normalizados para `-` no identificador do menu, preservando o nome original nos argumentos executados.

`args` executa diretamente o programa, sem expansão de shell. `${VAR}` referencia uma variável de ambiente obrigatória e `{name}` recebe `--name`, usado na criação de migrations. O Mudarro não carrega `.env` automaticamente. Um `shell` explícito é executado com Bash; não coloque segredos em comandos versionados.

## Compose com PostgreSQL e Prisma

```yaml
version: 1
name: mudarro
services:
  - id: api
    dir: .
    language: typescript
    manager: npm
    infrastructure:
      kind: docker              # podman também é aceito
      mode: compose
      file: compose.yaml
      port: 3000
      generate: true
    database:
      kind: postgresql
      tool: prisma
      url_env: DATABASE_URL
      generate: true
    commands:
      start:
        args: [npm, run, dev]
        group: aplicacao
```

`generate` cria Dockerfile, Compose, exemplos de ambiente e estrutura Prisma quando os destinos não existem. Use `db-install` explicitamente para instalar Prisma 7; revise a estrutura caso o projeto existente use uma versão anterior. A geração nunca migra versões de ORM nem altera automaticamente seus manifests existentes.

Configure `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` e `DATABASE_URL`. Dentro do Compose, o host de banco é `db`; para comandos de migrations executados no host, use `localhost` e `POSTGRES_PORT`. Uma aplicação Python deve expor seu servidor em `0.0.0.0` dentro do container por comando explícito; o Mudarro não altera as flags escolhidas pelo projeto.

As ações de banco executam no diretório local do serviço. Para rodar migrations dentro do container, substitua a ação:

```yaml
commands:
  db-migrate:
    args: [docker, compose, -f, compose.yaml, exec, -T, api, npx, --no-install, prisma, migrate, deploy]
    group: banco
```

## Kubernetes

```yaml
infrastructure:
  kind: kubernetes
  mode: manifests              # ou kustomize; helm integra chart existente
  file: k8s
  context: kind-meu-cluster
  namespace: desenvolvimento
  image: minha-org/minha-app:v1
  port: 3000
  generate: true
```

O cluster e o namespace devem existir; a imagem deve estar disponível no cluster. A geração cria Deployment/Service e, se solicitado, PostgreSQL com PVC e referências a um Secret existente. O arquivo de instruções descreve as chaves necessárias, sem valores. `down` escala os workloads declarados a zero, preservando PVCs. DaemonSets exigem uma ação de parada personalizada.

## Banco local

`kind: sqlite` utiliza o arquivo definido em `database.path` (padrão `app.db`). Execute `db-init` para criar o arquivo, preservando-o se já existir. Para Prisma e Alembic, declare o mesmo caminho na URL específica do ORM. SQLite não exige provisionar um servidor.

PostgreSQL com infraestrutura local e `database.generate: true` oferece `db-init`, `db-up`, `db-status`, `db-create` e `db-down`. Requer binários PostgreSQL instalados e variáveis `POSTGRES_*`. Os dados ficam no serviço em `.mudarro/postgres`; a aplicação continua com ações independentes.

Prisma usa modelos do schema; Django usa models existentes; Alembic cria revisões editáveis e permite ligar `target_metadata`; Goose cria migrations SQL. Os templates não criam entidades de negócio. O fragmento `mudarro_database.py` deve ser importado explicitamente no settings Django existente. Seeds Prisma/Alembic/Goose precisam ser implementados; a fixture inicial Django é vazia.

## Comandos próprios

```yaml
infrastructure:
  kind: custom
commands:
  up:
    args: [make, up]
    group: infraestrutura
  limpar-dados:
    shell: ./scripts/reset-db.sh
    requires: [bash]
    group: banco
    destructive: true
```

Uma ação destrutiva exige digitar `APAGAR <id-do-serviço>`. `db-reset` sempre exige confirmação, mesmo quando sobrescrito.

## Identidade do projeto consumidor — MUD-022

O título do menu usa `Config.Name` do projeto consumidor, não o nome da ferramenta Mudarro. O sample controlado agora usa Aurora (package/config); demo e gravação definem Aurora explicitamente por padrão, com override `MUDARRO_SAMPLE_NAME`. A CLI já usava Config.Name; as entradas de demonstração foram corrigidas. Aceite visual do usuário recebido; somente screenshot gráfica genuína continua pendente.

A sanitização remove controles e caracteres Unicode Cf de formatação. Estimativa conservadora de células considera CJK/fullwidth/emoji duas células e marcas combinantes zero; não garante tratamento perfeito de grafemas/largura em todos os terminais. Glyphs bitmap não suportados ou wordmark longo demais usam fallback textual legível, preservando identidade correta.

A etapa atual tem 43 funções Test e 24 casos nome/modo rich/narrow/plain/NO_COLOR: AtlasAPI, Aurora, Café, 漢字😀, nomes longos e injeção de controles, além de integração CLI com nome consumidor explícito. Race passou; cobertura final e gates de build passaram, registrados abaixo. Medições anteriores permanecem históricas. Mídias genuínas atuais ficam em `evidence/project-name-visual/` no sample canônico; `evidence/packages-visual/` permanece histórico.

## Validação atual do nome consumidor — MUD-022

43 funções Test passaram com race, incluindo 24 casos nome/modo. Vet, build Linux, compilação cruzada dos testes terminal Darwin arm64 e sintaxe Bash passaram. Cobertura: **1.063/1.451 statements = 73,26% (Go exibe 73,3%)**, 388 não executados. MUD-017 anterior 1.035/1.427=72,53% é histórico; mudança de código/denominador impede interpretar diferença como cobertura equivalente de requisitos. Runtime macOS não executado nesta etapa local histórica; CI atual executou e falhou, registrada acima.

GIF final genuíno: 983×739, 8 frames, 15,04 s. Frame composto inspecionado: lettering AURORA e PROJETO / Aurora. É frame da gravação PTY real, não screenshot gráfica ainda pendente. Aceite do desenho recebido. SHA-256 binário: `ef36e6c1142bd7294e543208895cea21e125c6898d14ec168e40c1c7eb327890`.

Evidências: [project-name-tests.txt](../samples/menu-auto/evidence/project-name-tests.txt), [project-name-coverage.out](../samples/menu-auto/evidence/project-name-coverage.out), [project-name-coverage-functions.txt](../samples/menu-auto/evidence/project-name-coverage-functions.txt), [project-name-coverage-summary.tsv](../samples/menu-auto/evidence/project-name-coverage-summary.tsv), [project-name-vet.txt](../samples/menu-auto/evidence/project-name-vet.txt), [project-name-visual/frame-menu-recording.png](../samples/menu-auto/evidence/project-name-visual/frame-menu-recording.png).

Implementação atual: [configuração UI](ui-configuration.md). Config/i18n/tema/preview implementados; fonte/aceite mouse em emulador real continuam pendentes separados. Gates finais:74 testes,76,38% statements agregados; ver [cobertura](coverage.md).


## Mix e Make explícitos

Serviço language:elixir/manager:mix oferece compile/test somente com seleção. mix_umbrella:true é declaração manual; requer mix.exs regular/legível e não infere membros/grafo. C/C++ permanece custom com targetsMake escolhidos pelo usuário; semMake, evidência pendente. [Sample e limites](../samples/menu-auto/evidence/mix-native/README.pt-BR.md).

```yaml
version: 1
name: Explicit Mix umbrella
services:
  - id: umbrella
    dir: .
    language: elixir
    manager: mix
    mix_umbrella: true
    infrastructure:
      kind: local
    commands:
      compile:
        args: [mix, compile]
      test:
        args: [mix, test]
```
