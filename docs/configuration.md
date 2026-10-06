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
name: Meu serviço
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
