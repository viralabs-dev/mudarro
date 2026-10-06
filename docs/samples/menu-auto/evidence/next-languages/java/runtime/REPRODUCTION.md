# JDK runtime and Maven static contract / JDK real e contrato Maven estático

Authorized exact distributions: Temurin JDK25.0.4.1+1 and Apache Maven3.10.0. Official source URLs were downloaded unchanged; requested SHA256/SHA512 matched. Archive member paths and link targets were verified contained before private extraction. See artifact-validation.json; distributions remain outside repository samples.

Eighteen actual checks passed. JDK javac compiles Main.java, java asserts 21*2==42 and prints REAL_JDK_SELFTEST_42; deliberate failure exits7. Mudarro runs only explicitly configured javac/java actions and their wrapper. Config bytes and repeated generation hashes were preserved.

Maven executed ONLY -version with private JAVA_HOME/PATH/user.home/repository. Maven compile/test were discovered, explicitly selected and PREVIEWED ONLY; no lifecycle, help plugin, dependency resolution, build/test plugin or repository artifact was executed/downloaded. Maven repository is empty. JDK self-test does not prove Maven build/test support. Gradle, PHP and Ruby runtimes were not authorized/executed in this task.

Reproduce source fixture checks with runtime.py after replacing private root/bin paths as appropriate. Do not execute the selected Maven compile/test actions: plugin closure remains outside authorization. sample-list.json contains actual sources/config/wrappers, excludes compiled classes/downloads/tools/caches/state; custom config embeds private absolute JDK paths and requires deliberate adjustment when moved.

Dezoito verificações reais aprovadas. Compilação/self-test são javac/java explícitos, não testes Maven. Maven somente -version; ações compile/test selecionadas foram apenas visualizadas. Nenhum plugin ou dependência Maven baixado. Hashes e caminhos foram verificados antes da extração; instalações privadas sem alteração de configurações globais. Samples com caminhos JDK absolutos precisam revisão ao copiar.
