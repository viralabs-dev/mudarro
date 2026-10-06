English · [Português (Brasil)](pt-BR/language-decision.md)

# ADR-MUD-014 — reassess Go

Status: assessment complete, keep Go+Bash; 2026-10-06. No migration or rewrite of original ADR-MUD-001. Current requirements: distributable CLI, embedded parsers, offline scan/templates, structured argv, supervision and tests. Existing operational evidence supports retention, not universal language superiority. No alternative prototypes or comparative benchmarks were performed.

| Criterion | Go | Alternatives |
|---|---|---|
| Distribution | Linux/macOS binary with embedded parsers/runtime | Rust binary also viable; Python/Node packaging exists; Bash parsing adds external tools |
| Offline scan | Implemented/tested | Any language can work offline; no alternative prototype |
| Processes/signals | os/exec and tested Unix supervisor | Rust/Python/Node APIs possible; rewriting adds cost; Bash quoting/state/errors become harder at this scale |
| Portability | Bash/Unix dependencies; WSL runtime untested | Changing language does not remove current Unix assumptions |
| Maintenance | go test/race/vet and small interfaces | Rust no-GC useful for specific constraints; team ecosystem can favor Python/Node |
| Terminal UI | Small ANSI/ASCII layer; optional Go libraries | Other ecosystems viable, colors do not justify migration |

Cost/maintenance comparisons are inferences from current requirements/code, not universal productivity measurements. Reconsider Go with strict measured binary/memory/latency constraints, no-GC/low-level requirements, predominantly web/plugin ecosystem, or a team exclusively skilled elsewhere. Native Windows requires new launcher/supervisor whatever language is chosen. Current critical path includes filesystem/external tools; no CPU/GC bottleneck demonstrated.

Keep Go+Bash and consumer-owned interfaces. Application runtimes (Node/Python/database/Docker/kubectl) remain action requirements, not scanner requirements. Darwin arm64 cross-build passed; it is not macOS runtime execution. Revisit through a recorded requirement, size/memory/scan/startup baseline, minimal alternative prototype and target-platform operational tests.

Primary sources: [Go FAQ](https://go.dev/doc/faq), [os/exec](https://pkg.go.dev/os/exec), [Rust Hello World](https://doc.rust-lang.org/book/ch01-02-hello-world.html), [Python zipapp](https://docs.python.org/3/library/zipapp.html), [Node single executable applications](https://nodejs.org/api/single-executable-applications.html). Go binaries contain runtime/GC; Python zipapp requires an interpreter and limits native modules; Node also supports single-executable packaging. See [matrix](support-matrix.md), [architecture](architecture.md), [visuals](terminal-visual.md).
