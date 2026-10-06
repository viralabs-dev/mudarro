# JVM e .NET — mínimos e bloqueios de fechamento

Somente metadados/POM/checksums textuais/docs oficiais; nenhum JAR/NuGet/SDK baixado ou executado. Maven3.10.0/JDK e SDK.NET10.0.401 já verificados no plano raiz, não repetidos aqui.

## Maven

Lista fechada de plugins top-level para fases compile/test, fixar no POM (não depender defaults Maven): org.apache.maven.plugins:maven-resources-plugin:3.3.1; maven-compiler-plugin:3.15.0; maven-surefire-plugin:3.5.5. POMs oficiais e hashes JAR.sha1 textuais foram obtidos. Endpoints JAR.sha512 destes três retornaram404; não inventar SHA256/SHA512. Relatórios oficiais versionados de dependências arquivados quando disponíveis.

**Não é ainda closure transitiva executável verificada.** POMs incluem parents/properties/dependencyManagement/scopes e providers Surefire selecionados em runtime; listar só dependencies diretas não resolve effective graph. declared-dependencies.json preserva expressões não resolvidas claramente. Dependências scope=test do próprio plugin não devem virar dependências do sample; provided normalmente entregues pelo Maven. Sem Maven instalado/executado nesta pesquisa, não afirmar resolução/cache/compatibilidade completa.

Zero dependências de aplicação elimina JUnit/TestNG do sample. `mvn test` com zero testes é gate de lifecycle, NÃO teste real. Recomendar java stdlib console self-test explícito após compile/testCompile (testes positivos/negativos com exitcode), mantendo resultado separado. Para teste via Surefire real escolher provider/framework pinned e fechar deps adicionais antes autorização; não incluir JUnit implicitamente.

Proposta mínima após autorização: Maven -gs <private-globalsettings> -s <private-usersettings> -Dmaven.repo.local=<private-cache> -Duser.home=<private-java-home> --batch-mode compile; depois java -cp target/classes Sample --self-test. Nenhum mvn package/install/dependency:go-offline/clean implícito: adicionariam plugins além da lista. Rede exclusivamente MavenCentral, stop caso resolverartifactfora allowlist; arquivar resolved paths/checksums/versions e obter approvalclosureantes executar qualquer plugin ainda não auditado. Para lista transitiva fechada exigir effective POM+graph com ferramenta verificada em estágio separadamente autorizado, não confiar na resolução parcial manual.

Fontes: https://repo.maven.apache.org/maven2/org/apache/maven/plugins/ ; https://maven.apache.org/plugins-archives/maven-compiler-plugin-3.15.0/dependencies.html ; https://maven.apache.org/plugins-archives/maven-resources-plugin-3.3.1/dependencies.html ; https://maven.apache.org/surefire-archives/surefire-3.5.5/maven-surefire-plugin/dependencies.html

## .NET

Recomendar console net10.0 com **self-test explícito**, sem Microsoft.NET.Test.Sdk/adapters/framework NuGet adicionais. global.json SDK10.0.401 rollForward:disable; NuGet.Config privado clearfeeds; csproj sem PackageReference; restore --configfile privado, build/run --no-restore, DOTNET_CLI_HOME/NUGET_PACKAGES/telemetria isolados como plano raiz. Código stdlib afirma21*2==42 e caso negativo retorna7. Resultado rotulado compile+console self-test; NÃO afirmar dotnet test/framework execution. Não dotnet new test, workload, certificados, globaltools ou downloads extraSDK.

A documentação Microsoft confirma que dotnet test/VSTest normalmente usa host/framework viaNuGet: https://learn.microsoft.com/en-us/dotnet/core/tools/dotnet-test-vstest . Portanto a lista adicional recomendada agora é vazia. Se usuário quiser framework real, nova lista fechada provider+adapter+testSDK+transitivas antes download.

Instalações/runtime NÃO autorizados por esta pesquisa. Permissão agrupada futura: bináriosJDK/Maven/.NET privados verificados; downloadplugin apenasclosureauditada; console self-test explícito sem pacotes extras; falhas não autorizam resolverdepsinstalarmais automaticamente.

## Atualização: tabelas transitivas oficiais verificadas

versioned-full-dependency-list.json agora contém listas integrais com section e scope dos relatórios oficiais versionados, incluindo Project Transitive Dependencies. runtime-artifact-checksum-metadata.json reúne35artefatos únicos compile/runtime+3plugins. Todos35 têm checksumSHA1 oficial Central; apenas5 têm SHA256/SHA512 disponível. Não tratar SHA1 como gate criptográfico forte. Provided/test estão preservados na lista integral mas excluídos dos35 runtime. A lista é closure publicada destes plugins, não resolutionruntime provada com Maven3.10.0 nem closure dos providersSurefire descobertos dinamicamente. Não dizer lista totalmentefechadaexecutável.

Permissão alternativa mínima: instalar somente JDK/Maven privados jácomhashforte, rodar --version e preparar resolução metadados effectivePOM em estágio futuro explicitamente autorizado, sem chamar lifecycleproject/plugins. Se Maven effectivePOM usar help:effective-pom, isso executa pluginhelp e exige aprovação/hashclosureprópria; não rotular como leitura passiva por ser comando de diagnóstico. JVM directjavac+java selftest evita plugins mas não comprova Mavenbuild.
