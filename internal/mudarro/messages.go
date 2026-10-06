package mudarro

import "strings"

// ownedMessagePairs contains application-owned prose only. Identifiers, argv,
// configuration keys and output from external programs are not catalog entries.
// Callers must apply localization only to messages owned by Mudarro, never to
// arbitrary subprocess output or configuration values.
var ownedMessagePairs = [][2]string{
	{"declare infrastructure.image para Dockerfile", "declare infrastructure.image for Dockerfile"},
	{"declare infrastructure.mode: compose ou dockerfile", "declare infrastructure.mode: compose or dockerfile"},
	{"declare context, namespace e file Kubernetes", "declare Kubernetes context, namespace and file"},
	{"mode Kubernetes deve ser manifests, kustomize ou helm", "Kubernetes mode must be manifests, kustomize or helm"},
	{"declare infrastructure.kind (detecção ambígua)", "declare infrastructure.kind (ambiguous detection)"},
	{"declare database.kind: postgresql ou sqlite", "declare database.kind: postgresql or sqlite"},
	{"ferramenta de migration não suportada", "unsupported migration tool"},
	{"ui inválido: ", "invalid ui."},

	{"version deve ser 1", "version must be 1"},
	{"name obrigatório", "name is required"},
	{"id inválido ou duplicado: ", "invalid or duplicate id: "},
	{"diretório inexistente: ", "directory does not exist: "},
	{"porta inválida: ", "invalid port: "},
	{"ação inválida: ", "invalid action: "},
	{"defina args OU shell", "define args OR shell"},
	{"executável vazio: ", "empty executable: "},
	{"url_env inválida", "invalid url_env"},
	{"caminho relativo inválido: ", "invalid relative path: "},
	{"caminho fora do projeto: ", "path outside project: "},
	{"symlink não permitido: ", "symlink not allowed: "},
	{"manifest grande demais: ", "manifest exceeds size limit: "},
	{"comando vazio", "empty command"},
	{"DaemonSet não escala a zero; declare ação down personalizada", "DaemonSet cannot scale to zero; declare a custom down action"},
	{"nenhum workload encontrado nos recursos declarados", "no workload found in declared resources"},
	{"Múltiplos lockfiles: declare manager", "Multiple lockfiles: declare manager"},
	{" colidem no nome ", " collide on name "},
	{": declare comando explícito", ": declare an explicit command"},
	{"Declare commands.start para execução local", "Declare commands.start for local execution"},
	{"Declare commands.start; módulo Python não é inferido", "Declare commands.start; Python module is not inferred"},
	{"Declare commands.start: nenhuma entrada Go única", "Declare commands.start: no unique Go entrypoint"},
	{"dois geradores reivindicam ", "two generators claim "},
	{"; configure um único proprietário", "; configure a single owner"},
	{"manifest inválido: ", "invalid manifest: "},
	{"arquivo gerado alterado manualmente: ", "generated file manually changed: "},
	{"; preserve sua versão ou restaure antes de gerar", "; preserve your version or restore before generating"},
	{"menu interrompido", "menu interrupted"},
	{"valor ausente: ", "missing value: "},
	{"opção desconhecida: ", "unknown option: "},
	{"entrada encerrada antes de salvar configuração", "input ended before saving configuration"},
	{"defina o gerenciador antes de incluir ", "define the manager before including "},
	{"script sugerido inexistente: ", "suggested script does not exist: "},
	{"uso: mudarro run serviço:ação", "usage: mudarro run service:action"},
	{"comando desconhecido: ", "unknown command: "},
	{"informe serviço:ação em projetos com vários serviços", "provide service:action in projects with multiple services"},
	{"ação inexistente: ", "action does not exist: "},
	{"não executável", "not executable"},
	{" pendência(s); as ações independentes continuam disponíveis", " pending item(s); independent actions remain available"},
	{"operação cancelada", "operation cancelled"},
	{"ação sem comando: ", "action has no command: "},
	{" não pertence a este serviço", " does not belong to this service"},
	{"processo não encerrou; verifique os logs antes de tentar novamente", "process did not stop; check logs before retrying"},
	{"operação local desconhecida: ", "unknown local operation: "},
	{"aplicação encerrou durante a subida; consulte logs", "application exited during startup; check logs"},
	{"supervisor inválido", "invalid supervisor"},
	{"serviço inexistente", "service does not exist"},
	{"commands.start ausente", "commands.start is missing"},
	{"generate de infraestrutura suporta docker, podman ou kubernetes", "infrastructure generation supports docker, podman or kubernetes"},
	{"commands.start obrigatório para gerar imagem", "commands.start is required to generate an image"},
	{"manager JS/TS deve ser npm, pnpm ou yarn", "JS/TS manager must be npm, pnpm or yarn"},
	{"declare commands.install.args para Python pip", "declare commands.install.args for Python pip"},
	{"manager Python não suportado", "unsupported Python manager"},
	{"declare Dockerfile existente para linguagem ", "declare an existing Dockerfile for language "},
	{"geração Kubernetes usa manifests ou kustomize; Helm integra charts existentes", "Kubernetes generation uses manifests or kustomize; Helm integrates existing charts"},
	{"geração Kubernetes exige image, file (diretório), port, context e namespace", "Kubernetes generation requires image, file (directory), port, context and namespace"},
	{"database.kind deve ser postgresql ou sqlite", "database.kind must be postgresql or sqlite"},
	{"database.tool não suportada: ", "unsupported database.tool: "},
	{"Declare comandos da aplicação", "Declare application commands"},
	{"Compose detectado: declare infrastructure.kind como docker ou podman", "Compose detected: declare infrastructure.kind as docker or podman"},
	{"Múltiplas infraestruturas: declare infrastructure.kind, mode e file", "Multiple infrastructures: declare infrastructure.kind, mode and file"},
	{"Declare context e namespace Kubernetes", "Declare Kubernetes context and namespace"},
	{"Múltiplas ferramentas de banco: declare database.tool", "Multiple database tools: declare database.tool"},
	{"Declare database.kind; conexão não será executada durante detecção", "Declare database.kind; connection is not executed during detection"},
	{"informe --name com letras, números, _ ou -", "provide --name with letters, digits, _ or -"},
	{"variável de ambiente ausente: ", "missing environment variable: "},
	{"--format inválido: ", "invalid --format: "},
	{"formato de configuração não suportado: ", "unsupported configuration format: "},
	{"--format conflita com a extensão de --config", "--format conflicts with --config extension"},
	{"configuração já existe; edite o arquivo preservado", "configuration already exists; edit the preserved file"},
	{"ui inválida: ", "invalid ui: "},
	{"mudarro.yaml e mudarro.json existem; selecione uma configuração explicitamente", "both mudarro.yaml and mudarro.json exist; select a configuration explicitly"},
	{"configuração deve ser arquivo regular: ", "configuration must be a regular file: "},
	{"configuração excede 4 MiB: ", "configuration exceeds 4 MiB: "},
	{"esperado um documento JSON", "expected one JSON document"},
	{"esperado um documento YAML", "expected one YAML document"},
	{"aninhamento JSON excede 256 níveis", "JSON nesting exceeds 256 levels"},
	{"chave de objeto JSON inválida", "invalid JSON object key"},
	{"chave JSON duplicada ", "duplicate JSON key "},
	{"delimitador JSON inesperado ", "unexpected JSON delimiter "},
	{"terminal interativo indisponível", "interactive terminal unavailable"},
}

// localizedMessage translates known owned clauses while retaining interpolated
// values. Locale defaults to English; unknown text remains verbatim.
func localizedMessage(text, locale string) string {
	pairs := make([]string, 0, len(ownedMessagePairs)*2)
	for _, pair := range ownedMessagePairs {
		if locale == "pt-BR" {
			pairs = append(pairs, pair[1], pair[0])
		} else {
			pairs = append(pairs, pair[0], pair[1])
		}
	}
	return strings.NewReplacer(pairs...).Replace(text)
}

// localizedError retains the original error for errors.Is/errors.As.
type localizedError struct {
	cause   error
	message string
}

func (e localizedError) Error() string { return e.message }
func (e localizedError) Unwrap() error { return e.cause }
