# Finalização dos planos — somente pesquisa

Instalações e runtimes **não autorizados** por estes documentos. Nenhum distribution/executable/package artifact foi baixado; somente checksums textuais, metadados e documentação oficial.

* [Provider podman-compose fechado](provider/report.md): três wheels exatos e hashes oficiais; [requirements](provider/requirements-hashed.txt), [manifest URLs/closure](provider/closure.json). EnginePodman e rootless continuam separados.
* [Maven/.NET](jvm-dotnet/report.md): plugins top-level exatos e relatórios/POMs oficiais; closuretransitiva não completamente auditada, gate explícito. .NET console self-test semtestSDK/pacotes adicionais recomendado; não chamar isso dotnet test.
* Hashes distribuições/pré-requisitos: frente hashes/, ver relatório específico quando presente. Hash recuperado não equivale a artifact executado.

Pedido agrupado recomendável ao usuário, somente após integração dos gates pelo pai: autorizar downloads verificados/extração/instalação em nova raizprivada /tmp dos itens com origem/hash/closure completos, semsudo/perfil/configglobal; autorizar compilação isolada apenasdos prerequisitosfechados, sem fallbackapt/brew; runtimes internos explicitamente descritos. Manter fora do lote Maven pluginscomclosurependente, qualquerfontecomhashnãoverificado, Podmanengine/uidmap/security e WSLenvinexistente. Não pedir permissão genérica para resolver dependências imprevistas. Root seleciona loteconcreto revisável antes da pergunta; nenhum desses comandos foi executado.
