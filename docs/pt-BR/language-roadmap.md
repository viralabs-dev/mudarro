[English](../language-roadmap.md) · Português (Brasil)

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


## Inventário anterior — antes deste lote estático

Locks Bun e runtime genuíno Bun1.4.0 validados localmente. Detecção Pipfile/Pipenv, scripts e argv de banco implementados/testados; runtime Pipenv2026.8.0 passou fixture isolada sem dependências. Nomes packageManager conhecidos validados; versões informadas agora exigem SemVer exato (metadata somente sintática). Entradas Go do host e IO limitado corrigidos; seleção targets e isolamento explícito passaram; parser completo de membros workspace continua aberto. Yarn1.22.22/4.18.1, uv0.12.23 e Poetry2.5.1 passaram runtime isolado sem dependências. O roadmap abaixo fica restrito ao trabalho restante; nenhuma demanda/adapter de linguagem nova inferida. [Evidência mais recente](../samples/menu-auto/evidence/isolated-runtimes/README.md). Mudanças locais não commitadas sobre base publicada b887f64.


## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Proposta anterior de prioridades (superada pelos contratos estáticos aprovados)

Não há demanda comprovada de linguagem nova no material consultado. Ordem abaixo é hipótese de valor/custo a validar com repositórios reais de Daniel; não ranking de popularidade. Nenhuma recomendação externa ou benchmark usado.

1. **Fechar managers/monorepos das três famílias atuais.** Maior retorno provável: provisionamento restante de versões packageManager; integrações reais de dependências/frameworks; grafo/Turbo/Nx JS restante (ownership/herança manager declarado entregues); integrações de grafo/framework Go (seleção de entrada e parser oficial de membros go.work entregues). Custo baixo/médio, depende de definir semântica de herança e escolha antes de código. Aceite: fixtures realistas, ambiguidade explícita, nenhum comando inventado, cobertura de geração/menu.
2. **Rust/Cargo, se houver projeto real alvo.** Manifest TOML já tem parser no produto; build/test/install e seleção de binários oferecem contrato plausível. Custo médio: workspace members/default-members, biblioteca vs múltiplos bins e caminho/target precisam regra explícita. Não afirmar suporte até adapter, fixtures e execução cargo real autorizada. Evitar start automático baseado somente Cargo.toml.
3. **PHP/Composer, se houver demanda web pessoal.** JSON facilita manifest/scripts; custom scripts podem reutilizar mecanismo existente. Custo médio: Composer runtime e PHP, frameworks/start variados. Começar scripts explícitos; Laravel/Symfony somente com evidência e critérios próprios.
4. **Java/Kotlin ou .NET somente com demanda concreta.** Valor alto em equipes desses ecossistemas, demanda aqui desconhecida; custo médio/alto por Maven/Gradle multimodule/wrapper, JVM targets ou csproj/solution múltiplos projetos. Separar build/test de start, evitar inferir entrypoint e executar wrappers na detecção.
5. **Ruby/Elixir/C/C++ depois de alvo confirmado.** Managers/frameworks/Make/CMake/Meson e múltiplos targets ampliam ambiguidades. Comandos custom já atendem parte do uso sem adapter automático; medir fricção real antes de adicionar.

Critério para priorizar: um repo real autorizado + ação repetitiva hoje manual + evidência estática suficiente + teste offline de detecção + E2E mínimo seguro. Registrar demanda/valor/complexidade/dependências no Kanban; um integrador escreve compartilhado, revisões e fixtures em snapshots separados. Não adicionar um grande catálogo nominal que só classifica arquivos sem gerar comandos confiáveis.

Estados sugeridos: inventariado, demanda confirmada, contrato definido, adapter implementado, geração testada, runtime exercitado, limites documentados. Nunca usar 'suportado' para apenas identificar manifest ou compilar o Mudarro.

## Registro atual

MUD-015 planeja; MUD-016 amplia testes atuais. Colisões de scripts JS foram corrigidas com regressão; as demais lacunas requerem contrato explícito e priorização. Não há suporte automático novo Rust/PHP/Java/.NET nesta rodada. Critérios e fixtures em [matriz](support-matrix.md).


Parser Go workspace, seleção explícita de entrada e política herdado/off entregues e validados localmente; resolução de replacements/toolchain/grafo e herança de serviços workspace permanecem fora da implementação. [Evidências](../samples/menu-auto/evidence/go-work-parser/README.md).


Ownership/herança somente manager e instalação raiz dos workspaces JavaScript declarados entregues localmente; grafo/Turbo/Nx e equivalência completa da linguagem de globs permanecem separados. [Evidências](../samples/menu-auto/evidence/javascript-workspaces/README.md).


## MUD-037 — checkpoint local

Tarefas explícitas selecionáveis Turbo/Nx e ownership da raiz implementados localmente. Build/test/ordenação/cache local/repetição/mudança de entrada/falhas e wrappers reais passaram com Turbo2.11.7/Nx23.2.1 e samples internos.186racePASS/81,38%, nenhum teste pulado. JSON estrito; scan não infere grafos/plugins/defaults. macOS/WSL/Podman nativos e aceite físico continuam pendentes. Novo lote não publicado. [Evidências e limites](../samples/menu-auto/evidence/orchestration/README.md).
