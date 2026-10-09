[English](../operations.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Operação

## Fluxo

1. Execute `scan` e examine as evidências; `scan --json` serve para automação.
2. Inicialize a configuração e selecione scripts adicionais.
3. Resolva escolhas de infraestrutura, entrada da aplicação, banco e ferramenta de migrations.
4. Inspecione `generate --dry-run`; gere arquivos e execute `doctor`.
5. Use `run serviço:install`, `venv` ou `db-install` explicitamente quando necessários. Ferramentas de sistema como Docker, Podman e kubectl precisam ser instaladas pelo procedimento do seu sistema.
6. Abra `menu.sh` ou use `run serviço:ação`.

O `doctor` verifica executáveis e o provedor Compose. Não comprova saúde de containers, conexão com bancos, credenciais ou instalação de cada módulo dentro de um runtime.

## Arquivos e regeneração

Arquivos preexistentes são preservados e não passam a pertencer ao gerador. Arquivos já gerados são comparados com o manifesto antes de qualquer escrita. Se um deles foi editado, a geração informa o conflito; preserve sua versão ou restaure o conteúdo gerado antes de continuar. Não remova o manifesto para forçar atualização.

O Mudarro não apaga wrappers antigos automaticamente. Após remover uma ação/serviço da configuração, um wrapper antigo não consegue executar a ação ausente. Faça a limpeza desses arquivos deliberadamente.

## Execução local

O supervisor é iniciado em um grupo de processos próprio e registra um token aleatório. Antes de sinalizar um processo, a ferramenta verifica a identidade do supervisor e o projeto. Não procura processos por nome de aplicação.

Os logs ficam em `.mudarro/run/<serviço>/output.log`. `down` sinaliza o grupo do supervisor; a aplicação deve permanecer em foreground, sem daemonizar ou criar sessões independentes. O start confirma que o processo permaneceu vivo inicialmente, não que a aplicação está saudável.

## Containers e Kubernetes

Compose usa o arquivo explicitamente selecionado; Docker usa `docker compose` e Podman usa `podman-compose`, evitando que o wrapper Podman selecione Docker Compose por acidente. A descida nunca solicita exclusão de volumes. Dockerfile isolado oferece build da imagem e criação/retomada do container. Configure imagem e portas explicitamente. Para stacks compartilhadas no monorepo, mantenha um único serviço proprietário da operação Compose.

Kubernetes sempre recebe contexto e namespace configurados. O caminho de manifests deve conter os recursos desta aplicação. A parada opera sobre Deployments e StatefulSets; não exclui recursos de armazenamento. Helm depende dos labels convencionais de release para selecionar workloads na parada/reinício/logs. Não há provisionamento de cluster nem operação de produção nesta versão.

## Distribuição

`scripts/release.sh vX.Y.Z` compila seis binários sem CGO: Linux/macOS amd64/arm64 em `.tar.gz` e Windows amd64/arm64 em `.zip` (`mudarro.exe`), com `checksums.txt`. O workflow de release é acionado por tags. Linux/macOS usam Bash (`install.sh`); WSL recebe o binário Linux; Windows usa o `install.ps1` (PowerShell 5.1 ou 7+).

Os dois instaladores aceitam `MUDARRO_VERSION`, `MUDARRO_INSTALL_DIR` e `MUDARRO_REPOSITORY`. Não usa `sudo` nem altera o shell do usuário. Downloads vêm de HTTPS e o hash é verificado antes da troca atômica do binário.

## Referências dos templates

- [Configuração Prisma 7](https://www.prisma.io/docs/orm/v7/reference/prisma-config-reference): schema, migrations e URL no arquivo de configuração.
- [Tutorial oficial Alembic](https://alembic.sqlalchemy.org/en/latest/tutorial.html): ambiente e revisões.

As versões e combinações efetivamente exercitadas estão em [Validação](validation.md).
