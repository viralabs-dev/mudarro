# MUD-036 — análise somente leitura

Nota e Kanban consultados: `/vault/projetos/viralabs/mudarro/06-backlog/MUD-036 — Resolução de workspaces JavaScript (Mudarro).md`. Continua A fazer, sem responsável; depende de MUD-035 (também A fazer). Nenhuma atividade assumida, código/Vault/Git alterado ou instalação realizada.

## Estado observado

`scan.go` já percorre e ordena todos os diretórios não excluídos, ignora symlinks/node_modules e chama JavaScript.Detect uma vez por package.json. Logo root e filhos já aparecem como serviços por diretório: workspaces não precisa criar uma segunda lista de serviços. Hoje `JavaScript.Detect` só consulta evidências locais: packageManager explícito vence lockfile; lock único escolhe manager; ausência escolhe npm. Assim filho sem lock/packageManager num workspace pnpm/yarn recebe npm indevidamente. Não lê workspaces nem pnpm-workspace.yaml, não relaciona ownership, não herda manager. Nome/ID de serviço vêm do diretório; package.json name só define nome do projeto na raiz.

Contratos existentes: `TestSupportJavascriptManagers` exige defaults standalone npm, explicit manager vence lock, argv `manager run original-script`, scripts custom não selecionados automaticamente. `TestScript...` mantém colisões sanitizadas explícitas. Runner usa cwd do Service.Dir. MUD-036 deve preservar esses contratos para pacotes standalone.

## Escopo implementável recomendado

1. Inventariar declarações e memberships em passo determinístico no Scan, usando somente manifests guardados/limitados existentes. Resolver package.json workspaces array e forma object.packages (compatibilidade Yarn), pnpm-workspace.yaml packages via yaml.v3 já disponível. Nada de invocar npm/pnpm/Corepack durante scan. Não adicionar biblioteca nova.
2. Resolver membership contra diretórios package.json já encontrados, canonicalizados para slash; raízes/filhos/nested owners separados. Nearest declared owner que efetivamente inclua o filho; pacote fora dos padrões continua standalone. Excludes e diretórios ignorados têm precedência. Dedupe por diretório, nunca por nome de pacote. Múltiplas definições divergentes/nomes duplicados devem gerar evidence+warning/pending, sem escolher ownership silenciosamente.
3. Herança manager apenas para membro confirmado, a partir de evidência válida/unívoca do owner. Standalone sem evidência continua npm. Não herdar comandos/framework/deps/infra/banco da raiz. Regerar comandos locais com o manager resolvido antes de produzir sugestões; trocar somente Service.Manager depois de Detect deixaria argv npm antigo, bug provável.
4. Para npm/yarn usar package.json workspaces. Para pnpm, YAML é autoridade; package.json.workspaces sozinho não prova membership pnpm. Se ambos existem, não fazer união automática: escolha fonte conforme manager, registre conflito quando divergente. pnpm YAML packages ausente não significa todos os subpacotes.
5. Globs: suportar explicitamente literais, *, **, ? e classes, negativos ! para YAML; validar/normalizar `./`, slashes e negar escape fora do checkout. Não prometer semântica completa minimatch/extglob/brace sem implementação/testes. Padrão não suportado gera aviso e impede herança afetada, em vez de resolver aproximação. ** exige matcher próprio: filepath.Match não atravessa separadores. Recusa traversal é limite deliberado, mesmo que um manager aceite caminhos externos.
6. Ownership da instalação: raiz mantém seu comando install; membros herdados não devem ganhar instalações duplicadas implicitamente. Recomendação inicial: ausência de install no filho herdado + pending acionável indicando owner/install; execução e script run ficam no cwd do filho. Evitar introduzir novo Command.Dir nesta rodada. Se install local explicitamente configurado pelo usuário, preservar.

## Precedência/conflitos proposta conservadora

- Sem owner: contrato standalone existente.
- Membro + root manager válido + sem evidência local: herdar.
- Evidência local concordante: manter mesmo manager e registrar origem.
- Root manager ambíguo/desconhecido: não herdar nem cair silenciosamente em npm; ações automáticas dependentes bloqueadas, configuração explícita resolve.
- Root e child manager/lock divergentes: reportar conflito; não escolher automaticamente instalação/manager. Isso não impede o usuário de declarar manualmente manager/comandos no mudarro.yaml/json.
- Workspace aninhado: owner mais próximo que inclua o pacote; declarar explicitamente boundary, sem atravessar root da varredura.

Essas políticas são escolhas rotineiras conservadoras, não precisam interromper a fila para confirmação. Uma decisão humana só é necessária se quiser suportar configuração conflitante como comportamento automático (por exemplo root pnpm + child npm) ou exigir globs completos/instalação automática global, que expande escopo e altera semântica de ownership.

## Matriz e runtime viáveis sem instalação

Fixtures: npm root+filhos; pnpm root YAML e lock; yarn manifests sem runtime; standalone adjacente; nested owner; zero scripts na raiz; overlap/duplicate patterns; negativos; **; inválido/oversized/FIFO; excludes; symlink externo; pacote sem name; scoped/duplicate names; root manager ambíguo; child explicit manager conflict; scripts com colisão e nomes colon/dot; hashes antes/depois scan/init/generate repeated; wrapper cwd/argv com espaços.

Runtime disponível nesta sessão: npm 11.19.0 e pnpm 11.3.0 (`COREPACK_ENABLE_NETWORK=0 pnpm --version`); yarn ausente PATH. Testes npm/pnpm reais de script podem usar Node builtin e zero dependências, sem install, registry ou Corepack downloads. Executar script que registra cwd/args/marker para comprovar owner/cwd e script certo; preservar hashes package/locks. Yarn permanece contrato simulado/dry-run, não runtime aprovado. Não fazer `install` para provar ownership: teste plano/wrapper que aponta só owner, ou fake executor claramente identificado.

## Ordem sugerida

MUD-035 valida manager primeiro; MUD-036 parser/membership+herança+testes; runtime npm/pnpm sem dependências; docs EN/pt-BR e Vault reconciliados. Turborepo/Nx/orquestração/retries/cache/topologia continuam MUD-037, fora deste escopo.

## Fontes primárias verificadas

- npm declara workspaces no package.json e comandos no contexto de workspace: https://docs.npmjs.com/cli/v11/using-npm/workspaces/
- Yarn declara array de globs no package.json: https://yarnpkg.com/features/workspaces
- pnpm YAML define raiz/include/exclude; raiz incluída; sem packages só raiz; package.json workspaces não substitui YAML: https://pnpm.io/settings#packages

As fontes atuais podem descrever pnpm mais novo que o 11.3.0 instalado: não atribuir todas as regras recentes ao runtime local sem teste. Propostas de precedência/ownership acima são políticas do Mudarro, não afirmação de compatibilidade integral com managers.
