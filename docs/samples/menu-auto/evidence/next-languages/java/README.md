# Java/Maven static contract (MUD-040)

Samples are internal, with no application dependencies. The single-project adapter offers opt-in `mvn -f pom.xml compile` and `test`, with the service directory as cwd. Aggregators receive diagnostics only; contained modules are detected independently and require explicit selection. Gradle files are executable DSL evidence, never evaluated; no tasks or Kotlin support inferred.

Static gates validate XML/DTD/entity rejection, bounded regular files and symlinks, missing/dynamic/outside modules, exclusion before reading, no-execution traps, selection, generation twice and preservation of existing configuration. Start/main/plugins/profiles are not inferred. Namespaces are accepted by XML local names.

No javac/Maven runtime was executed. A future Maven invocation may download lifecycle plugins even though this sample has no application dependencies. Toolchain/plugin versions, private cache and network approval remain separate gates. `mvn test` does not prove Java tests ran in this zero-test sample.

Reproduce static gates with `go test ./test/internal/mudarro -run 'Test(Maven|Gradle)' -race -count=1`. Run from the repository with its approved Go cache/module configuration. No installation is implied.
