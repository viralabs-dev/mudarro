from pathlib import Path
import json,shutil,hashlib
r=Path('/home/danielsouza/dev/viralabs/mudarro');d=Path('/tmp/mudarro-parser-checkpoint');q=Path('/tmp/mudarro-mud038-parser');p=r/'docs/samples/menu-auto/evidence/go-work-parser';p.mkdir(exist_ok=True)
s=json.loads((d/'summary.json').read_text());s['gates']='mod verify/race/vet/Linux build/Darwin arm64 cross-build/smoke/menu/config/resize PASS';(d/'summary.json').write_text(json.dumps(s,indent=2)+'\n')
for f in d.iterdir():
 if f.is_file() and f.suffix in ['.txt','.json','.jsonl','.out','.cast','.py']:shutil.copy2(f,p/f.name)
for sub in ['runtime-recheck','resize']:shutil.copytree(d/sub,p/sub,dirs_exist_ok=True)
for name in ['dependency-download.json','dependency-download.stderr','installer-final.txt','release-local-summary.json','release-local.txt']:shutil.copy2(q/name,p/name)
for sub in ['tests','runtime']:
 out=p/sub;out.mkdir(exist_ok=True)
 for f in (q/sub).iterdir():
  if f.is_file() and f.suffix in ['.txt','.json','.jsonl','.log','.md','.cast','.gif','.png','.py']:shutil.copy2(f,out/f.name)
sample=p/'sample';source=q/'runtime/valid'
for f in source.rglob('*'):
 if not f.is_file() or f.is_symlink():continue
 rel=f.relative_to(source)
 if 'run' in rel.parts or f.suffix not in ['.go','.mod','.work','.json','.sh']:continue
 dest=sample/rel;dest.parent.mkdir(parents=True,exist_ok=True);shutil.copy2(f,dest)
