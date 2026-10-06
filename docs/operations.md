English · [Português (Brasil)](pt-BR/operations.md)

# Operations

1. Run scan and inspect evidence; `scan --json` supports automation.
2. Initialize config and select additional scripts.
3. Resolve infrastructure, application entrypoint, database and migration-tool choices.
4. Inspect `generate --dry-run`, generate and run doctor.
5. Explicitly invoke `run service:install`, venv or db-install when needed. Install Docker/Podman/kubectl through the system's own procedure.
6. Open menu.sh or use `run service:action`.

Doctor checks executables and Compose provider; it does not prove application health, database connection, credentials or every runtime module.

## Files and regeneration

Preexisting files remain outside generator ownership. Generated files are checked against the manifest before writes; manual changes produce conflicts. Preserve edits or restore generated content deliberately. Do not delete the manifest to force an update. Stale wrappers are not removed automatically; removed actions can no longer run through them. Clean up deliberately.

## Local execution

The supervisor has a separate process group and random token. Identity and project are checked before signaling, without searching application process names. Logs: `.mudarro/run/<service>/output.log`. Down signals the supervisor group; applications must stay foreground without daemonizing/creating independent sessions. Startup confirms initial survival, not application health.

## Containers and Kubernetes

Compose uses the selected file. Docker invokes docker compose; Podman invokes podman-compose to avoid accidental Docker-provider selection. Down does not request volume deletion. Standalone Dockerfile supports image build and owned-container create/resume; configure image/ports explicitly. Shared monorepo Compose stacks should have one owner.

Kubernetes always uses configured context/namespace; manifests must describe this application's resources. Stop scales Deployments/StatefulSets without deleting storage. Helm relies on conventional release labels for workload selection. No cluster provisioning or production operation is provided.

## Distribution

`scripts/release.sh vX.Y.Z` builds four CGO-free Linux/macOS binaries and tar.gz/checksums.txt. Tags trigger release workflow. WSL uses Linux binaries. Installer accepts MUDARRO_VERSION/MUDARRO_INSTALL_DIR/MUDARRO_REPOSITORY; no sudo or shell changes. HTTPS downloads are hash-verified before atomic replacement.

Template references: [Prisma 7 configuration](https://www.prisma.io/docs/orm/v7/reference/prisma-config-reference), [Alembic tutorial](https://alembic.sqlalchemy.org/en/latest/tutorial.html). Measured combinations: [validation](validation.md).
