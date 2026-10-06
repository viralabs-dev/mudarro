# Official checksum audit

Metadata only; no distribution, code archive, gem or installer downloaded. No installs or executions.

| Tool | File | SHA256/status |
|---|---|---|
| rust 1.99.0 | rust-1.99.0-x86_64-unknown-linux-gnu.tar.xz | 891c6366d7100feda0bca4c03ce63f3c9ac827cbebbc283e7433061d42c6a376 |
| postgres 18.6 | postgresql-18.6.tar.bz2 | 555610c24d53e4316da5b7d3fc25c279d96856d5e0e23ee308c328c5fa881d9f |
| bundler 4.0.22 | bundler-4.0.22.gem | d8d5ec84c8555e0af71db63ed7aee4d1a8fb839ec46d84212d61979242a5d75a |
| bison 3.8.2 | bison-3.8.2.tar.xz | 9bba0214ccf7f1079c5d59210045227bcf619519840ebfa80cd3849cff5a5bf2 |
| flex 2.6.4 | flex-2.6.4.tar.gz | blocked-no-official-digest |
| libyaml 0.2.5 | yaml-0.2.5.tar.gz | blocked-no-official-digest |

## Sources and limitations

- rust: https://static.rust-lang.org/dist/channel-rust-1.99.0.toml. Publisher metadata/checksum saved locally.
- postgres: https://ftp.postgresql.org/pub/source/v18.6/postgresql-18.6.tar.bz2.sha256. Publisher metadata/checksum saved locally.
- bundler: https://rubygems.org/api/v2/rubygems/bundler/versions/4.0.22.json. Publisher metadata/checksum saved locally.
- bison: https://lists.gnu.org/archive/html/bug-bison/2021-09/msg00056.html. Direct retrieval TLS EOF; full official announcement checksum visible via web indexed result; encoding converted locally.
- flex: https://api.github.com/repos/westes/flex/releases/tags/v2.6.4. Official GitHub release metadata digest is null; no publisher SHA256 established. Do not treat third-party checksums as official.
- libyaml: https://api.github.com/repos/yaml/libyaml/releases/tags/0.2.5. Official GitHub release metadata digest is null; no publisher SHA256 established. Do not treat third-party checksums as official.

Bundler requires Ruby >= 3.2.0 and RubyGems >= 3.4.1; no runtime gem dependencies in publisher metadata. Flex publishes detached signatures but verifying those requires its trusted signing key and a later authorized artifact download. libyaml metadata has no signature asset. Checksums identify future bytes; they do not attest an unperformed download.