checks={str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest() for f in p.rglob('*') if f.is_file() and f.name not in ['SHA256SUMS.json','README.md']};(p/'SHA256SUMS.json').write_text(json.dumps(checks,indent=2,sort_keys=True)+'\n')
(p/'README.md').write_text(f'''# Official Go workspace parser — integrated local evidence

**{s['top_level_tests']['pass']} race PASS ({s['substantive_groups']} substantive groups + three isolated helpers), 0 FAIL/SKIP; {s['covered_statements']}/{s['total_statements']} = {s['coverage_percent']}%.** [Summary](summary.json), [test events](tests.jsonl), [coverage](coverage.out). Module verification, vet, Linux build, Darwin arm64 cross-build, smoke, menu/configuration and six genuine PTY resize variants passed. Native macOS was not executed. Final Linux binary SHA256 `{s['linux_binary_sha256']}`; base b887f64; new work remains local, uncommitted and unpublished.

## Real flows and test matrix

| Flow | Actual result / evidence |
| --- | --- |
| Official parser; quoted paths/comments/use blocks; version/toolchain/replace metadata | Nine parser groups /14 subcases pass with fresh race run; [summary](tests/summary.json), [fresh output](tests/after-followup-fresh.txt) |
| Missing/invalid member go.mod, missing declaration, syntax, outside-root relative/absolute paths, symlinks, duplicate directory/module names | Precise warnings preserve independent services; baselines and negative results retained in [tests](tests/README.md) |
| Regular-file4MiB bound, FIFO, scan excludes/ignored directories before member reads | Bounded/static checks pass; excluded entries explicitly unvalidated; isolated FIFO helpers joined without watchdog timeout |
| Scan remains offline without executing Go | Trap-executable regression passes; module/toolchain/replacement resolution remains a runtime responsibility |
| Real Go workspace import, inherited context and explicit off | [21 real checks](runtime/after-summary.json) pass, including expected import failure with off; actual wrapper and supervisor run, generated hashes/manual config preserved, owned service stopped |
| Final binary recheck | [Real Go scan/import/test](go-recheck.json) and [all five installed runtime scans/quality actions](runtime-recheck/summary.json) pass |
| Distribution and installer | [Four actual cross-built archives](release-local-summary.json) contain notices/licenses; [installer fixture tests](installer-final.txt) pass legacy/new/repeat/partial/corrupt/symlink/preservation cases |

Members are syntactically valid root-contained modules, not a promise that the workspace executes. Distinct directories declaring the same module remain listed with an explicit conflict warning; Go rejects that workspace. Canonical duplicate paths are warned/deduplicated. Excluded members appear separately and are not read/validated. No JS/Go inheritance, aggregate orchestration, packageManager provisioning, imported-module graph, replacement target or toolchain compatibility is inferred.

## Genuine recording and saved sample

![Actual Go workspace execution](runtime/workspace-menu.png)

![Actual Go workspace menu GIF](runtime/workspace-menu.gif)

[Original cast](runtime/workspace-menu.cast) and [PTY bytes](runtime/workspace-menu.pty.txt) record real execution with `REAL_WORKSPACE_IMPORT_OK`. The GIF was decoded completely:15 frames,860×647; PNG is frame7 composed from that recording and visually inspected. [Validation](runtime/media-validation.json). Headless PTY input was injected into the real application; these are terminal recording frames, not graphical desktop screenshots. No browser was used. Original recording binary SHA256 `0e3e43f4b8cff5218b4e174c0a9d683155da3490994c6fbe698e2dbfefcad193` predates the last exclusion/conflict diagnostic refinements; final-bin rechecks above verify actual unchanged positive runtime paths without relabelling the recording.

[Saved real workspace sample](sample/go.work) includes app space, library, independent module, nested module, explicit config and generated menu/wrappers. Binary, caches, transient supervisor state and personal content are excluded. Reproduce in a fresh copy; never overwrite the source sample:

```bash
# From the repository root, using an existing Go SDK on PATH:
go build -o /tmp/mudarro-workspace-demo ./cmd/mudarro
fixture_dir=$(mktemp -d /tmp/mudarro-workspace-demo.XXXXXX)
cp -R docs/samples/menu-auto/evidence/go-work-parser/sample/. "$fixture_dir/"
GOPROXY=off GOTOOLCHAIN=local GOFLAGS=-buildvcs=false /tmp/mudarro-workspace-demo scan --root "$fixture_dir" --json
GOPROXY=off GOTOOLCHAIN=local GOFLAGS=-buildvcs=false /tmp/mudarro-workspace-demo run app-space-go:probe --root "$fixture_dir"
GOPROXY=off GOTOOLCHAIN=local GOFLAGS=-buildvcs=false /tmp/mudarro-workspace-demo menu --root "$fixture_dir"
```

Select app-space-go → probe to see the cross-module marker. Existing [selection/context contract](../manager-go-decisions/README.md) preserves inherited go.work by default and supports explicit go_workspace:off. Selection promotes one Go entry to commands.start, preserving existing manual configuration. Strict SemVer version syntax remains as documented there.

## Dependency integrity and packaging

Authorized project dependency `golang.org/x/mod v0.25.0` is pinned in go.mod/go.sum, compatible with the project Go1.23 minimum; official proxy/SumDB verification is saved in [download metadata](dependency-download.json). No new Go SDK or global tool installation. [Module verification](mod-verify-retry.txt) passes. The first offline attempt lacked metadata for the existing YAML test dependency; only already-cached check.v1/x.tools .mod/.info metadata was copied read-only into the private cache to recover, without downloading/installing global tools. [Original failure](mod-verify.txt).

[Official parser source](https://github.com/golang/mod/blob/v0.25.0/modfile/work.go) · [Module tag](https://proxy.golang.org/golang.org/x/mod/@v/v0.25.0.info) · [SumDB record](https://sum.golang.org/lookup/golang.org/x/mod@v0.25.0).

[Third party notices](../../../../../THIRD_PARTY_NOTICES.md) and LICENSES preserve the three direct compiled dependency licenses. Release archives include these files; the installer preserves legacy binary-only archives and stores new notices under its target `mudarro-licenses/` directory. Tests use private fixtures/fake download transport, not a published install. Archives were built under a private /tmp source copy, not released. Binary replacement is atomic separately; binary plus notices are not a single transaction and later filesystem copy failure may leave a partial update. Cross-builds do not prove native Darwin/arm64 execution.

Per-file4MiB reading is not an aggregate memory cap; existing stat/open filesystem races remain. Real dependency/framework-heavy projects, native macOS, WSL, native PostgreSQL, local Podman and physical emulator/mouse/desktop capture remain separate environment/acceptance limits. [Artifact hashes](SHA256SUMS.json) permit integrity checking; no secrets, module cache, environments or transient process state are included.
''')
for lang in ['', 'pt-BR/']:
 for name in ['coverage.md','validation.md','support-matrix.md']:
  f=r/'docs'/lang/name;t=f.read_text();t=t.replace('## Latest code checkpoint — confirmed manager/Go decisions','## Historical checkpoint156 — confirmed manager/Go decisions').replace('## Checkpoint de código mais recente — decisões managers/Go','## Checkpoint histórico156 — decisões managers/Go');pos=t.find('\n## ');link='../samples' if lang else 'samples';title='## Latest integrated checkpoint — official Go workspace parser' if not lang else '## Checkpoint integrado mais recente — parser oficial Go workspace';para='165 race PASS (162 substantive groups + three helpers), 0 FAIL/SKIP; **2199/2746 = 80.08%**. Official go.work/go.mod parsing, membership/conflict/exclusion diagnostics, inherited context and explicit isolation pass. All final local gates and genuine Go/five installed runtime rechecks pass. Four cross-built archives and installer fixtures preserve licenses. Native macOS and physical graphical acceptance remain unexecuted.' if not lang else '165 PASS race (162 grupos substantivos + três helpers), 0 falhas/skips; **2199/2746 = 80,08%**. Parser oficial go.work/go.mod, diagnósticos de membros/conflitos/exclusões, contexto herdado e isolamento explícito passaram. Gates finais e rechecks reais Go/cinco runtimes instalados passaram. Quatro archives cross-build e fixtures do instalador preservam licenças. macOS nativo e aceite gráfico físico não executados.';t=t[:pos]+'\n'+title+'\n\n'+para+f' [Evidence / reprodução]({link}/menu-auto/evidence/go-work-parser/README.md).\n'+t[pos:]
  t=t.replace('version text still opaque pending contract','supplied versions require exact SemVer; metadata is syntax-only, with no hash verification/provisioning').replace('go.work emits limitation warning, membership/context/selection unresolved','official go.work/go.mod parser validates bounded root-contained members; conflicts/exclusions warned; selected entry promotes start; inherited context default, explicit off available').replace('go.work tem aviso de limitação, parser de membros ainda aberto; seleção targets e go_workspace off explícitos disponíveis','parser oficial go.work/go.mod valida membros bounded dentro da raiz; conflitos/exclusões diagnosticados; seleção promove start, contexto herdado por padrão e off explícito disponível')
  f.write_text(t)
