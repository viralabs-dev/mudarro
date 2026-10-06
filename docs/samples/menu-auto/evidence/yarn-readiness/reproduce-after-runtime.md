# MUD-047 — reprodução futura, NÃO executada

Yarn/yarnpkg não encontrado no PATH ou candidatos consultados; nenhum runtime Yarn em locais conhecidos do cache Corepack. Corepack 0.34.6 existe, mas não foi usado para baixar, preparar ou ativar Yarn. Instalação/ativação/download exige autorização separada com versão e destino definidos. Não assumir Yarn Classic e Yarn Modern equivalentes.

A fixture real em `fixture/` contém package.json com manager explícito yarn e script Node sem dependências externas. Nenhum yarn.lock fictício foi criado. Scan/init/generate do binário real passaram; tentativa de ação terminou exit 1 por executável ausente. `runtime-marker.txt` não existe: execução Yarn NÃO aprovada.

Após disponibilizar Yarn autorizado em caminho isolado, confirmar a versão e usar somente PATH desta execução (sem corepack enable/configuração global):

```bash
export PATH="/CAMINHO_APROVADO_COM_YARN:/tmp/mudarro-validation-next:$PATH"
yarn --version
/tmp/mudarro-validation-next/mudarro run app-javascript:test --root /tmp/mudarro-yarn-readiness/fixture
bash /tmp/mudarro-yarn-readiness/fixture/.mudarro/scripts/app-javascript/test.sh
cat /tmp/mudarro-yarn-readiness/fixture/runtime-marker.txt
```

Aceite: ambos os comandos devem terminar 0, imprimir a mensagem de probe e criar marker contendo exatamente o cwd da fixture. Comparar hashes package.json/config/launcher antes/depois; preservar logs e versão/hash do runtime. Nenhum teste de dependência/linking/install é coberto por isso.

Se Yarn exigir installation state, parar e registrar a ação necessária; não executar install, Corepack download ou lock generation sem autorização. Para Yarn Modern, declarar versão compatível em uma nova fixture pode ser necessário; não alterar evidência histórica desta tentativa ausente.
