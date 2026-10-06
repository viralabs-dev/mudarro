# Proposta exata de plugins Maven — autorização separada

**Plugins não executados.** Ferramentas do loteA autorizadas20:35:52 agora foram instaladas: JDK/Java stdlib real e Maven --version validados; nenhum lifecycle Maven. Plugins continuam fora da autorização. Nenhum JAR/plugin baixado ou executado por esta proposta.

## Cenário mínimo e teste real sem dependências de aplicação

A documentação oficial versionada [Surefire3.5.5 POJO](https://maven.apache.org/surefire-archives/surefire-3.5.5/maven-surefire-plugin/examples/pojo-test.html) admite classes `*Test`, métodos `test*` e assertions JDK, sem JUnit/TestNG aplicação. Propor provider explícito `org.apache.maven.surefire:surefire-junit3:3.5.5` no plugin; componentes extras `common-junit3:3.5.5` e `common-java5:3.5.5`. POMs declaram JUnit3.8.1 como **provided**; não inserir dependência JUnit no sample. CompatibilidadePOJO/provider semJUnit precisa runtimeverificação futura, não comprovada sópor docs. Não escrever classes falsas `junit.*` para contornar.

Sample POM fixa resources3.3.1/compiler3.15.0/surefire3.5.5, releaseJava explícito e pluginprovider; zero appdeps. Um POJO testa21*2==42 e outrocontrole negativo lançaAssertionError. Exigir relatórioSurefire com testes executados>0 e negativo real exitnãozero; `mvn test` NoTests não é aprovação. Não package/install/clean/exec/help/dependencyplugins implícitos.

## Pedido agrupado concreto para revisão

Após decisão explícita sobre hashes fracos e fechamento dos três providers, autorizar somente download em cache privado e execução de resources/compiler/surefire nas versões abaixo e dependências exatas deste inventário, apenas fixture interno aprovado. Maven/JDK privados, settings global/user próprios semcredenciais, repo/cache/user.home privados, HTTPS MavenCentral, sem sudo/global/shellRC. `mvn --batch-mode -gs <private> -s <private> -Dmaven.repo.local=<private> -Duser.home=<private> compile` e `test` somente após conferênciaresolvedinventory; bináriosdirectjava para self-test são validadostool-only separadamente. Stop em novo artifact/versão/source; não adicionar downloads automaticamente.

**Gate restante:** 35artefatos do relatório +3candidatosprovider não equivalem closure executávelprovada. Apenas5/35base têm SHA256/512 oficial; demais30 e3providers somenteSHA1. SHA1 detecta corrupção mas não é gateforte. Pedir escolha explicitamente: mantergateforte e adiar plugins, ou aceitar provenanceCentralHTTPS+SHA1 limitado para esses artifactsexatos e exigir inventário/checksumlocalSHA256 arquivado apósdownload. SHA256local não vira checksumoficial. Não instalar enquanto essa escolha/autorização não vier.

## Inventário exato publicado — 35 artefatos

| Coordinate | Checksum oficial / força | Origem |
| --- | --- | --- |
| `org.apache.maven.shared:maven-shared-incremental:1.1` | `sha1:9d017a7584086755445c0a260dd9a1e9eae161a5` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/shared/maven-shared-incremental/1.1/maven-shared-incremental-1.1.jar.sha1) |
| `org.apache.maven.shared:maven-shared-utils:3.4.2` | `sha1:bfa28296272a5915b08de9f11f34a94b0a818fd0` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/shared/maven-shared-utils/3.4.2/maven-shared-utils-3.4.2.jar.sha1) |
| `org.codehaus.plexus:plexus-compiler-api:2.16.2` | `sha256:db34d13c8d688063a946922f4de448c909ba43fec355ea15514495f33072b031` | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-compiler-api/2.16.2/plexus-compiler-api-2.16.2.jar.sha256) |
| `org.codehaus.plexus:plexus-compiler-manager:2.16.2` | `sha256:99630ac196571a2754baa143a12723795b13020631902a6dff18dfb997e59b0e` | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-compiler-manager/2.16.2/plexus-compiler-manager-2.16.2.jar.sha256) |
| `org.codehaus.plexus:plexus-java:1.5.2` | `sha256:1e6a4298e145c1e23af430b04ac53d76dc11077e0f3d36ef9c027ce790d96505` | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-java/1.5.2/plexus-java-1.5.2.jar.sha256) |
| `org.codehaus.plexus:plexus-utils:4.0.2` | `sha1:9526a9548b302572f23337fcc217fb4cc713b9c3` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-utils/4.0.2/plexus-utils-4.0.2.jar.sha1) |
| `org.ow2.asm:asm:9.9.1` | `sha256:6f3828a215c920059a5efa2fb55c233d6c54ec5cadca99ce1b1bdd10077c7ddd` | [Central metadata](https://repo.maven.apache.org/maven2/org/ow2/asm/asm/9.9.1/asm-9.9.1.jar.sha256) |
| `org.codehaus.plexus:plexus-compiler-javac:2.16.2` | `sha256:e48141c146d6cb96619aafb07b2e10e1ac08f339e96ae4067ddd9c2d0f626672` | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-compiler-javac/2.16.2/plexus-compiler-javac-2.16.2.jar.sha256) |
| `com.thoughtworks.qdox:qdox:2.2.0` | `sha1:39651eb3ce73d6e506490ea352e1e13eab6b55e8` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/com/thoughtworks/qdox/qdox/2.2.0/qdox-2.2.0.jar.sha1) |
| `commons-io:commons-io:2.11.0` | `sha1:a2503f302b11ebde7ebc3df41daebe0e4eea3689` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/commons-io/commons-io/2.11.0/commons-io-2.11.0.jar.sha1) |
| `javax.inject:javax.inject:1` | `sha1:6975da39a7040257bd51d21a231b76c915872d38` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/javax/inject/javax.inject/1/javax.inject-1.jar.sha1) |
| `org.slf4j:slf4j-api:1.7.36` | `sha1:6c62681a2f655b49963a5983b8b0950a6120ae14` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/slf4j/slf4j-api/1.7.36/slf4j-api-1.7.36.jar.sha1) |
| `org.apache.maven.surefire:maven-surefire-common:3.5.5` | `sha1:253a483f39cf49e30e2a1f5ba9d695c9cd5fb268` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/maven-surefire-common/3.5.5/maven-surefire-common-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:surefire-api:3.5.5` | `sha1:51e809679c002c53b3ae4a6f7330be4c031c5961` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-api/3.5.5/surefire-api-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:surefire-extensions-api:3.5.5` | `sha1:095a518fba2d699e9719f7332d654b669bf8513a` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-extensions-api/3.5.5/surefire-extensions-api-3.5.5.jar.sha1) |
| `org.apache.maven.resolver:maven-resolver-api:1.4.1` | `sha1:ceee6b7ea1bc252afa585fa32f76c2cda206bdcd` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/resolver/maven-resolver-api/1.4.1/maven-resolver-api-1.4.1.jar.sha1) |
| `org.apache.maven.resolver:maven-resolver-util:1.4.1` | `sha1:3f6d4f4bc3e24b46a776b47ccfeaed9d2ed01549` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/resolver/maven-resolver-util/1.4.1/maven-resolver-util-1.4.1.jar.sha1) |
| `org.apache.maven.shared:maven-common-artifact-filters:3.4.0` | `sha1:25855b1fa26fa5a1b387375de4b14ac39df23981` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/shared/maven-common-artifact-filters/3.4.0/maven-common-artifact-filters-3.4.0.jar.sha1) |
| `org.apache.maven.surefire:surefire-booter:3.5.5` | `sha1:20e766e6aa17c9af68c68730cbb7b7974beff842` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-booter/3.5.5/surefire-booter-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:surefire-extensions-spi:3.5.5` | `sha1:6e200c98153c3ac79e8b3102bf8b5395728dfada` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-extensions-spi/3.5.5/surefire-extensions-spi-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:surefire-logger-api:3.5.5` | `sha1:28302b8c8bec2a0c38fec085c55a1b097b31a391` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-logger-api/3.5.5/surefire-logger-api-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:surefire-shared-utils:3.5.5` | `sha1:9bccbee7fbcbc7f853aab24aff64e122bac0dedc` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-shared-utils/3.5.5/surefire-shared-utils-3.5.5.jar.sha1) |
| `commons-codec:commons-codec:1.19.0` | `sha1:8c0dbe3ae883fceda9b50a6c76e745e548073388` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/commons-codec/commons-codec/1.19.0/commons-codec-1.19.0.jar.sha1) |
| `commons-io:commons-io:2.21.0` | `sha1:52a6f68fe5afe335cde95461dd5c3412f04996f7` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/commons-io/commons-io/2.21.0/commons-io-2.21.0.jar.sha1) |
| `org.apache.commons:commons-compress:1.28.0` | `sha1:e482f2c7a88dac3c497e96aa420b6a769f59c8d7` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/commons/commons-compress/1.28.0/commons-compress-1.28.0.jar.sha1) |
| `org.apache.commons:commons-lang3:3.20.0` | `sha1:65897b3e5731220962e659e001904af3c3cbeba9` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/commons/commons-lang3/3.20.0/commons-lang3-3.20.0.jar.sha1) |
| `org.apache.maven.shared:maven-shared-utils:3.3.4` | `sha1:f87a61adb1e12a00dcc6cc6005a51e693aa7c4ac` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/shared/maven-shared-utils/3.3.4/maven-shared-utils-3.3.4.jar.sha1) |
| `org.apache.commons:commons-lang3:3.12.0` | `sha1:c6842c86792ff03b9f1d1fe2aab8dc23aa6c6f0e` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/commons/commons-lang3/3.12.0/commons-lang3-3.12.0.jar.sha1) |
| `org.apache.maven.shared:maven-filtering:3.3.1` | `sha1:7b613072bcce1d949b6d82f714af08b4535aae2b` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/shared/maven-filtering/3.3.1/maven-filtering-3.3.1.jar.sha1) |
| `org.codehaus.plexus:plexus-utils:3.5.1` | `sha1:c6bfb17c97ecc8863e88778ea301be742c62b06d` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-utils/3.5.1/plexus-utils-3.5.1.jar.sha1) |
| `org.codehaus.plexus:plexus-interpolation:1.26` | `sha1:25b919c664b79795ccde0ede5cee0fd68b544197` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/codehaus/plexus/plexus-interpolation/1.26/plexus-interpolation-1.26.jar.sha1) |
| `org.sonatype.plexus:plexus-build-api:0.0.7` | `sha1:e6ba5cd4bfd8de00235af936e7f63eb24ed436e6` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/sonatype/plexus/plexus-build-api/0.0.7/plexus-build-api-0.0.7.jar.sha1) |
| `org.apache.maven.plugins:maven-compiler-plugin:3.15.0` | `sha1:22d13874d815e28e90134a5977eb2f321fa6b35d` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/plugins/maven-compiler-plugin/3.15.0/maven-compiler-plugin-3.15.0.jar.sha1) |
| `org.apache.maven.plugins:maven-resources-plugin:3.3.1` | `sha1:5a0e59faaaec9485868660696dd0808f483917d0` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/plugins/maven-resources-plugin/3.3.1/maven-resources-plugin-3.3.1.jar.sha1) |
| `org.apache.maven.plugins:maven-surefire-plugin:3.5.5` | `sha1:4f7f8e64b4d6a64ac19beb353ae23a41ead27208` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/plugins/maven-surefire-plugin/3.5.5/maven-surefire-plugin-3.5.5.jar.sha1) |

## Três candidatos extras para provider POJO

| Coordinate | Checksum / força | Origem |
| --- | --- | --- |
| `org.apache.maven.surefire:surefire-junit3:3.5.5` | `sha1:1abdaed598b9fbba07212e80ea286e53c0038a76` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/surefire-junit3/3.5.5/surefire-junit3-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:common-junit3:3.5.5` | `sha1:29f25e8d5e42f8dacf667e95a01a0d3d2a06ef51` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/common-junit3/3.5.5/common-junit3-3.5.5.jar.sha1) |
| `org.apache.maven.surefire:common-java5:3.5.5` | `sha1:082ac7ba576a079a06fb977b1d6cba95b51c421e` — fraco | [Central metadata](https://repo.maven.apache.org/maven2/org/apache/maven/surefire/common-java5/3.5.5/common-java5-3.5.5.jar.sha1) |

Detalhes legíveis/máquina: [request JSON](plugin-approval-request.json), [lista oficial por seção/scope](versioned-full-dependency-list.json), [checksumsbase](runtime-artifact-checksum-metadata.json), [provider metadata](pojo-provider-metadata.json). POMs/documentação oficiais arquivados; nenhum pacote/código baixado. Provided/test do próprio plugin não entram automaticamente no runtimefixture. Distinção: --version valida ferramenta; javac/console testaJDK; compile testaMavenplugins; testcomPOJO testaSurefire real; nenhumdeles comprova Gradle/Kotlin/frameworks.
