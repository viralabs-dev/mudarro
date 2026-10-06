# MUD-042 — runtime .NET privado verificado

**31 comandos reais PASS**, SDK **10.0.401**, sem mudança de produção. Escopo: compilação console net10.0 e self-test explícito com biblioteca padrão; **não houve dotnet test nem execução de framework de testes**.

[Registros reais](command-results.json), [resumo](summary.json), [metadados oficiais](official-sdk-metadata.json), [hash/extração](installation-verification.json), [ambiente](runtime-environment.json), [restore](restore-policy.json) e [processos próprios](process-check.json). O tar oficial teve SHA512 correspondente ao autorizado e5.631 membros contidos. SDK/CLI home/caches ficaram privados em /tmp; HOME herdado foi preservado. Sem sudo, certificados, workloads, globaltools ou serviços persistentes.

[Sample console](samples/console/Sample.csproj) sem PackageReference, [SDK pin](samples/console/global.json) rollForward:disable e [NuGet.Config](samples/console/NuGet.Config) com feeds/fallbacks vazios. Restore resolveu zero bibliotecas/pacotes externos; caches privados sem nupkg. Build usa no-restore, build servers desabilitados e compilação sem servidor compartilhado.

Scan repetido, seleção explícita apenas build, comandos console/self-test/falha declarados, generate2 com hashes iguais, preservação de configuração, preview/run/wrappers, repetição, falha real de compilação com fonte restaurada e recuperação passaram. Self-test afirma21*2==42; falha controlada retorna7 no console e1 no Mudarro/wrapper. Múltiplos projetos, XML inválido e manifest excluído oversized foram gates estáticos. Nenhum framework/start inferido.

Uma expectativa inicial do coletor foi corrigida: XML inválido produz aviso e fallback custom sem comandos C#, conforme comportamento do aplicativo. Execução anterior com HOME privado foi substituída por gate completo final31 comandos preservando HOME herdado. Históricos ficam privados; apenas evidência final exportada. Não foi necessário corrigir adapter.

[Instruções completas](README.md#reproduction) e [driver real](runtime.py). Para reproduzir use SDK privado já verificado, cópia dos cinco arquivos do sample, ambiente isolado preservando HOME, restore com NuGet.Config explícito, generate e run das ações app-csharp:build/console/self-test. A ação controlled-failure falha propositalmente. Nenhum SDK/bin/obj/cache ou fixture4MiB foi copiado. Mídias reais ficam sob responsabilidade do integrador.
