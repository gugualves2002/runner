# Sistema Runner

Runner é um projeto acadêmico de integração entre uma CLI Go (`assinatura`) e um serviço Java (`assinador.jar`) para simular operações de assinatura digital.

O objetivo é demonstrar:
- integração HTTP entre um cliente e um servidor Java;
- orquestração de lifecycle do serviço (`start`, `stop`, `sign`, `validate`);
- modelagem de payloads JSON e tratamento de erros estruturados;
- separação de responsabilidades entre CLI, cliente HTTP e serviço de assinatura.

## Estrutura do repositório
- `assinatura/`: CLI Go que controla o `assinador.jar`, inicia/paralisa servidor e envia requisições HTTP.
- `assinador/`: serviço Java que expõe endpoints REST para assinatura e validação.
- `internal/cli/`: cliente HTTP reutilizável e modelos de requisição do CLI.
- `docs/`: documentação de requisitos, planejamento, critérios e decisões de arquitetura.

## Como usar
- `assinatura start` inicia o servidor Java em `localhost:7070`.
- `assinatura stop` encerra a instância do servidor.
- `assinatura sign <dados_base64> <algoritmo>` cria uma assinatura via servidor HTTP.
- `assinatura validate <dados_base64> <assinatura_base64> <algoritmo>` valida uma assinatura via servidor HTTP.

## Documentação relevante
- [Especificação](especificacao.md)
- [Design](design.md)
- [Plano de implementação](docs/plano-revisitado-v2.md)
- [Critérios de aceitação](./docs/criterios.md)
- [ADR HTTP Integration](./docs/adr-http-integration.md)
- [Contrato de integração](./docs/integration-contract.md)

## Estado atual
O repositório foi organizado para suportar integração via HTTP entre o CLI Go e o serviço Java, com payloads JSON e controle de processo por porta padrão.
