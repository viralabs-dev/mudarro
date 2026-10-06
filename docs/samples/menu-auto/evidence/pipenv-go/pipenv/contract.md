MUD-034 proposta de contrato antes da implementação

Runtime pipenv não encontrado PATH, ~/.local/bin ou Linuxbrew/bin/opt consultados. Não instalar. Separar aprovação de detecção/argv da aprovação runtime.

Pipfile TOML válido escolhe manager pipenv quando não há evidência uv/Poetry conflitante. Pipfile.lock JSON válido apoia escolha quando junto ao Pipfile. Lock isolado não deve inventar projeto utilizável: avisar/pendência; sem install/start automático. Arquivo malformed torna diagnóstico explícito, sem fallback silencioso pip. Não validar hash de resolução lock vs Pipfile sem implementar algoritmo oficial.

Install: Pipfile+lock -> pipenv sync (lock-preserving); só Pipfile -> pipenv install (ação proposta, não executada por scan). Incluir dev? padrão instalar conjunto normal, configuração explícita para --dev; não decidir automaticamente.

Django: pipenv run python manage.py runserver/test, sem inferir .venv local. Não criar ação venv python3 -m venv para Pipenv. Manager ambíguo não pode inferir prefixo .venv/python nem comandos start/test/install; pendência acionável.

Pipfile [scripts] strings não vazias são sugestões pipenv run <original-key>, sem interpretar/splitar seu texto. Scripts custom continuam opt-in --select/all/interactive. Nomes inválidos ou colisões de sanitização ':'/'.' -> '-' ficam pendentes, não escolher arbitrariamente. Formas inline call podem ser suportadas como sugestão pela chave, pois Pipenv interpreta; se não suportadas, warning explícito, sem inferir execução. Não precisa decidir usuário para esse recorte conservador.

Conflitos Pipfile vs uv.lock/poetry.lock requerem manager explícito; não selecionar ferramenta automaticamente. pyproject scripts pertencentes a outros managers também não devem escapar da ambiguidade em sugestões executáveis.

Cobertura pretendida: Pipfile/lock, malformedfiles/scripts, Djangoargv, conflito, scriptselecionado/colisão, scansemexec/semwrite, initpreserva e generateidempotente. Runtime indisponível é limitação separada, mocks não equivalem à execução Pipenv.

Fontes: https://pipenv.pypa.io/en/latest/pipfile.html e https://pipenv.pypa.io/en/latest/scripts.html
