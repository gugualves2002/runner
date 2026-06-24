# Sprint 2: Assinatura Digital Simulada (Modo Local)

Este documento descreve as implementações da Sprint 2.

## Estrutura Implementada

### Aplicação Java (assinador)

**Localização**: `assinador/`

Projeto Maven que implementa a lógica de assinatura digital simulada.

#### Classes Principais

- `SignatureService` - Interface que define operações de assinatura
- `FakeSignatureService` - Implementação simulada de assinatura (sem criptografia real)
- `SignatureException` - Exceção customizada para erros
- `SignatureResponse` - Classe de resposta JSON
- `Main` - Ponto de entrada que funciona como servidor CLI

#### Comandos do assinador.jar

```bash
# Criar assinatura
java -jar assinador.jar sign --payload "conteúdo" --key-alias "minha-chave"

# Validar assinatura
java -jar assinador.jar validate --payload "conteúdo" --signature "..." --key-alias "minha-chave"

# Ver ajuda
java -jar assinador.jar --help
```

#### Compilar o projeto

```bash
cd assinador
mvn clean package
```

Gera: `assinador/target/assinador.jar`

#### Executar testes

```bash
cd assinador
mvn test
```

### CLI em Go (assinatura)

**Localização**: `cmd/assinatura/`

CLI multiplataforma que invoca assinador.jar localmente.

#### Comandos do CLI

```bash
# Ver versão
assinatura version

# Criar assinatura (invoca assinador.jar localmente)
assinatura sign --payload "Olá Mundo" --key-alias minha-chave

# Validar assinatura
assinatura validate --payload "Olá Mundo" --signature "..." --key-alias minha-chave

# Ver ajuda
assinatura --help
```

#### Compilar CLI

```bash
# Build para sua plataforma
go build -o assinatura ./cmd/assinatura

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o assinatura-linux-amd64 ./cmd/assinatura
GOOS=windows GOARCH=amd64 go build -o assinatura-windows-amd64.exe ./cmd/assinatura
GOOS=darwin GOARCH=amd64 go build -o assinatura-darwin-amd64 ./cmd/assinatura
```

#### Executar testes

```bash
go test ./cmd/assinatura/...
```

### Pacotes Internos

#### internal/jdk

**Arquivo**: `internal/jdk/manager.go`

Gerencia detecção e provisionamento automático de JDK.

- Detecta Java no PATH
- Procura em `~/.hubsaude/jdk/`
- Baixa JDK do Adoptium (Eclipse Temurin) se necessário

#### internal/invoker

**Arquivo**: `internal/invoker/invoker.go`

Gerencia invocação do assinador.jar em modo local.

- Localiza assinador.jar
- Executa operações de sign/validate
- Parse e formatação de respostas JSON

## Fluxo de Execução

### Operação de Assinatura

```
┌──────────────────┐
│   Usuário        │
└──────────────────┘
        │
        │ assinatura sign --payload "..." --key-alias "..."
        ▼
┌──────────────────────────────────┐
│   CLI Go (cmd/assinatura)       │
│  - Parse de argumentos           │
│  - Validação de flags            │
└──────────────────────────────────┘
        │
        │ java -jar assinador.jar sign ...
        ▼
┌──────────────────────────────────┐
│   assinador.jar                  │
│  - Validação de parâmetros       │
│  - FakeSignatureService.sign()   │
│  - Resposta JSON                 │
└──────────────────────────────────┘
        │
        │ JSON com assinatura
        ▼
┌──────────────────┐
│   Usuário        │
│  (resultado)     │
└──────────────────┘
```

## Validação de Parâmetros

### sign

- `--payload`: Obrigatório, não pode ser vazio
- `--key-alias`: Obrigatório, não pode ser vazio

Erro retornado como JSON:
```json
{
  "status": "error",
  "message": "Payload não pode ser vazio"
}
```

### validate

- `--payload`: Obrigatório
- `--signature`: Obrigatório, deve ser Base64 válido
- `--key-alias`: Obrigatório

## Exemplo de Uso Completo

### Pré-requisitos

1. Java 21+ instalado
2. Go 1.25+ instalado
3. Maven (para compilar assinador.jar)

### Passos

```bash
# 1. Compilar assinador.jar
cd assinador
mvn clean package
cd ..

# 2. Compilar CLI
go build -o assinatura ./cmd/assinatura

# 3. Criar assinatura
./assinatura sign --payload "Documento importante" --key-alias "chave-1"

# Resposta esperada (JSON):
# {
#   "signature": "RkFLRV9TSUdOQVRVUkVfLTIwMjE4OTU0X...",
#   "status": "success",
#   "message": "Assinatura criada com sucesso"
# }

# 4. Validar assinatura (copie a signature da resposta anterior)
./assinatura validate \
  --payload "Documento importante" \
  --signature "RkFLRV9TSUdOQVRVUkVfLTIwMjE4OTU0X..." \
  --key-alias "chave-1"

# Resposta esperada (JSON):
# {
#   "status": "valid",
#   "message": "Assinatura válida"
# }
```

## Status da Sprint 2

### Tarefas Completadas

- ✅ US-02.1: Simulação de criação de assinatura digital
  - Interface SignatureService definida
  - Implementação FakeSignatureService
  - Testes unitários cobrem todos os cenários

- ✅ US-02.2: Validação de parâmetros de criação
  - Validação no assinador.jar
  - Mensagens de erro claras

- ✅ US-02.3: Simulação e validação de assinatura
  - Método validate implementado
  - Validação de parâmetros

- ✅ US-01.2: Parsing de comandos e parâmetros no CLI
  - Cobra framework integrado
  - Flags obrigatórias e opcionais

- ✅ US-01.3: Invocação do assinador.jar no modo local
  - CLI invoca assinador.jar via java -jar
  - Propagação de argumentos

- ✅ US-01.4: Exibição legível de resultados
  - Respostas em JSON formatado
  - Saída estruturada

- ✅ US-04.1: Detecção e provisionamento automático do JDK
  - Manager para detecção de Java
  - Busca em PATH e ~/.hubsaude/
  - Download automático (estrutura preparada)

## Próximas Sprints

- **Sprint 3**: Modo servidor HTTP para assinador.jar
- **Sprint 4**: CLI para gerenciar Simulador do HubSaúde
