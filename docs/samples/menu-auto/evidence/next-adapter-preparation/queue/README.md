# Auditoria da próxima fila — nove atividades

Leitura atual: todos os nove cartões permanecem **A fazer**, sem responsável atribuído. Esta auditoria não assumiu nem iniciou atividades. Samples internos padronizados já aprovados eliminam a antiga dependência de repositório externo; demanda comercial e suporte completo não foram inventados.

| Card | Bloqueio atual | Próximo passo |
| --- | --- | --- |
| MUD-001 — Smoke WSL | Host Linux nãoWSL; wsl.exe ausente | Solicitar acesso a ambiente WSL já existente; sem substituir por Linux/cross-build |
| MUD-002 — PostgreSQL nativo host | initdb/pg_ctl/postgres/psql ausentes na checagem | Preparar plano oficial pinned de instalação isolada para autorização; nunca tocar cluster pessoal |
| MUD-005 — Frameworks/runtime modules | Python disponível; Flask/FastAPI/pytest/uvicorn ausentes em distribuição desse Python | Pode preparar fixtures sintéticos de metadados e testes scanner; não afirmar runtime; definir subset framework antes código |
| MUD-039 — Rust/Cargo | cargo/rustc ausentes | Contrato + fixtures internos zero deps são viáveis; runtime requer plano pinned e autorização instalação |
| MUD-040 — Java/Kotlin | JRE8u502 somente; javac/mvn/gradle ausentes | Recomendar primeiro Maven/Java zero deps; Kotlin/Gradle são capacidade adicional, não support claim automático; planejar JDK+Maven isolados |
| MUD-041 — PHP/Composer | php/composer ausentes | Fixtures internos JSON/testes são viáveis; preparar versões/toolchain isolado antes runtime |
| MUD-042 — C sharp/.NET | dotnet ausente | Contrato mínimo SDK console zero package deps; manifests internos e testes semruntime; restore pode acessar rede deve ser declarado no plano |
| MUD-043 — Ruby | ruby/bundle ausentes | Fixtures internos e negativos são possíveis; parser subset/sem avaliação requer contrato explícito, runtime requer instalação autorizada |
| MUD-046 — Podman local | podman/podman-compose ausentes | Plano instalação precisa distinguir toolchain de requisitos kernel/rootless/storage; não alterar segurança/global config; CIcontainers prova distinta |

Decisões de contrato mínimas: definir subset JVM (Maven/Java recomendado primeiro), Ruby sem avaliação de DSL e frameworks Python com metadados/sugestões opt-in. Os demais bloqueios principais são ambientais; autorização de instalação exige planos concretos com versão/origem/local/comandos, não uma permissão genérica antecipada.

WSL necessita ambiente WSL real. PostgreSQL nativo e Podman local não são comprovados por CI containers. JRE presente não equivale a JDK/toolchain JVM. Checagem de frameworks consultou somente metadados de distribuição Python, sem importar aplicativo.

Matriz mínima futura por adapter: manifest inválido/ambíguo, excludes e symlinks, scan sem execução, ownership/argv/cwd, targets explícitos sem startup inventado, generate repetido e preservação manual, runtime real de toolchain pinned em fixture interno sem deps externas, falha real e limpeza. Casos manifest/code DSL devem permanecer identificação limitada até contrato seguro aprovado.

[Relatório estruturado](report.json), [disponibilidade não exaustiva](availability.json), [versão Java/metadados Python](runtime-details.json). Fontes exatas dos cartões estão no JSON. Nenhuma busca ampla em arquivos pessoais, instalação, alteração de configuração ou publicação.
