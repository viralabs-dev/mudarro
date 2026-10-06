# Provider privado — dependências fechadas

Pesquisa somente; nenhum pacote baixado/instalado. PyPI oficial confirma podman-compose1.6.0 exige Python>=3.9 e somente python-dotenv/PyYAML no runtime. Fixados python-dotenv1.2.4 (>=3.10) e PyYAML6.0.3 (>=3.8); Python existente3.14.7 satisfaz. Não solicitar extras devel/cli: click é somente extraCLI de dotenv; ferramentas devel não entram.

requirements-hashed.txt contém exatamente três wheels com hashes oficiais. PyYAML selecionado cp314-cp314 GNUlinuxx64, não cp314t/free-threaded. closure.json guarda URLs, filenames e fontes. Antes execução verificar ABI Python3.14normal+glibc/x64; mudança de plataforma exige novo hash/wheel oficial.

Proposta NÃO EXECUTADA após aprovação: Pythonexistente -mvenv <private>/provider; pip --isolated --cache-dir <private>/cache download --only-binary=:all: --require-hashes -r requirements-hashed.txt --dest <private>/wheelhouse --index-url https://pypi.org/simple; verificar filenames/hash; pip --isolated --cache-dir <private>/cache install --no-index --find-links <private>/wheelhouse --require-hashes -r requirements-hashed.txt. Não global/home/sudo/shellRC. Arquivar pip-versão e resolvermanifest. Pip/bootstrap do venv deve ser stdlib ensurepip existente, não baixar upgrade implicitamente.

Provider disponível não comprova enginePodman nem rootless/rede/storage; instalaçãoengine continua decisão distinta e não autorizada. Sem execução provider contra host neste plano.
