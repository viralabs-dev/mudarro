# MUD045: preparação C/C++ com ferramentas existentes

PASS: 40 comandos locais reais, 20 por sample independente. GNU Make4.3 e GCC/G++13.3.0 de /usr/bin compilaram programas sem dependências externas. C imprime21; C++ imprime42. O target test executa e confere a saída exata. Repetição preservou hash/mtime dos binários e manteve testes aprovados. Fontes e configuração manual foram preservadas; geração de wrappers foi idempotente.

A receita de falha sai7; GNU Make retorna2. Target inexistente também retorna2. Mudarro relata a falha do processo e seu wrapper retorna1. É propagação esperada, não defeito de produção. Logs por família em logs/c e logs/cpp; results.json contém comandos/saídas e acceptance.json registra as verificações.

O scanner atual retorna serviço custom genérico e sugestões Make existentes. Nenhum adapter C/C++ foi implementado. mudarro.yaml declara build/test/fail manualmente; preview/wrappers demonstram somente integração explícita. Sem inferência de startup, supervisor, instalação de pacote, CMake ou Meson. CMake/Meson não estão disponíveis nesta máquina. Sem instalação, rede, configuração global/pessoal ou alterações no repositório/Vault/Git.

Allowlist dos samples reutilizáveis: main.c/main.cpp, Makefile, mudarro.yaml e instruções. Excluir build/, .mudarro/, menu.sh, binários, saídas/caches. Não há lock de dependências: somente biblioteca padrão e compilador/Make existentes.

Contrato ainda depende da decisão pai/usuário para classificação de fontes C/C++ misturadas e projetos sem Make. Detecção futura proposta: leituras limitadas/contidas/regulares, scan sem executar ferramentas, targets Make opt-in, sem adivinhar install/start/framework/banco. Esta preparação não implementa nem atesta suporte CMake/Meson.
