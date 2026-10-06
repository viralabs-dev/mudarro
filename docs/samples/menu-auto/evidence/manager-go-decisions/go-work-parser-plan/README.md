# Remaining MUD-038: official workspace parser proposal

**Historical checkpoint/proposal:** subsequent explicit parser authorization superseded the earlier dependency constraint. The official pinned parser is now implemented and verified; [current integrated result](../../go-work-parser/README.md). Previous baselines and approval state below are retained as historical evidence.


Not implemented or installed. The previous instruction recorded in the MUD-038 Vault activity forbids an improvised partial parser or a new dependency. The three confirmed choices settle version syntax, Go target promotion and default/isolated execution context; they do not explicitly lift that dependency constraint. Existing workspace warnings remain, so MUD-038 is not claimed complete.

Concrete proposed project dependency: `golang.org/x/mod v0.25.0`, using `modfile.ParseWork` for bounded, regular-file go.work bytes and the official modfile parser for member go.mod files. This tag requires Go1.23.0, matching this project's declared minimum. No new Go SDK or global tool installation is proposed. Its module metadata declares x/tools v0.13.0 for tagged tooling; resolve the actual required graph transparently, without installing global binaries.

Expected changes: go.mod/go.sum add the pinned project dependency; offline scanning parses syntax, enumerates use directives, checks root-contained regular member modules and diagnoses missing/invalid/outside-root members. Preserve discovery of independent go.mod services and inherited runtime context. Do not infer manager/service inheritance or dependency resolution. Keep replacements/version/toolchain information declarative and disclose limits; never execute Go during scan. Static FIFO, oversized, symlink, nested-workspace, missing/outside members, syntax and valid-member fixtures must verify preservation, deterministic results and scan with a trap Go executable. Native workspaces and resulting CLI/runtime behavior need actual regression runs before completion.

Verified official checksums, copied from SumDB lookup:

- Module: `h1:n7a+ZbQKQA/Ysbyb0/6IbB1H/X41mKgbhfv7AfG/44w=`
- go.mod: `h1:IXM97Txy2VM4PJ3gI61r1YEk/gAj6zAHN3AdZt6S9Ww=`

[Official source go.mod](https://github.com/golang/mod/blob/v0.25.0/go.mod) · [SumDB lookup](https://sum.golang.org/lookup/golang.org/x/mod@v0.25.0) · [saved module metadata](go.mod) · [saved checksum response](checksum.txt).

Required approval: permit this exact project dependency and corresponding go.mod/go.sum changes despite the earlier no-new-dependency constraint. This is separate from the already approved five runtime installations. No dependency code was downloaded; only metadata/checksums were read.
