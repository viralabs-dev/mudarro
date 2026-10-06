# Validação

Estado inicial medido em 2026-10-05. Os resultados abaixo distinguem execução real, testes automatizados e plataformas ainda não exercitadas.

| Verificação | Evidência / estado |
|---|---|
| Testes unitários e concorrência | `go test -race ./...` passou em Linux; detecção, ambiguidade, exclusões, YAML inválido, argumentos, geração, edição manual, menu e confirmação destrutiva |
| Análise estática | `go vet ./...` e sintaxe Bash |
| Processo local | Smoke real: up repetido, status, logs, restart e down, em diretório com espaço |
| Docker Compose | Stack temporária Alpine: up, status, logs, restart e down passaram |
| Kubernetes | Cluster kind exclusivo: up, status, logs, restart e down passaram; PVC preservado; cluster removido |
| PostgreSQL e SQLite | As oito combinações com Goose, Alembic, Django e Prisma passaram em ambientes temporários |
| Seeds | Django/Alembic: inserção e repetição idempotente; Goose/Prisma: chamada do hook de seed definido pelo projeto |
| Instalador | Download simulado, instalação, repetição/atualização, checksum corrompido e remoção passaram |
| Sem rede | Script dedicado usa Docker `--network none` para scan/init/generate |
| Podman, macOS | Jobs próprios na CI; resultado remoto deve ser consultado antes de afirmar validação |
| WSL | Usa binário Linux; smoke real em WSL ainda não executado |
| PostgreSQL nativo, Helm, Dockerfile isolado | Implementados; testes locais completos ainda não executados |

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
