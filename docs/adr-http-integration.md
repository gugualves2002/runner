# ADR 0001: HTTP Server Mode for CLI ↔ JAR Integration

## Status
Proposed

## Context
O projeto Runner integra um CLI Go (`assinatura`) com um serviço Java (`assinador.jar`) que realiza operações de assinatura digital simulada. O CLI precisa invocar a lógica do Java de forma confiável, testável e com boa experiência de uso.

Duas abordagens foram consideradas:

1. Invocação direta (cold start): o CLI executa o JAR em cada comando, passando argumentos e lendo saída padrão.
2. Modo servidor HTTP: o JAR é iniciado uma vez como um servidor persistente e o CLI se comunica via HTTP com o serviço.

## Decision
Adota-se o modo servidor HTTP como padrão de integração entre o CLI e o JAR.

## Rationale
- **Desempenho**: evita a latência de startup do JDK para cada comando. O servidor permanece em execução e atende múltiplas requisições.
- **Confiabilidade**: permite health checks e políticas de retry antes de executar chamadas de assinatura ou validação.
- **Separação de responsabilidades**: o CLI gerencia lifecycle e invocação HTTP; o JAR foca no serviço de assinatura e validação.
- **Testabilidade**: a comunicação HTTP pode ser testada com `httptest` no Go e com mocks/servidores de teste no Java.
- **Escalabilidade**: o servidor pode crescer em features sem exigir mudanças na interface de chamada do CLI, mantendo um contrato HTTP estável.

## Benefits
- `assinatura start` abre o servidor e `assinatura stop` controla a instância, tornando o uso mais próximo de um serviço de produção.
- O modo servidor facilita a integração com múltiplos clientes e permite reuso do mesmo processo Java.
- Erros são tratados como respostas HTTP estruturadas, simplificando diagnóstico e logs.

## Consequences
- Requer que o JAR seja executado e mantido em memória durante a sessão.
- A porta padrão deve ser gerenciada e liberada ao parar o servidor.
- A implementação precisa garantir que `assinador.jar` escute corretamente em `localhost:7070` e exponha a API JSON.

## Alternatives considered
- **Cold start direto**: bom para execução única, mas menos eficiente e mais frágil para comportamento de servidor.
- **IPC nativo**: complexo e não necessário para o escopo acadêmico; a comunicação HTTP é suficiente e familiar.

## Notes
O contrato técnico da integração está documentado em `docs/integration-contract.md`.
