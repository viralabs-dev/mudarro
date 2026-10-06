# Proposta agrupada de ambientes Mudarro

## Atualização de autorização e execução — lote A

A autorização20:35:52 permitiu Rust/JDK/Maven/.NET privados, agora instalados e validados nos subsets registrados em [evidência runtime](../../../next-languages/runtime/README.md). Maven somente --version, sem plugins/lifecycle. B (compilaçõesPHP/Ruby/PostgreSQL/pré-requisitos), C (Podman) e D (WSL) não autorizados/executados por esse lote. A pesquisa inicial abaixo é histórica; suplementos continuam delimitando hashes/closures/gates.

[Suplemento de finalização](finalization/README.md) · [Hashes oficiais](finalization/hashes/report.md) · [Checksums Ubuntu e limites](finalization/hashes/ubuntu-supplement.md) · [Closure provider](finalization/provider/report.md) · [Maven/.NET](finalization/jvm-dotnet/report.md). O texto original abaixo registra a pesquisa inicial; os suplementos corrigem hashes recuperados e delimitam os gates ainda pendentes. Nenhuma instalação autorizada ou executada.

Verificação: 2026-10-06, Pop!_OS x86_64/glibc2.39. **Somente pesquisa/leitura/documentação. Todos os comandos abaixo são propostas, não foram executados.** Contratos e publicação044045 aguardam Daniel. Sem sudo, downloads de distribuições, compilação, instalação, namespaces ou serviços nesta frente. Consultas HTTP foram feitas pela ferramenta web; curl local falhou DNS, sem tentar escalada/permissão.

## Decisão agrupada

| Lote | Pins verificados em fontes oficiais | Aprovação proposta / limite |
|---|---|---|
| A — binários privados | Rust/Cargo1.99.0; Temurin JDK25.0.4.1+1 + Maven3.10.0; .NET SDK10.0.401 | Download verificado e instalação somente em diretório privado. Sem sudo, perfil shell ou serviço. Runtime futuro aprovado separadamente do planejamento. |
| B — compilar em prefixo privado | PHP8.5.11 + Composer2.10.3; Ruby4.0.7 + Bundler4.0.22 + libyaml0.2.5; PostgreSQL18.6; pré-requisitos Bison3.8.2/Flex2.6.4 | Autorizar também estas dependências privadas. Build sem sudo, CPU/disco temporários. Falha de configure não autoriza apt/brew nem novas dependências. PostgreSQL só cluster próprio temporário, Unix socket privado. |
| C — Podman, decisão distinta | Upstream Podman6.1.3 e provider podman-compose1.6.0 | Distribuição oficial não oferece engine Linux x64 pronto nos assets consultados. Pacote local apt candidato4.9.3+ds1-1ubuntu0.2 é versão diferente. Não fingir pin6.1.3 em apt. Instalação sistêmica de uidmap e engine/dependências ou cadeia privada de build exigem outra proposta/aprovação. |
| D — ambiente | WSL real | PopOS atual não é WSL; wsl.exe ausente, nenhum WSL_* relevante/kernel Microsoft. Nenhum host remoto/Windows consultado. Daniel precisa identificar um ambiente WSL já autorizado; não criar VM/WSL. |

## Isolamento comum

Quando e somente quando aprovados, criar nova raiz (exemplo abaixo), sem reaproveitar instalações/pastas pessoais. Não alterar HOME, perfis, alternativas, brew, apt, Git ou serviços. Variáveis específicas são passadas apenas aos subprocessos. Guardar versão/origem/hash/exit em manifesto. Prefixos privados continuam locais, não entram em publicação.

```bash
MUD_RUNTIME_ROOT=$(mktemp -d /tmp/mudarro-approved-runtimes.XXXXXX)
chmod 700 "$MUD_RUNTIME_ROOT"
mkdir -p "$MUD_RUNTIME_ROOT"/{downloads,src,prefix,cache,evidence,fixtures}
```

