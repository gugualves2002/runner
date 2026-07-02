# Documentação Técnica da Integração HTTP

## Visão Geral
O CLI `assinatura` se integra ao serviço Java `assinador.jar` via HTTP JSON. O JAR expõe endpoints REST em `http://localhost:7070/api/v1`, e o CLI envia requisições POST para os seguintes recursos:

- `/sign`
- `/validate`

O servidor deve estar em execução para que o CLI funcione nesses comandos. O modo servidor é iniciado por `assinatura start`.

## Endpoint: POST /api/v1/sign

### Payload
```json
{
  "data": "<dados_base64>",
  "algorithm": "<algoritmo>",
  "pkcs11ConfigPath": "<caminho-opcional>",
  "pin": "<pin-opcional>",
  "alias": "<alias-opcional>"
}
```

### Resposta de sucesso
```json
{
  "signature": "<assinatura_base64>"
}
```

### Resposta de erro
- `400 Bad Request` quando o payload é inválido ou parâmetros ausentes.
- Corpo de erro JSON:
```json
{
  "message": "<descrição do erro>"
}
```

## Endpoint: POST /api/v1/validate

### Payload
```json
{
  "data": "<dados_base64>",
  "signature": "<assinatura_base64>",
  "algorithm": "<algoritmo>",
  "pkcs11ConfigPath": "<caminho-opcional>",
  "pin": "<pin-opcional>",
  "alias": "<alias-opcional>"
}
```

### Resposta de sucesso
```json
{
  "valid": true
}
```
ou
```json
{
  "valid": false
}
```

### Resposta de erro
- `400 Bad Request` quando o payload é inválido ou parâmetros ausentes.
- Corpo de erro JSON igual a:
```json
{
  "message": "<descrição do erro>"
}
```

## Contrato CLI ↔ JAR

### Requisições do CLI
- O CLI constrói JSON a partir dos argumentos de linha de comando.
- Os campos obrigatórios são:
  - `data`
  - `algorithm`
- Para validação, o CLI também envia `signature`.
- Campos PKCS#11 (`pkcs11ConfigPath`, `pin`, `alias`) são opcionais, mas se `pkcs11ConfigPath` for fornecido, `alias` também deve ser fornecido.

### Tratamento de erros no CLI
- Erros de flag são tratados antes da chamada HTTP.
- Se o servidor retornar código não `200`, o CLI tenta decodificar um JSON de erro e retorna mensagem estruturada.
- Mensagem de fallback: `servidor retornou erro (<status>): <body>`.

### Exemplo de fluxo
1. Usuário roda `assinatura start`.
2. O CLI inicia o `assinador.jar` em background e aguarda health check.
3. Usuário roda `assinatura sign <dados_base64> <algoritmo>`.
4. O CLI faz POST `/sign` para `http://localhost:7070/api/v1/sign`.
5. O servidor responde com a assinatura.
6. Usuário roda `assinatura validate <dados_base64> <assinatura_base64> <algoritmo>`.
7. O CLI faz POST `/validate` para `http://localhost:7070/api/v1/validate`.

## Códigos de erro e mensagens esperadas
- `400 Bad Request`: payload inválido, JSON malformado ou parâmetros de validação ausentes.
- `500 Internal Server Error`: erro interno do servidor Java, possivelmente causado por falha no PKCS#11 ou no processo.
- `502 Bad Gateway` / `503 Service Unavailable`: não aplicável diretamente, mas o CLI trata falhas de conexão como servidor indisponível.

## Observações de segurança
- O servidor roda apenas em `localhost` por design de desenvolvimento.
- O contrato usa JSON e não inclui autenticação no escopo atual.
- Dados em Base64 transitam no payload; recomenda-se transporte seguro se evoluir para ambiente real.
