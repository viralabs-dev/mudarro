# Runtime real Rust/Cargo1.99.0 privado

PASS:46 comandos reais. Rust/Cargo1.99.0 instalado apenas no prefixo /tmp autorizado. Tar HTTPS oficial datado: SHA256891c6366d7100feda0bca4c03ce63f3c9ac827cbebbc283e7433061d42c6a376 conferido contra pin aprovado e checksum oficial antes da extração.67711 membros inspecionados, sem links externos/arquivos especiais. install.sh standalone apenas rustc,cargo,rust-std e --disable-ldconfig; sem rustup,sudo,perfil/configuração global.

Cargo usou CARGO_HOME/CARGO_TARGET_DIR/RUSTUP_HOME/TMPDIR privados, CARGO_NET_OFFLINE=true e --offline/--locked nos comandos diretos após generate-lockfile real offline. Somente crates locais library/app; nenhum arquivo de dependência registry. Cargo.lock real arquivado. Nenhum download de dependências de aplicação.

Build/test e repetições passaram; teste da biblioteca confere21, bin selecionado imprimiu42. Bin de falha saiu7; CLI/wrapper Mudarro saíram1 com exit status7. Bin inexistente e erro intencional de compilação saíram101. build.rs negativo não foi executado. Evidência comprova este sample, sem ampliar suporte genérico a grafo/plugins/build scripts.

Mudarro scan apresentou seis sugestões build/test opt-in, sem comandos/start automáticos. Seleção explícita, rejeição de overwrite, generate idempotente, previews, seis runs e seis wrappers passaram. Startup manual seleciona owned-server; falha controlada seleciona owned-fail. Serviço próprio passou up/status/logs/restart/down/stopped, logs com OWNED_SERVER_42, configuração preservada. Nenhum processo próprio restante antes da gravação pelo integrador; nenhum processo pessoal sinalizado.

Configuração runtime autoritativa: sample/mudarro.json. manual-start.example.yaml é exemplo legado de bin finito, não o serviço persistente validado. Excluídos prefix/toolchain, tar/src,caches,target,binários,.mudarro state e tokens de identidade. Originais locais preservados em /tmp/mudarro-approved-runtimes.6AsmDd/rust.

Pai cuida de mídia PTY, gates gerais, Vault/publicação. Nenhum código/teste de produção alterado nesta rodada runtime; sem navegador/GUI.