Antes de extrair, verificar checksum oficial e membros de cada tar (sem path absoluto/.., sem sobrescrever fora da raiz). Não usar curl|sh. Para arquivos cujo hash não ficou acessível nesta pesquisa, obter e salvar o checksum/assinatura oficial **antes de executar ou extrair**; divergência/interrupção encerra esse item. Versão verificada não equivale a artefato baixado/validado. Caminhos de downloads Rust/PG abaixo seguem convenções oficiais, mas disponibilidade HTTP e hashes desses arquivos ainda precisam do gate de download autorizado.

## A1 Rust/Cargo1.99.0

[Release oficial1.99.0](https://blog.rust-lang.org/2026/10/01/Rust-1.99.0/), [Cargo instalação](https://doc.rust-lang.org/cargo/getting-started/installation.html), [métodos oficiais](https://rust-lang.github.io/rustup/installation/other.html). Escolha tar oficial standalone: evita instalar rustup ou editar perfis. GCC/linker já presentes. URL proposta: `https://static.rust-lang.org/dist/rust-1.99.0-x86_64-unknown-linux-gnu.tar.xz`; checksum `.sha256`, assinatura `.asc`. Hash oficial recuperado no fechamento: `891c6366d7100feda0bca4c03ce63f3c9ac827cbebbc283e7433061d42c6a376`; metadados em finalization/hashes/rust.metadata.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl --fail --location --proto '=https' --tlsv1.2 -o rust.tar.xz https://static.rust-lang.org/dist/rust-1.99.0-x86_64-unknown-linux-gnu.tar.xz
curl --fail --location --proto '=https' --tlsv1.2 -o rust.sha256 https://static.rust-lang.org/dist/rust-1.99.0-x86_64-unknown-linux-gnu.tar.xz.sha256
# Validar SHA256 contra rust.tar.xz usando digest do arquivo oficial; registrar ambos.
# Após validar digest e membros:
tar -xJf rust.tar.xz -C "$MUD_RUNTIME_ROOT/src"
cd "$MUD_RUNTIME_ROOT/src/rust-1.99.0-x86_64-unknown-linux-gnu"
./install.sh --prefix="$MUD_RUNTIME_ROOT/prefix/rust-1.99.0" --disable-ldconfig
"$MUD_RUNTIME_ROOT/prefix/rust-1.99.0/bin/rustc" --version
"$MUD_RUNTIME_ROOT/prefix/rust-1.99.0/bin/cargo" --version
```

Runtime futuro: PATH privado, CARGO_HOME/cache/cargo e CARGO_TARGET_DIR/cache/target próprios; `cargo build --offline`, `cargo test --offline`, bin explicitamente escolhido. Sample sem crates externas/build.rs; caso build.rs apenas negativo scan não executável. Downloads crates adicionais não implícitos.

## A2 JDK25.0.4.1+1 e Maven3.10.0

[Temurin release](https://github.com/adoptium/temurin25-binaries/releases/tag/jdk-25.0.4.1+1), [asset oficial](https://github.com/adoptium/temurin25-binaries/releases/expanded_assets/jdk-25.0.4.1%2B1), [Maven download](https://maven.apache.org/download.cgi), [instalação](https://maven.apache.org/install.html). JDK inclui javac; JRE8 existente não satisfaz compilação. Tar Temurin Linux glibc x64, aproximadamente135MB, SHA256 `dbb698396d478e7fa2b1e50f4103324b2a99b90569ee27c33f2261f9215cf41e`. Maven SHA512 oficial `908b1501bfb420bf7c8affb855534a9c407fd6099367bfb9f2f2dcb8e9799102bffb84518cde74c679bd76870247c6528683abdd620581bffa90f95d92d175aa`.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl -fL --proto '=https' -o jdk.tar.gz 'https://github.com/adoptium/temurin25-binaries/releases/download/jdk-25.0.4.1%2B1/OpenJDK25U-jdk_x64_linux_hotspot_25.0.4.1_1.tar.gz'
printf '%s  %s\n' dbb698396d478e7fa2b1e50f4103324b2a99b90569ee27c33f2261f9215cf41e jdk.tar.gz | sha256sum -c -
curl -fL --proto '=https' -o maven.tar.gz https://dlcdn.apache.org/maven/maven-3/3.10.0/binaries/apache-maven-3.10.0-bin.tar.gz
printf '%s  %s\n' 908b1501bfb420bf7c8affb855534a9c407fd6099367bfb9f2f2dcb8e9799102bffb84518cde74c679bd76870247c6528683abdd620581bffa90f95d92d175aa maven.tar.gz | sha512sum -c -
# Depois de verificar membros:
mkdir -p "$MUD_RUNTIME_ROOT/prefix/jdk25"
tar -xzf jdk.tar.gz --strip-components=1 -C "$MUD_RUNTIME_ROOT/prefix/jdk25"
tar -xzf maven.tar.gz -C "$MUD_RUNTIME_ROOT/prefix"
JAVA_HOME="$MUD_RUNTIME_ROOT/prefix/jdk25" "$MUD_RUNTIME_ROOT/prefix/apache-maven-3.10.0/bin/mvn" --version
"$MUD_RUNTIME_ROOT/prefix/jdk25/bin/javac" -version
```

Maven runtime requer plugins mesmo sem dependências de aplicação. Sample futuro fixa versões de plugins no POM; revisão/downloads de plugins e transitivas em Maven Central com cache privado fazem parte de autorização explícita, ainda não são lista fechada nesta proposta. Usar `-Dmaven.repo.local=<cache>` e settings global/user privados sem credenciais, declarar `user.home=<cache>` para Java. Sem daemon Maven. Fase offline só depois de cache provado. Kotlin/Gradle permanece fora destes comandos.

## A3 .NET SDK10.0.401

[SDK oficial](https://dotnet.microsoft.com/en-us/download/dotnet/10.0), [metadados/hash](https://builds.dotnet.microsoft.com/dotnet/release-metadata/10.0/releases.json), [requisitos Ubuntu](https://learn.microsoft.com/en-us/dotnet/core/install/linux-ubuntu-install). Instalação tar privada sem dotnet-install script. SHA512 Linux glibc x64: `51c8b999af9e8dd9998c9edc5944e19a90788862068acd38694e098889054ce8c23d4f0c5cccfa16bf187d044562359e5ee69a9f8ad0bbe913ba90311fbce25b`. Pacotes libc6/libgcc/libstdc++/ICU74/OpenSSL3/zlib/CA consultados presentes; validação dinâmica final ainda necessária.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl -fL --proto '=https' -o dotnet.tar.gz https://builds.dotnet.microsoft.com/dotnet/Sdk/10.0.401/dotnet-sdk-10.0.401-linux-x64.tar.gz
printf '%s  %s\n' 51c8b999af9e8dd9998c9edc5944e19a90788862068acd38694e098889054ce8c23d4f0c5cccfa16bf187d044562359e5ee69a9f8ad0bbe913ba90311fbce25b dotnet.tar.gz | sha512sum -c -
mkdir -p "$MUD_RUNTIME_ROOT/prefix/dotnet-10.0.401"
# Validar membros antes:
tar -xzf dotnet.tar.gz -C "$MUD_RUNTIME_ROOT/prefix/dotnet-10.0.401"
DOTNET_ROOT="$MUD_RUNTIME_ROOT/prefix/dotnet-10.0.401" DOTNET_CLI_HOME="$MUD_RUNTIME_ROOT/cache/dotnet-home" DOTNET_CLI_TELEMETRY_OPTOUT=1 DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1 DOTNET_ADD_GLOBAL_TOOLS_TO_PATH=0 "$MUD_RUNTIME_ROOT/prefix/dotnet-10.0.401/dotnet" --info
```

Sample manual console net10.0 sem PackageReference, NuGet.Config próprio com `<clear/>` em fontes, global.json fixa10.0.401/rollForward:disable. NUGET_PACKAGES/cache privados; restore com config privado, build/run `--no-restore`. `dotnet test` com teste real geralmente exige SDK de teste/adapters/pacotes: nenhum xUnit/NUnit/MSTest automaticamente autorizado. Compile+console self-test explícito cobre subset; teste framework posterior requer lista de pacotes aprovada. Sem certificados dev/workload/globaltools.

## B0 Pré-requisitos privados de compilação

Bison3.8.2 [tag do mantenedor](https://github.com/akimd/bison/releases/tag/v3.8.2), Flex2.6.4 [release oficial](https://github.com/westes/flex/releases/tag/v2.6.4), libyaml0.2.5 [release](https://github.com/yaml/libyaml/releases/tag/0.2.5). m4/GCC/make/perl já disponíveis. Bison/flex ausentes; libyaml-dev ausente. Compilar arquivos release com configure pré-gerado, não tar de checkout Git que exige bootstrap adicional. Hashes/assinaturas precisam ser obtidos e registrados antes da execução; páginas GNU/alguns checksums falharam nesta pesquisa.

URLs propostas: `https://ftp.gnu.org/gnu/bison/bison-3.8.2.tar.xz` + `.sig`; `https://github.com/westes/flex/releases/download/v2.6.4/flex-2.6.4.tar.gz`; `https://github.com/yaml/libyaml/releases/download/0.2.5/yaml-0.2.5.tar.gz`. Depois de download verificado/extrair cada fonte em src:

```bash
# Cada bloco executado no respectivo diretório extraído:
./configure --prefix="$MUD_RUNTIME_ROOT/prefix/build-tools"
make -j2
make check
make install
# libyaml usa prefixo separado:
./configure --prefix="$MUD_RUNTIME_ROOT/prefix/libyaml-0.2.5" --disable-shared
make -j2
make check
make install
```

Mudar PATH apenas subprocesso para build-tools/bin. Não instalar apt/brew mesmo se configure falhar. Checksum desconhecido é gate, não passar automaticamente só por HTTPS.

## B1 PHP8.5.11 e Composer2.10.3

[Tar/hash PHP](https://www.php.net/downloads.php?source=Y), [CLI instalação](https://www.php.net/manual/en/install.unix.commandline.php), [Composer versão/hash](https://getcomposer.org/download/). PHP source xz14MB SHA256 `d9be75c08e8c316f4c8f4194d8fbe1750a15f6a6d9d4e3fe72082abeeb800360`. Não instalar Apache/FPM. OpenSSL/curl/zlib pkg-config presentes; usar distribuição release com parser gerado e Bison/Flex privados se requeridos. SQLite/oniguruma/readline ausentes; subset CLI não os exige, sem mbstring/SQLite/servidor.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl -fL --proto '=https' -o php.tar.xz https://www.php.net/distributions/php-8.5.11.tar.xz
printf '%s  %s\n' d9be75c08e8c316f4c8f4194d8fbe1750a15f6a6d9d4e3fe72082abeeb800360 php.tar.xz | sha256sum -c -
# Validar membros, extrair e entrar em php-8.5.11:
./configure --prefix="$MUD_RUNTIME_ROOT/prefix/php-8.5.11" --disable-all --enable-cli --disable-cgi --disable-phpdbg --enable-phar --enable-filter --enable-tokenizer --with-openssl --with-curl --with-zlib
make -j2
# Suite selecionada CLI/Phar/JSON, registrar falhas/skips; não alegar toda suite.
make install
"$MUD_RUNTIME_ROOT/prefix/php-8.5.11/bin/php" -n -v
curl -fL --proto '=https' -o "$MUD_RUNTIME_ROOT/prefix/php-8.5.11/bin/composer.phar" https://getcomposer.org/download/2.10.3/composer.phar
printf '%s  %s\n' 7a2d379d5b8ffdaa028580ef26494c36d2feef4b178d3dd1473a4dbc5e17c8d6 "$MUD_RUNTIME_ROOT/prefix/php-8.5.11/bin/composer.phar" | sha256sum -c -
"$MUD_RUNTIME_ROOT/prefix/php-8.5.11/bin/php" -n "$MUD_RUNTIME_ROOT/prefix/php-8.5.11/bin/composer.phar" --version
```

Criar launcher composer privado invocando PHP privado/phar com argv preservado; COMPOSER_HOME/COMPOSER_CACHE_DIR privados. `composer validate --no-check-publish`, scripts sample explícitos sem dependências externas. Se precisar install: somente sample confiável revisado, `--no-plugins --no-scripts --no-interaction`; depois `run-script` específico explicitamente opt-in. Sem globalinstall/selfupdate/hooks/frameworks.

## B2 Ruby4.0.7 e Bundler4.0.22

[Ruby tar/hash](https://www.ruby-lang.org/en/downloads/), [build oficial](https://www.ruby-lang.org/en/documentation/installation/), [Bundler4.0.22 e requisitos](https://rubygems.org/gems/bundler/versions/4.0.22). Bundler requer Ruby>=3.2/RubyGems>=3.4.1. Ruby tar.gz SHA256 `911ace20f90d068ca0e4dda6d0e4f0f81e52e52f2dd4f4004c721e253412e82d`. OpenSSL/libffi/zlib/GMP disponíveis; libyaml privado evita instalar pacote sistêmico. Não ativar YJIT, Rust não é requisito desse subset.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl -fL --proto '=https' -o ruby.tar.gz https://cache.ruby-lang.org/pub/ruby/4.0/ruby-4.0.7.tar.gz
printf '%s  %s\n' 911ace20f90d068ca0e4dda6d0e4f0f81e52e52f2dd4f4004c721e253412e82d ruby.tar.gz | sha256sum -c -
# Depois de verificar/extrair entrar ruby-4.0.7:
./configure --prefix="$MUD_RUNTIME_ROOT/prefix/ruby-4.0.7" --disable-install-doc --disable-yjit --with-libyaml-dir="$MUD_RUNTIME_ROOT/prefix/libyaml-0.2.5"
make -j2
make test
make install
"$MUD_RUNTIME_ROOT/prefix/ruby-4.0.7/bin/ruby" -v
"$MUD_RUNTIME_ROOT/prefix/ruby-4.0.7/bin/ruby" -ropenssl -rpsych -e 'puts OpenSSL::OPENSSL_VERSION; puts Psych::VERSION'
GEM_HOME="$MUD_RUNTIME_ROOT/cache/gems" GEM_PATH="$MUD_RUNTIME_ROOT/cache/gems" "$MUD_RUNTIME_ROOT/prefix/ruby-4.0.7/bin/gem" install bundler --version 4.0.22 --no-document --install-dir "$MUD_RUNTIME_ROOT/cache/gems" --bindir "$MUD_RUNTIME_ROOT/prefix/ruby-4.0.7/bin" --source https://rubygems.org
```

Antes de geminstall, obter .gem oficial/hash/provenance e preferir `gem install --local <arquivo verificado>`; comando online acima descreve origem/versão, não substitui gate. Não executar Gemfile/gemspec alheio; zeroGems sample aprovado/stdlib self-test. GEM_HOME/GEM_PATH/BUNDLE_USER_HOME/BUNDLE_USER_CACHE/BUNDLE_USER_CONFIG/BUNDLE_PATH privados. Não chamar Rails/Rake, ruby managers/globalconfig.

## B3 PostgreSQL18.6 nativo

[Versão/source oficial](https://ftp.postgresql.org/pub/source/v18.6/), [requisitos18](https://www.postgresql.org/docs/18/install-requirements.html), [build](https://www.postgresql.org/docs/18/install-make.html), [initdb](https://www.postgresql.org/docs/18/app-initdb.html), [pg_ctl](https://www.postgresql.org/docs/18/app-pg-ctl.html). Stable18.6, não19beta4. Flex/Bison obrigatórios conforme docs18, ausentes: B0 é pré-requisito explícito. ICU/zlib disponíveis, readline não; configure --without-readline muda só conforto psql, registrar limitação. Não usar container como prova de host.

```bash
cd "$MUD_RUNTIME_ROOT/downloads"
curl -fL --proto '=https' -o postgresql.tar.bz2 https://ftp.postgresql.org/pub/source/v18.6/postgresql-18.6.tar.bz2
curl -fL --proto '=https' -o postgresql.sha256 https://ftp.postgresql.org/pub/source/v18.6/postgresql-18.6.tar.bz2.sha256
# Hash oficial não recuperado na pesquisa; verificar digest antes de extrair.
# Depois de validar/extrair entrar postgresql-18.6, PATH privado build-tools/bin:
./configure --prefix="$MUD_RUNTIME_ROOT/prefix/postgresql-18.6" --without-readline --with-openssl
make -j2
make check
make install
```

Proposta de execução futura sem daemon persistente: usuário atual não root, cluster/socket/testSQL exclusivos. Sem TCP por padrão, socket modo0700, porta própria55432 no diretório privado (não disputa socket pessoal). Autenticação local trust restrita ao socket inacessível a outros usuários; nenhuma edição pg_hba pessoal.

```bash
mkdir -m 700 "$MUD_RUNTIME_ROOT/pg-socket"
PG_BIN="$MUD_RUNTIME_ROOT/prefix/postgresql-18.6/bin"
"$PG_BIN/initdb" -D "$MUD_RUNTIME_ROOT/pg-data" -U mudarro_sample --encoding=UTF8 --locale=C --auth-local=trust --auth-host=reject
"$PG_BIN/pg_ctl" -D "$MUD_RUNTIME_ROOT/pg-data" -l "$MUD_RUNTIME_ROOT/evidence/postgresql.log" -o "-h '' -k $MUD_RUNTIME_ROOT/pg-socket -p 55432" -w start
"$PG_BIN/psql" -X -v ON_ERROR_STOP=1 -h "$MUD_RUNTIME_ROOT/pg-socket" -p 55432 -U mudarro_sample -d postgres -c 'SELECT version(); SELECT 21+21;'
"$PG_BIN/pg_ctl" -D "$MUD_RUNTIME_ROOT/pg-data" -m fast -w stop
```

Driver futuro deve usar trap/try-finally para parar somente cluster próprio, registrar pid/hash e confirmar encerramento. Nenhum rm -rf global, serviço systemd, porta exposta, alteraçãocluster existente. Fixtures Mudarro podem exigir TCP localhost para integração: autorização distinta deve especificar bind127.0.0.1/porta própria/credenciais temporárias, sem publicar segredos.

## C Podman6.1.3 / provider1.6.0 — proposta ainda condicionada

[Release6.1.3](https://github.com/podman-container-tools/podman/releases/tag/v6.1.3), [assets](https://github.com/podman-container-tools/podman/releases/expanded_assets/v6.1.3), [instalação](https://podman.io/docs/installation), [Makefile pin](https://raw.githubusercontent.com/podman-container-tools/podman/v6.1.3/Makefile), [rootless](https://docs.podman.io/en/latest/markdown/podman.1.html), [provider1.6.0](https://github.com/containers/podman-compose/releases/tag/v1.6.0). Assets oficiais consultados têm instaladores Windows/macOS e clients; não engineLinuxx64. Não substituir por podman-remote: requer serviço remoto e não valida host.

Disponibilidade somente leitura: subuid/subgid já têm linha do usuário (números/pessoas não expostos); newuidmap/newgidmap, conmon, pasta, fuse-overlayfs ausentes; slirp4netns/runc existem. Kernel max_user_namespaces positivo não comprova usernamespace permitido; nenhum unshare executado, nenhum sysctl/AppArmor/seccomp alterado. Podman upstreambuildCGO exige dependências (seccomp/etc.) ainda não fechadas; nunca remover seccomp/build security tags para contornar.

Duas escolhas distintas para Daniel:

1. Pacotes da distribuição: candidato local **podman4.9.3+ds1-1ubuntu0.2**, uidmap1:4.13+dfsg1-4ubuntu3.2, conmon2.1.10+ds1-1build2, passt0.0~git20240220.1e6f92b-1, fuse-overlayfs1.13-1. Índice local pode estar antigo; pedir revisão de `apt-get -s install` e atualizações/CVEs antes de comando sistêmico final. `sudo apt-get install ...` não autorizado. Pacoteuidmap envolve helpers privilegiados; não copiar/setuid manualmente. Nenhuma mudança subuid necessária demonstrada.
2. Manterpinupstream6.1.3: cadeia buildprivada engine/conmon/network/storage com todasorigens/versões/hash fechadas em outra proposta, mais instalação confiáveluidmap pelo administrador. Comando upstream orientativo `make bin/podman` com GOexistente e PREFIXprivado; **não executável como plano fechado** enquanto dependências/segurança/completude não aprovadas. `make install` default pode escrever /etc e systemd; proibido usar genericamente.

Provider futuro1.6.0 em venv privado: `python3 -m venv "$MUD_RUNTIME_ROOT/prefix/podman-compose"`; pipdownload binários/wheels `podman-compose==1.6.0` e dependências exatas com hashes fechados antes de pipinstall `--no-index --require-hashes`. pyproject oficial define PyYAML/python-dotenv; transitivas/pins ainda precisam resolução, não inventados. Sem pipglobal.

Após aprovaçãofechada: XDG_CONFIG_HOME/XDG_DATA_HOME/XDG_RUNTIME_DIR0700 e TMPDIR privados, CONTAINERS_CONF/CONTAINERS_STORAGE_CONF locais, root/runroot separados. Primeiro `podman version`/`info`, só depois namespace/runtime. Workload interno/imagem OCI comdigest explícito aprovado, rede `none` se possível, semprivileged/hostnetwork/volumepessoal/socketDocker. Executar provider via caminhoabsoluto/PODMAN_COMPOSE_PROVIDER. Sem serviço `podman system service`, systemd ou linger. Teardown somente container/volume/pod exclusivo comnomeID guardado; nunca prune/reset. Imagem PostgreSQL/porta/rede demanda plano específico adicional, não coberta por instalar provider.

## Verificações feitas e pendentes

Feito: comandosdisponíveis, pkg-configlibs, pacotesruntime.NET, subuid/subgid sóbooleanousuário, kernel/WSLsinais, apt-cachecandidatos; nenhuma varsecret/externalhost/homeprojects examinado. Evidências JSON readonly anexas.

Pendente: aceitedownloads/instalação; hashes Flex/libyaml ainda não estabelecidos; Rust/PG/Bison/Bundler recuperados em finalization/hashes/report.md; buildconfigure/compatibilidade, pluginMaven/test.NET depsnãofechadas; Podmandistribuiçãoouupstream/cadeia/uidmap; WSLenvreal. Nada desta pesquisa é adapter implementado/runtimepass. Cards permanecem9Afaire/0andamento/4bloqueados/45concluídos; concluídos não reabertos. Manifesto04404589files permanece congelado e separado.


## Fechamento posterior

Leia [hashes oficiais recuperados](finalization/hashes/report.md), [provider três wheels fechados](finalization/provider/report.md), [Maven/.NET e gates](finalization/jvm-dotnet/report.md). PostgreSQL18.6 SHA256555610c24d53e4316da5b7d3fc25c279d96856d5e0e23ee308c328c5fa881d9f; Bison3.8.2 SHA2569bba0214ccf7f1079c5d59210045227bcf619519840ebfa80cd3849cff5a5bf2; Bundler4.0.22 SHA256d8d5ec84c8555e0af71db63ed7aee4d1a8fb839ec46d84212d61979242a5d75a. Provider1.6.0/PyYAML6.0.3/python-dotenv1.2.4 fechado por metadata/wheels oficiais paraCPython3.14glinuxx64. Flex/libyaml hashes oficiais ainda gates, Maventransitivas nãofechadas. Nãoexecutar os blocos anteriores com gates pendentes.
