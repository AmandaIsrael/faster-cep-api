# Faster CEP API

Uma API REST em Go que busca informações de CEP de forma otimizada, utilizando **multithreading** para consultar simultaneamente duas APIs externas e retornar o resultado mais rápido.

## 🚀 APIs Utilizadas

- **BrasilAPI**: `https://brasilapi.com.br/api/cep/v1/{cep}`
- **ViaCEP**: `http://viacep.com.br/ws/{cep}/json/`

## ⚡ Como Funciona

1. Recebe uma requisição HTTP com um CEP
2. Valida o formato do CEP (8 dígitos numéricos)
3. Dispara duas goroutines simultaneamente para consultar ambas as APIs
4. Retorna o primeiro resultado que chegar
5. Aplica timeout de 1 segundo (configurável)
6. Exibe logs detalhados no terminal

## 📦 Estrutura do Projeto

```
faster-cep-api/
├── cmd/
│   └── main.go                    # Ponto de entrada da aplicação
├── configs/
│   └── config.go                  # Configurações da aplicação
├── internal/
│   ├── dto/
│   │   └── cep.go                # DTO de resposta
│   ├── entity/
│   │   ├── brasilapi_cep.go      # Entidade BrasilAPI
│   │   └── via_cep.go            # Entidade ViaCEP
│   └── infra/
│       ├── gateway/
│       │   └── cep_gateway.go    # Gateway para APIs externas
│       └── handlers/
│           └── cep_handler.go    # Handler HTTP
├── pkg/
│   └── validations.go            # Validações utilitárias
├── test/
│   └── cep.http                  # Arquivo de teste HTTP
├── go.mod                        # Dependências Go
└── README.md                     # Este arquivo
```

## 🔧 Pré-requisitos

- Go 1.19 ou superior
- Conexão com a internet (para acessar as APIs)

## 📖 Como Executar

1. **Clone o repositório:**
```bash
git clone <url-do-repositorio>
cd faster-cep-api
```

2. **Instale as dependências:**
```bash
go mod tidy
```

3. **Execute a aplicação:**
```bash
go run cmd/main.go
```

4. **A API estará disponível em:** `http://localhost:8080`

## 🌐 Endpoints

### `GET /{cep}`

Busca informações de um CEP específico.

**Exemplo de requisição:**
```bash
curl http://localhost:8080/01153000
```

**Exemplo de resposta:**
```json
{
  "cep": "01153-000",
  "logradouro": "Rua Vitorino Carmilo",
  "bairro": "Campos Elíseos",
  "localidade": "São Paulo",
  "uf": "SP",
  "estado": "SP",
  "regiao": "Sudeste",
  "ibge": "3550308",
  "ddd": "11"
}
```

**Códigos de status:**
- `200`: Sucesso
- `400`: CEP inválido ou malformado
- `500`: Erro interno (falha em ambas as APIs)
- `504`: Timeout (nenhuma API respondeu em 1 segundo)

## ⚙️ Configurações

A aplicação suporta configuração via variáveis de ambiente:

| Variável | Descrição | Padrão |
|----------|-----------|---------|
| `PORT` | Porta do servidor | `8080` |
| `TIMEOUT` | Timeout das requisições | `1s` |
| `BRASILAPI_URL` | URL da BrasilAPI | `https://brasilapi.com.br/api/cep/v1/%s` |
| `VIACEP_URL` | URL da ViaCEP | `http://viacep.com.br/ws/%s/json/` |

**Exemplo de uso:**
```bash
export PORT=9090
export TIMEOUT=2s
go run cmd/main.go
```

## 🧪 Teste

Use o arquivo `test/cep.http` para testar a API:

```http
GET http://localhost:8080/65055356 HTTP/1.1
```

Ou use curl:
```bash
curl http://localhost:8080/01153000
```

## 📋 Logs

A aplicação exibe logs detalhados no terminal:

```
=== RESULTADO DA CONSULTA CEP ===
API Utilizada: BrasilAPI
CEP: 01153-000
Logradouro: Rua Vitorino Carmilo
Bairro: Campos Elíseos
Localidade: São Paulo
UF: SP
```

## 🏗️ Arquitetura

### Camadas:
- **Handler**: Processa requisições HTTP
- **Gateway**: Abstrai comunicação com APIs externas
- **Entity**: Representa estruturas das APIs externas
- **DTO**: Objeto de transferência de dados
- **Config**: Centralizador de configurações

### Padrões Utilizados:
- **Clean Architecture**: Separação clara de responsabilidades
- **Dependency Injection**: Injeção de dependências
- **Gateway Pattern**: Abstração de APIs externas
- **Race Condition**: Primeira resposta vence

## 🚀 Características Técnicas

- **Multithreading**: Goroutines para requisições simultâneas
- **Context**: Controle de timeout e cancelamento
- **Graceful Shutdown**: Encerramento elegante da aplicação
- **Validação**: CEP deve ter exatamente 8 dígitos numéricos
- **Configuração Flexível**: Via variáveis de ambiente
- **Logs Estruturados**: Informações detalhadas de debug

## 🔄 Dependências

```go
require (
    github.com/go-chi/chi/v5 v5.x.x
)
```

## 📄 Licença

Este projeto está sob licença MIT.