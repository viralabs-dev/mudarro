[English](../language-roadmap.md) · Português (Brasil)

# Proposta de próximas capacidades

Não há demanda comprovada de linguagem nova no material consultado. Ordem abaixo é hipótese de valor/custo a validar com repositórios reais de Daniel; não ranking de popularidade. Nenhuma recomendação externa ou benchmark usado.

1. **Fechar managers/monorepos das três famílias atuais.** Maior retorno provável: Bun/packageManager validado; Pipfile sem falso pip; managers JS root/filhos e workspace; Go múltiplos bins/go.work. Custo baixo/médio, depende de definir semântica de herança e escolha antes de código. Aceite: fixtures realistas, ambiguidade explícita, nenhum comando inventado, cobertura de geração/menu.
2. **Rust/Cargo, se houver projeto real alvo.** Manifest TOML já tem parser no produto; build/test/install e seleção de binários oferecem contrato plausível. Custo médio: workspace members/default-members, biblioteca vs múltiplos bins e caminho/target precisam regra explícita. Não afirmar suporte até adapter, fixtures e execução cargo real autorizada. Evitar start automático baseado somente Cargo.toml.
3. **PHP/Composer, se houver demanda web pessoal.** JSON facilita manifest/scripts; custom scripts podem reutilizar mecanismo existente. Custo médio: Composer runtime e PHP, frameworks/start variados. Começar scripts explícitos; Laravel/Symfony somente com evidência e critérios próprios.
4. **Java/Kotlin ou .NET somente com demanda concreta.** Valor alto em equipes desses ecossistemas, demanda aqui desconhecida; custo médio/alto por Maven/Gradle multimodule/wrapper, JVM targets ou csproj/solution múltiplos projetos. Separar build/test de start, evitar inferir entrypoint e executar wrappers na detecção.
5. **Ruby/Elixir/C/C++ depois de alvo confirmado.** Managers/frameworks/Make/CMake/Meson e múltiplos targets ampliam ambiguidades. Comandos custom já atendem parte do uso sem adapter automático; medir fricção real antes de adicionar.

Critério para priorizar: um repo real autorizado + ação repetitiva hoje manual + evidência estática suficiente + teste offline de detecção + E2E mínimo seguro. Registrar demanda/valor/complexidade/dependências no Kanban; um integrador escreve compartilhado, revisões e fixtures em snapshots separados. Não adicionar um grande catálogo nominal que só classifica arquivos sem gerar comandos confiáveis.

Estados sugeridos: inventariado, demanda confirmada, contrato definido, adapter implementado, geração testada, runtime exercitado, limites documentados. Nunca usar 'suportado' para apenas identificar manifest ou compilar o Mudarro.

## Registro atual

MUD-015 planeja; MUD-016 amplia testes atuais. Colisões de scripts JS foram corrigidas com regressão; as demais lacunas requerem contrato explícito e priorização. Não há suporte automático novo Rust/PHP/Java/.NET nesta rodada. Critérios e fixtures em [matriz](support-matrix.md).
