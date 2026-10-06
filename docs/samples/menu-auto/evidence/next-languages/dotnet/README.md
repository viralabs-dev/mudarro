# C#/.NET static contract (MUD-042)

Internal console and multiple-project samples have no PackageReference. Project XML is bounded and validated without DTDs/entities; MSBuild properties, imports and targets are never evaluated by detection. Explicit opt-in `dotnet build ./<project>` and `dotnet test ./<project>` suggestions preserve file argv and service cwd. Multiple projects use deterministic sorted `build-project-N` / `test-project-N` choices; the evidence field identifies each file. Prefix `./` prevents option injection for unusual filenames.

Solutions are evidence only: legacy .sln syntax and all solution graphs are unresolved, while .slnx checks XML shape only. Projects in other directories are detected independently; no referenced project is read or followed from a solution. No ASP.NET, startup, run command or aggregate selection inferred.

Static gates cover malformed XML/DTD/entities, bounded regular/symlink reads, unavailable manifests, no-execution traps, ambiguous choices, safe unusual names, exclusions, selection, idempotent generation and existing configuration preservation.

No .NET SDK/build/test runtime executed. Restore may consult feeds even without application dependencies; SDK, private caches, network and any test packages need separate approval. `dotnet test` on these samples is not evidence that a test runner executed tests.

Reproduce with `go test ./test/internal/mudarro -run TestDotNet -race -count=1` under the repository's approved Go environment.

## Authorized runtime checkpoint

The earlier static-only checkpoint above is historical. SDK10.0.401 private console build and explicit standard-library self-test now passed31 genuine commands; no dotnet test/framework execution. See [runtime evidence](runtime/README.md).
