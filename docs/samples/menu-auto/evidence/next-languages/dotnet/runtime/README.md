# MUD-042 — verified private .NET runtime

**PASS: 31 genuine recorded commands** with SDK **10.0.401**, Mudarro binary SHA256 `86210bde93100944d280caa3dbd2394a88cea6f0287c2041573328d86e634d91`. This covers console compilation and an explicit standard-library self-test; **dotnet test and test frameworks were not executed**.

[Command records](command-results.json) contain actual argv, cwd, stdout, stderr, exits and durations. [Summary](summary.json), [official metadata](official-sdk-metadata.json), [archive verification](installation-verification.json), [environment](runtime-environment.json), [restore policy](restore-policy.json), and [owned process check](process-check.json) support the result.

The official Microsoft tarball was downloaded only after metadata matched the approved URL/version/SHA512. SHA512 matched; all 5,631 archive members were contained, with no archive links. Installation stayed under `/tmp/mudarro-approved-runtimes.6AsmDd/dotnet`. No sudo, shell profile, services, certificates, workloads or global tools changed. HOME was inherited verbatim; DOTNET_CLI_HOME and every NuGet cache stayed private.

The [console sample](samples/console/Sample.csproj) has no PackageReference. [global.json](samples/console/global.json) pins the SDK with rollForward disabled; [NuGet.Config](samples/console/NuGet.Config) clears package sources and fallback folders. Restore used that explicit file, resolved zero application libraries and zero sources, and left zero nupkg files in private caches. Build used `--no-restore --disable-build-servers -p:UseSharedCompilation=false -nodeReuse:false`.

Runtime checks include repeated scan; opt-in selection of build only; explicit console/self-test/failure argv in [mudarro.json](samples/console/mudarro.json); generation twice with identical wrapper hashes; existing configuration preservation; preview/run/wrappers; repeat builds; real compile failure followed by restored source and successful recovery; deterministic multiple-project choices; invalid XML diagnostics; and exclusion before reading an oversized manifest. The self-test verifies `21 * 2 == 42`. Controlled console failure exits7 directly and1 through Mudarro/wrapper, as recorded. Scan inferred no start/framework.

An initial collector assertion expected invalid XML to produce no services, but Mudarro intentionally produces its generic custom fallback plus warning and no C# commands. The collector expectation was corrected; original attempt evidence remains privately in /tmp. An earlier private-HOME run was superseded by the complete final31-command run preserving inherited HOME; only final records are exported. No production fix was needed.

## Reproduction

Use the verified private SDK, an approved Mudarro build and the five source/config files under `samples/console`; regenerate wrappers with `mudarro generate --root <copy>`. Set the SDK/NuGet environment from `runtime-environment.json`, preserving inherited HOME. Perform `dotnet restore ./Sample.csproj --configfile ./NuGet.Config --disable-parallel --verbosity minimal -p:NuGetAudit=false`, then `mudarro run --root <copy> app-csharp:build`, `app-csharp:console`, and `app-csharp:self-test`. `app-csharp:controlled-failure` intentionally returns1; direct `dotnet bin/Debug/net10.0/Sample.dll --controlled-failure` returns7.

[runtime.py](runtime.py) is the exact executed driver with fresh private run directories. It does not download/install tools. Its fixed private SDK/bin paths must already exist. Invalid/excluded/multiple cases are static checks; only the console is compiled. SDK archive, binaries, obj/bin, caches, large negative fixture, raw intermediate logs and personal content are excluded from this repository evidence. Root integrator captures genuine menu/CLI media separately.
