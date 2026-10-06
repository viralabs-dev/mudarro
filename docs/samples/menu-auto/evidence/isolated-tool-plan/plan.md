# MUD-058 — instalação isolada proposta, não executada

Leitura de metadados oficiais em 2026-10-06. Nenhum binário/pacote baixado, nenhum installer executado. Raiz pedirá autorização em lote. Downloads futuros executam código de terceiros no usuário atual; isolamento de diretórios não é sandbox. Sem sudo, global install, shellRC, troca de credenciais ou settings de segurança. Recursos só numa raiz /tmp nova criada com mktemp, chmod700, caches locais e env por subprocesso.

## Pins verificados e origem

| Ferramenta | Pin | Compatibilidade declarada |
|---|---|---|
| Pipenv | 2026.8.0 | Python>=3.10 |
| Poetry | 2.5.1 | Python>=3.10,<4 |
| UV | 0.12.23 | Python>=3.8; wheel Linuxx86_64 manylinux2.17 |
| Yarn Classic | 1.22.22 | Node>=4 |
| Corepack existente, opção Modern | 0.34.6 | Node^20.10.0/^22.11.0/>=24; Node24.15.0 local verificado |
| Yarn Modern | 4.18.1 | Node>=18.12 |

Python3.14.7 e Node26 satisfazem ranges declarados; não comprova execução bem-sucedida de dependências transitivas. summary.json contém hashes/integrity e links; *-metadata.json conserva resposta oficial. Hashes diretos PyPI: Pipenv8da889c636a8cab59a435c1366f30d252ecee465372341eb8cd4da43b055ad9d; Poetry3f7910abc667ce0c33a67e306baf017acdec1c8eeec4aa64d176e530d8338a61; UV565c6e2874dbeae86c02f3dea97255e878fec672659a73d4930c6b93fcab2fff. Esses hashes são de wheels diretos, NÃO lock de toda árvore transitiva. npm integritysha512 está em summary.json.

## Comandos propostos após autorização

```bash
runtime_root=$(mktemp -d /tmp/mudarro-runtimes.XXXXXX)
chmod 700 "$runtime_root"
mkdir -p "$runtime_root/config"
: > "$runtime_root/config/npm-user.rc"
: > "$runtime_root/config/npm-global.rc"
python_bin=/home/linuxbrew/.linuxbrew/bin/python3
# Não ativar venv nem mudar PATH global.
for tool in pipenv poetry uv; do
  "$python_bin" -m venv "$runtime_root/$tool/tool-venv"
done
"$runtime_root/pipenv/tool-venv/bin/python" -m pip --isolated --cache-dir "$runtime_root/cache/pipenv" install --only-binary=:all: --index-url https://pypi.org/simple --report "$runtime_root/pipenv/install-report.json" pipenv==2026.8.0
"$runtime_root/poetry/tool-venv/bin/python" -m pip --isolated --cache-dir "$runtime_root/cache/poetry" install --only-binary=:all: --index-url https://pypi.org/simple --report "$runtime_root/poetry/install-report.json" poetry==2.5.1
"$runtime_root/uv/tool-venv/bin/python" -m pip --isolated --cache-dir "$runtime_root/cache/uv" install --only-binary=:all: --index-url https://pypi.org/simple --report "$runtime_root/uv/install-report.json" uv==0.12.23
```

Antes da instalação efetiva, preferir pip download/metadata/report em wheelhouse isolada, conferir hashes diretos oficiais e gerar requirements travado com hashes para transitivas; então install --no-index --find-links wheelhouse --require-hashes. Não fazer fallback sdist/build se wheel incompatível: informar bloqueio. Pins superiores fixos não garantem transitivas fixas; report/freeze+hases documentam resolução escolhida. --isolated ignora config pessoal, mas PIP_CACHE_DIR também pode ser ignorado: passar --cache-dir explicitamente no comando efetivo. Assim nenhuma escrita de cache pessoal é pretendida. Use pip --cache-dir "$runtime_root/cache/<tool>" --isolated ... como forma final.

UV standalone oficial é alternativa sem venv Python; exigiria assetrelease Linuxarquitetura correta +checksum verificado e extração local. Evitar curl|sh: installer poderia alterar perfis/PATH e exige escopo adicional. Preferência proposta pelo menor número de fluxos: UVwheel PyPIvenv acima, executável compilado empacotado; standalone não instalado simultaneamente sem necessidade.

### Yarn: duas opções explicitamente diferentes

Classic1.22.22 via npm local encaixa scripts/install atuais e menor escopo, mas não valida YarnModern4/workspaces/PnP. Modern4.18.1 via Corepack0.34.6 já existente com cache privado valida linha atual, com mecanismo/artefatos/cache adicionais e sem garantir compatibilidade templateMudarro. Escolher Classic para primeiro gate limitado e Modern separado, OU ambos se usuário autorizar o escopo comparativo; não chamar Classic de validação de YarnModern.

