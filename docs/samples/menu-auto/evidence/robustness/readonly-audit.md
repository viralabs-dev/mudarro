# Auditoria somente leitura — MUD-035 / MUD-050 / MUD-052

Não assumidos cards, não implementado código. Base: leitura atual language_adapters.go e docs/coverage.md/pt-BR, support-matrix. Relatório não é execução de testes nem nova medição.

## MUD-035: comportamento atual e contrato proposto

Hoje packageManager não vazio vence locks e usa strings.Split(valor,"@")[0]. Nenhuma validação de nome/versão. String whitespace vira executável whitespace; @9 resulta manager vazio; qualquer nome arbitrário vira comando install/run. Múltiplos locks só produzem Pending sem packageManager. JSON packageManager não-string já falha unmarshal. Bun locks não reconhecidos; isso é MUD-033 separado, sem escolher política aqui.

Matriz proposta antes de implementar:

| Caso | Comportamento contratual a decidir/testar |
|---|---|
| Ausente ou string vazia, sem lock | Preservar npm default compatível |
| Ausente + lock único npm/pnpm/yarn | Manager correto |
| Ausente + locks conflitantes | Pending explícito, sem install/start inventados |
| npm/pnpm/yarn conhecidos com versão | Extrair nome mantendo evidência versão; definir sintaxe de versão aceita, sem download/consulta registry |
| Nome conhecido sem versão | Decidir permitir forma explícita legada ou exigir versão; testar antes de mudar compatibilidade |
| Espaços leading/trailing/interiores | Definir trim estrito ou rejeição; nunca argv executável vazio/whitespace |
| Nome desconhecido/caminho/URL/flags/separador shell | Erro/pending acionável; nenhum subprocesso; preservar serviço irmão válido |
| @version ou manager@ sem nome/versão | Rejeição clara sem comandos vazios |
| Tipo number/null/object/array | number/object/array warning parse; null política explícita (unmarshal string atual trata vazio) |
| packageManager conhecido contradiz lock único ou múltiplo | Definir precedence já existente vs exigência conflito; origem explícita e teste não silencioso |
| packageManager Bun ou bun.lock/bun.lockb | Não decidir ambiguidade Bun aqui; dependência MUD-033/contrato separado |
| Config explícito custom manager | Separar validação do manifest inferido de escape hatch commands explícitos; não bloquear comando custom de projeto por lista de adapters |

Aceite mínimo: scan determinístico/offline, não roda runtime; warnings/pending legíveis/localizados; init selection não salva config parcial em ambiguidade; generate duas vezes sem diffs; argv nominal seguro; evidência origem preservada. Runtime npm/pnpm/Yarn/Bun é estágio separado conforme ferramentas/autorização.

## MUD-050 aceites e docs a atualizar depois dos gates

Negativos supervisor: estado inexistente/stale/truncado/JSON inválido, PID supervisor inválido/morto, token/projeto/identidade divergentes, processo estrangeiro nunca sinalizado, startup encerra cedo/sem start, árvore de processos, concorrência up/down. Fixtures precisam isolamento e evitar PID pessoal; nenhuma limpeza genérica. Testes não devem substituir E2E real de lifecycle. Registrar quais ramos foram exercitados sob coverpkg e quais só em binário externo.

Atualizar docs/coverage EN/PT risco supervisor e números SOMENTE do novo perfil; docs/validation e support-matrix distinguindo testes estado/sinais versus runtime; notes MUD050/Vault pelo integrador. Antiga afirmação local_runner0% histórica não significa falta de E2E e não deve ser transportada para perfil novo.

## MUD-052 aceites e docs a atualizar depois dos gates

Falhas de leitura/manifest, criação/rename/tempwrite/close, destino arquivo-diretório, permissão se teste viável sem confiar em root, symlink e preflight inválido devem preservar conteúdo/manifest do usuário e não criar arquivos parciais. OSExecutor Output stdout+erro/exitstatus vazio, execução ausente, args vazios; preservar Unwrap/errors.As e saída externa. Não criar abstrações somente para cobertura. Documentar eventual comportamento de tempfile residual/cleanup, sem afirmar transação atômica total se escrita multiarquivo continua sequencial após preflight.

Atualizar coverage EN/PT atomicWrite/Output/SafePath/ReadManifest percentuais novos com evidência, matriz de negativos e limitações; commands de coverprofile devem apontar perfil final existente. README só se comportamento visível mudar; ADR só para desenho real alterado.

## Drift atual constatado

Coverage EN/PT acumula vários cabeçalhos 'Current' históricos (MUD01741/72.53, UI74/76.38). Reclassificar cada etapa por data/checkpoint e um único resumo final após gates. Números de risco atomicWrite53.3/Output40 podem ser históricos; comparar novo functions.txt antes atualizar. Shell032 tem denominador maior; não interpretar aumento/diminuição como requisitos. A referência de comando atual aponta project-name-coverage.out (MUD022), não perfil mais recente; ajustar só após entrega de novo perfil. Preservar histórico bruto e distinção CI macOS falhou vs fix Linux/crosscompile aprovado, sem executar CI.
