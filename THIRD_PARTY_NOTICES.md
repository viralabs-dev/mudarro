# Third-party notices

Mudarro itself is licensed under the [MIT License](LICENSE). Release archives (`.tar.gz` and the Windows `.zip`) include `LICENSE`, and both installers (`install.sh`, `install.ps1`) keep it at `mudarro-licenses/LICENSE`.

The Mudarro binary includes these project dependencies. Preserve their copyright notices and complete license texts when distributing binaries or source incorporating them.

| Module | Version | License text |
|---|---|---|
| golang.org/x/mod | v0.25.0 | [Go Authors BSD license](LICENSES/golang.org-x-mod.txt) |
| github.com/pelletier/go-toml/v2 | v2.2.3 | [go-toml license](LICENSES/go-toml-v2.txt) |
| gopkg.in/yaml.v3 | v3.0.1 | [YAML licenses](LICENSES/gopkg.in-yaml.v3.txt) |

The Go workspace parser uses the official x/mod modfile API. No dependency is installed as a global executable. Release archives include this notice and the license directory. The installer preserves them under its chosen install directory in `mudarro-licenses/`; runtime validation artifacts remain local until separately authorized for publication.
