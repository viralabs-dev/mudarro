# Personal Git attribution correction — MUD-056

Daniel explicitly authorized rewriting the eight published commits and updating only main with force-with-lease. The former global Git email mapped author/committer to daneesouza despite gh authentication as daneiel. Repository-local name/email now use the confirmed personal profile; global settings were not changed.

All eight corrected commits are published; authenticated GitHub API confirms both author.login and committer.login **daneiel** for every commit. Base049a405 was not rewritten. Trees, author/committer dates and message bodies are identical; only the first three subjects gained [skip ci], already present in the other five. No squash, tags, backup refs, workflow/settings/protection change or CI dispatch/rerun. This note and other newer local project work are **not** part of that push.

| Previous SHA | Corrected SHA | Tree |
|---|---|---|
| `d04b75b22ac218ae201f6b63857a83952ad460a5` | `f5b43b60d2477f9cdc30a48401ea5471da3fcc93` | Identical |
| `7e0562c6def635071288e52ef265c20863b5ce26` | `89522db7418976993feca7661630f3a9816969b5` | Identical |
| `bade86ef9aa087f88b783abe1ad77dbb95e10418` | `51053d736dc98bc91bdf22e55a8c8217804f57e3` | Identical |
| `9869dc6302d4e131a4033e5a080f93fcb5c7449b` | `849e0e0976f4e7cc37eee9fb4bbdef12e5e50bb4` | Identical |
| `bfbc80587bc56ca2f69cdad1f47ea3e273ad2d53` | `2406bfb71ba07a35e6f29e7bacd7a13152bcd7da` | Identical |
| `b78f2b4866313511a5903129bf5d61e80deffcbb` | `d82daa19ea4a0bcf40b26330187d895ef51c3688` | Identical |
| `246048a1d04cf0de969b4b630bf6367273472cdb` | `33218af1ef0c594a66c18d49e89163b83474dd88` | Identical |
| `56dadbaf62e3515e4bb31bb3c3ff2a583cd642ed` | `b887f64ebdfba400e312ce821122283693d868d4` | Identical |

[Full map](commit-map.json), [API attribution/tree/first Actions check](github-verification-first.json), [push with explicit expected SHA](push-proof.txt).

Backups remain outside the repository at /tmp/mudarro-identity-repair: history.bundle, workingtree.tar, index-before, config-before, original status/diff and restored verification copies. Bundle verification passed; archive restoration was tested against488 entries by bytes, modes and symlink targets. A generic filtered extraction initially removed group-write in the temporary verification copy; validated own-archive extraction preserved modes and passed. No original data or backup was deleted.

[Backup proof](backup-proof.json), [post-push integrity](workingtree-after-proof.json): original workingtree/index/status exactly preserved. Local main ref was changed directly only after successful push because old/new tip trees are identical; no stash/reset/checkout was necessary. Additional local documentation was written afterward and remains uncommitted.

Recover safely into a **new** directory rather than overwriting current work:

```bash
git clone /tmp/mudarro-identity-repair/history.bundle /tmp/mudarro-recovery-copy
# Validate archive member paths against workingtree-before.json before extracting.
# Extract workingtree.tar into the new recovery directory with original modes.
```

Other clones now have divergent history and need deliberate alignment preserving their local changes; no commands were executed against other checkouts. Old SHA citations in historical logs/CI runs remain valid historical references; use this map for their corrected equivalents. Absence of new Actions runs is an observation, not CI approval. Native macOS and existing physical acceptance limits remain unchanged.

Final repeated Actions read: 2026-10-06T14:08:30.371804+00:00; all eight SHAs total_count=0/runs=[]. [Final read](actions-final.json).