for lang in ['', 'pt-BR/']:
 f=r/'docs'/lang/'language-roadmap.md';t=f.read_text().replace('Go multiple binaries/go.work','Go module graph/framework integrations (entry selection and official go.work membership parsing delivered)').replace('Go múltiplos bins/go.work','integrações de grafo/framework Go (seleção de entrada e parser oficial de membros go.work entregues)');t+='\n\n'+('Current Go workspace parser, explicit entry selection and inherited/off policy are delivered and validated locally; replacement/toolchain/dependency graph resolution and workspace service inheritance remain outside this implementation. ' if not lang else 'Parser Go workspace, seleção explícita de entrada e política herdado/off entregues e validados localmente; resolução de replacements/toolchain/grafo e herança de serviços workspace permanecem fora da implementação. ')+('[Evidence](samples/menu-auto/evidence/go-work-parser/README.md).\n' if not lang else '[Evidências](../samples/menu-auto/evidence/go-work-parser/README.md).\n');f.write_text(t)
for f in [r/'docs/samples/menu-auto/evidence/manager-go-decisions/README.md',r/'docs/samples/menu-auto/evidence/manager-go-decisions/go-work-parser-plan/README.md']:
 t=f.read_text();idx=t.find('\n');t=t[:idx]+'\n\n**Historical checkpoint/proposal:** subsequent explicit parser authorization superseded the earlier dependency constraint. The official pinned parser is now implemented and verified; [current integrated result]('+('../' if f.parent.name=='go-work-parser-plan' else '')+'../go-work-parser/README.md). Previous baselines and approval state below are retained as historical evidence.\n'+t[idx:];f.write_text(t)
print('Integrated parser evidence, real sample, and EN/PT documentation')
