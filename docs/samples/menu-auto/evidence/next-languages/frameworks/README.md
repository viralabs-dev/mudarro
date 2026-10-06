# MUD-005: declared framework metadata, explicit commands

Static scan recognizes only already supported Node names (`next`, `@nestjs/core`, `vite`, `express`) and approved Python declared names (`Flask`, `FastAPI`, `django`, with `pytest` as a tool). Evidence comes from bounded contained regular manifests: package.json dependencies/devDependencies/optionalDependencies; pyproject.toml PEP621 dependencies/optional groups and Poetry dependency/group tables; requirements.txt direct named requirement lines. Includes, command options, imported source, installed packages and unknown framework names are not resolved. Multiple declared frameworks leave the singular framework empty and request explicit choice. No commands or startup are added by this enrichment.

## Executed evidence

Six Go test groups and eleven named subcases pass with `-race`; tests cover approved metadata forms, no inferred start/test, ambiguity, unknown/empty names, wrong types, exclusions, symlink, oversize and hostile import/executable fixtures and preservation of explicitly declared Python entrypoint suggestions. [Structured test events](focused-tests.jsonl).

[12 actual CLI commands](runtime-results.json): Python/Node scan, generate twice with all files/config byte-identical, preview, explicit user-configured check, and controlled failure for each. Python uses only unittest; Node uses node:assert/strict. These are standard-library command/runtime checks, **not Flask/Express or other installed-framework runtime validation**. No install or server was run. Runtime binary was built locally from the current shared worktree; root integration will provide final aggregate gates.

Internal samples contain intentionally declared but uninstalled Flask/pytest and Express dependencies. Reproduce: build `go build -o /tmp/mudarro-mud005-repro ./cmd/mudarro`, then run `/tmp/mudarro-mud005-repro scan --json --root <sample>`; generate twice; preview/run `owned-python:check` or `owned-node:check`. Python sample requires existing python3; Node sample requires existing node. Do not run install to validate metadata. Configured checks are explicit argv with no framework import.

## Limits

Legacy manage.py Django commands and existing Node script handling predate this change; their behavior is not invented by this hook. Pipfile/Pipfile.lock metadata enrichment, transitive dependencies, requirements includes, installed availability, HTTP healthchecks, dependency resolution, DSL/import evaluation and inferred application modules are outside this minimum. Framework runtime awaits separately authorized packages.
