# Rust/Cargo: mínimo estático (MUD039)

Nove grupos/21 subcasos passaram com race. Dezessete comandos reais do Mudarro validaram scan, seleção explícita, preservação de configuração, generate duas vezes, preview e dry-run. Sentinelas cargo/rustc nunca foram chamadas. Isso comprova adapter/CLI estático; não compilação/runtime Rust. Nenhum toolchain foi instalado.

Cargo.toml regular, limitado e contido é lido pela biblioteca TOML existente. Package nomeado, paths explícitos de lib/bin e members/default-members literais são validados. Não há resolução de AST, targets implícitos Cargo, grafo de dependências, valores herdados do manifest, resolver/toolchain, globs ou cargo metadata. Globs recebem diagnóstico; declarações inválidas/ambíguas/externas/excluídas ficam pending sem sugestões executáveis automáticas da declaração problemática. Membros válidos continuam serviços independentes. build.rs nunca executa no scan.

Serviços rust/cargo têm Commands vazio; build/test são sugestões opt-in com cwd próprio. Sem inferência de install/run/start/servidor/framework/banco. Startup e bin dependem de configuração explícita. Ownership do workspace é informativo em Pending; WorkspaceRoot, restrito a JS no schema, não foi alterado. Workspace virtual tem suffix de ID próprio; colisão raiz+membro app foi testada em scan/init.

Sample contém library/app, default app, dependência de path local, valor21 e app esperado42, além de bin de falha esperado7 para futuro runtime autorizado. Nenhuma dessas saídas foi executada nesta entrega. manual-start.example.yaml é exemplo manual, não inferência. negative-build-script demonstra descoberta sem avaliar build.rs.

Sem target/, Cargo.lock inventado, caches ou ferramentas instaladas. Runtime exige plano de toolchain exato e autorização futura; registrar essa limitação ao integrar os artefatos.

Final owner regression: nonempty package.workspace is Pending with no suggestions; no owner file is read and no nearest-ancestor ownership is inferred, including contained, oversized, excluded, symlink and outside paths.
