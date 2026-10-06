[English](../documentation.md) · Português (Brasil)

## Correção atual da CI — MUD-030

Verificada [CI de push37416072569](https://github.com/viralabs-dev/mudarro/actions/runs/37416072569), head `bade86e`: concluída com **falha**. Jobs Ubuntu/bancos/containers passaram; macOS executou e falhou no tratamento EOF/cleanup com race, incluindo panic com índice n=-1. macOS não é mais globalmente “não executado”: cross-compile Darwin local passou, runtime real da CI falhou. Podman local continua indisponível; CI containers passou, distinta da execução local.

Leitura autenticada de permissões Actions retornou enabled=true/allowed_actions=all; não identifica quem alterou settings. Nenhuma alteração de settings/workflow/config, habilitação, dispatch ou rerun ocorreu. Este commit corretivo usa o marcador oficial [skip ci] para o push autorizado somente de main, respeitando o pedido de não executar Actions sem mudar settings. Gates/evidências brutos anteriores preservados como históricos, superados quanto ao estado atual da CI. Linux local74 testes/76,38% permanece válido para checkpoint registrado. Fix de portabilidade passou 75 testes locais Linux com race, vet e compilação Linux/Darwin arm64. Runtime macOS corrigido permanece não validado; pular CI não significa aprovação.

# Idiomas e manutenção da documentação

Inglês é padrão do repositório: README.md e docs/**/*.md. Português brasileiro fica em docs/pt-BR; README.md nesta pasta é a entrada traduzida. O Vault permanece somente em português e tem manutenção separada.

Assets, samples gerados e evidências mantêm uma única cópia fora de docs/pt-BR. Documentos traduzidos referenciam essa cópia. Não traduzir comandos, flags, chaves/grupos de configuração, token de confirmação destrutiva, logs, IDs, hashes, caminhos de fixtures, licenças ou valores fornecidos pelo usuário. Exemplos originais preservam comentários em português deliberadamente.

Mudanças de comportamento atualizam inglês e pt-BR na mesma alteração. Conciliar tabelas, estados, limites e links sem transportar resultados antigos. Ambas versões apontam uma à outra. Nova localização de testes/build exige novos resultados reais; tradução não transforma contrato em runtime nem frame composto em screenshot.

Se tradução ficar pendente, marcar no início `Tradução pendente: <origem, seção alterada e data>` e registrar ação exata na atividade. Revisão confere números/comandos/restrições, navegação e links relativos. Não exige serviço de tradução nem publicação.

A localização inicial preserva documentos portugueses completos e fornece orientação operacional equivalente em inglês. Narrativas históricas podem ter mais detalhes em português; evidências e limites de execução compartilhados são referência. Buscar equivalência semântica, não contagem idêntica de frases.

- [Plano aprovado de configuração UI](ui-configuration-plan.md)

Implementação atual: [configuração UI](ui-configuration.md). Config/i18n/tema/preview implementados; fonte/aceite mouse em emulador real continuam pendentes separados. Gates finais:74 testes,76,38% statements agregados; ver [cobertura](coverage.md).