```bash
# Opção Classic, sem -g; prefix/cache privados, scripts de instalação desativados.
npm --userconfig "$runtime_root/config/npm-user.rc" --globalconfig "$runtime_root/config/npm-global.rc" --cache "$runtime_root/cache/npm-classic" install --prefix "$runtime_root/yarn-classic" --ignore-scripts --no-audit --no-fund yarn@1.22.22
"$runtime_root/yarn-classic/node_modules/.bin/yarn" --version
# Opção Modern: reutilizar Corepack EXISTENTE, sem instalar quinta ferramenta.
node_existing=/home/danielsouza/.nvm/versions/node/v24.15.0/bin/node
corepack_existing=/home/danielsouza/.nvm/versions/node/v24.15.0/lib/node_modules/corepack/dist/corepack.js
COREPACK_HOME="$runtime_root/cache/corepack" COREPACK_ENABLE_PROJECT_SPEC=0 "$node_existing" "$corepack_existing" yarn@4.18.1 --version
# Se Mudarro precisar resolver "yarn" via PATH, criar shim SOMENTE em
# "$runtime_root/yarn-modern/bin/yarn" invocando os caminhos absolutos acima;
# PATH scoped ao subprocesso. Não corepack enable ou npm -g.
```

npm pode resolver dependências/downloads registry; --ignore-scripts evita lifecycle install, mas runtime posterior executará código instalado. ModernCorepack baixa Yarnpin no cache privado; conferir integridade por metadados e lockfile local antes reproduzir. Não corepack enable/use global, não modificar packageManager do repo real. Comandos ainda não executados: confirmar CLI flags/arquivos do pacote instalado e registrar falhas em vez de prometer sucesso.

## Runtime após disponibilizar ferramentas

Fixtures novas sem dependências, copies fora do repo real. UV: cache/interpretersystem explícitos, UV_PYTHON_DOWNLOADS=never, --offline; gerar lock real, sync/test/run com arquivos locais. Poetry: POETRY_CACHE_DIR privado, POETRY_VIRTUALENVS_IN_PROJECT=true, keyringfalse, env use Python existente; package-mode=false para evitar build backend desnecessário. Pipenv: PIPENV_VENV_IN_PROJECT=1, PIPENV_DONT_LOAD_ENV=1, PIPENV_IGNORE_VIRTUALENVS=1, PIPENV_PYTHON explícito, cache privado; não pyenv/asdf install. Yarn: cwdfixture package.jsonscriptsnode, sem deps; cacheoffline quando viável. Lockfiles devem ser gerados pela ferramenta real, sem usar sentinelas de detecção como runtime.

Gates: versões reais, manifest/lock, scan/init/generate twice, argv/menu/wrappers e ação runtime mínima; nenhum install geral do Mudarro automático. Logs separados por ferramenta; não fixtures/ports concorrentes. Falha de backend/cache/python não autoriza download extra. Variáveis de cache/config por processo, sem expor tokens pessoais. Root integra evidências/Vault; nenhum push/CI nesta preparação.

## Fontes primárias

- [Pipenv PyPI](https://pypi.org/project/pipenv/2026.8.0/) e [documentação](https://pipenv.pypa.io/).
- [Poetry PyPI](https://pypi.org/project/poetry/2.5.1/) e [instalação oficial](https://python-poetry.org/docs/#installation).
- [UV PyPI](https://pypi.org/project/uv/0.12.23/) e [instalação oficial Astral](https://docs.astral.sh/uv/getting-started/installation/).
- [Yarn instalação](https://yarnpkg.com/getting-started/install), [Classic](https://classic.yarnpkg.com/en/docs/install/), [Corepack](https://github.com/nodejs/corepack/releases/tag/v0.36.0), [Yarn4.18.1](https://github.com/yarnpkg/berry/releases/tag/@yarnpkg/cli/4.18.1).

## Revisão mínima — Corepack existente

Corepack0.34.6 encontrado em /home/danielsouza/.nvm/versions/node/v24.15.0/bin/corepack; --version passou e node --version retornou v24.15.0. Metadados oficiais engines^20.10.0/^22.11.0/>=24 permitem essa versão; Yarn4.18.1 engines>=18.12 também. Bin dist/corepack.js esperado pelo pacote; conferir existência antes execução autorizada. Cache Yarn ainda ausente segundo frente readiness; invocação Modern baixará Yarn, requer autorização apesar de Corepack já existir. Não executar chamada Yarn durante preparação.

Pedido em lote: quatro famílias Pipenv/Poetry/UV/Yarn, NÃO instalar Corepack novo. Classic1.22.22 é mínimo possível, Modern4.18.1 é opção de cobertura distinta reutilizando Corepack. Autorizar ambos canais Yarn aumenta downloads/cobertura e deve ser explícito. Metadados0.36.0 anteriores ficam históricos, não são target de instalação. Nada instalado/baixado fora de JSON de metadados.


## Harmless command review correction

A read-only `npm config get prefix` probe confirmed using /dev/null for both userconfig/globalconfig fails with double-loading config. Final proposal uses two distinct empty private config files created only under the fresh runtime root. No installation was attempted.


## Authorized execution after this proposal

User approved all five exact versions on2026-10-06. Actual installations and zero-dependency runtimes completed in a fresh private temporary root; see [verified results, samples and genuine recordings](../isolated-runtimes/README.md). Earlier “not executed” text describes this proposal before approval, not current availability. No additional tools/versions or global changes were authorized.
