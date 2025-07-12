# 🏗️ Clean Architecture - Sistema de Orders

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Architecture](https://img.shields.io/badge/Architecture-Clean-blue?style=flat-square)
![APIs](https://img.shields.io/badge/APIs-REST%20%7C%20gRPC%20%7C%20GraphQL-orange?style=flat-square)

Este projeto implementa um sistema de gerenciamento de orders usando **Clean Architecture** em Go, com múltiplas interfaces de comunicação (REST, gRPC e GraphQL).

## 📦 **Repositório**

**GitHub**: https://github.com/danielencestari/pos_03

Clone o projeto:
```bash
git clone https://github.com/danielencestari/pos_03.git
cd pos_03/CleanArch
```

## 📋 **Funcionalidades**

- ✅ **Criar Orders** via REST, gRPC e GraphQL
- ✅ **Listar Orders** com paginação via REST, gRPC e GraphQL  
- ✅ **Ordenação** por diferentes campos (id, price, tax, final_price)
- ✅ **Paginação** configurável (page, limit)
- ✅ **Eventos assíncronos** com RabbitMQ
- ✅ **Clean Architecture** com separação clara de responsabilidades

## 🗂️ **Arquitetura do Projeto**

```
├── cmd/ordersystem/          # Ponto de entrada da aplicação
├── internal/
│   ├── entity/              # 🟡 Entities - Regras de negócio
│   ├── usecase/             # 🟢 Use Cases - Casos de uso
│   ├── infra/               # 🔵 Interface Adapters
│   │   ├── database/        # Repositórios (MySQL)
│   │   ├── web/             # Handlers REST
│   │   ├── grpc/            # Services gRPC
│   │   └── graph/           # Resolvers GraphQL
│   └── event/               # Sistema de eventos
├── configs/                 # Configurações
├── api/                     # Testes HTTP
└── sql/migrations/          # Scripts de migração
```

## 🚀 **Como Executar**

### **1️⃣ Pré-requisitos**
- Docker e Docker Compose
- Go 1.24+ (para desenvolvimento)

### **2️⃣ Subir a aplicação**
```bash
# Clona o repositório
git clone https://github.com/danielencestari/pos_03.git
cd pos_03/CleanArch

# Sobe todos os serviços
docker compose up -d
```

### **3️⃣ Executar migração do banco**
```bash
# Conecta no MySQL
docker exec -it mysql mysql -u root -proot orders

# Executa a migração
source sql/migrations/001_create_orders_table.sql;
```

### **4️⃣ Executar a aplicação**
```bash
# Via Docker (recomendado)
docker compose up app

# Ou localmente (para desenvolvimento)
go run cmd/ordersystem/main.go
```

## 🌐 **Portas dos Serviços**

| Serviço | Porta | URL | Descrição |
|---------|-------|-----|-----------|
| **REST API** | `8000` | http://localhost:8000 | API RESTful |
| **gRPC** | `50051` | localhost:50051 | Serviço gRPC |
| **GraphQL** | `8080` | http://localhost:8080 | Playground GraphQL |
| **MySQL** | `3306` | localhost:3306 | Banco de dados |
| **RabbitMQ** | `15672` | http://localhost:15672 | Management UI |

## 🧪 **Como Testar**

### **📡 REST API**

**Criar Order:**
```bash
POST http://localhost:8000/order
Content-Type: application/json

{
    "id": "order-001",
    "price": 100.5,
    "tax": 10.0
}
```

**Listar Orders:**
```bash
# Listar com paginação
GET http://localhost:8000/order?page=1&limit=10&sort=id

# Listar com parâmetros padrão
GET http://localhost:8000/order
```

### **⚡ gRPC**

**Usando evans:**
```bash
# Instalar evans
brew install evans

# Conectar ao servidor
evans -r repl

# Selecionar package e service
package pb
service OrderService

# Criar order
call CreateOrder
# Insira: id="order-001", price=100.5, tax=10.0

# Listar orders
call ListOrders
# Insira: page=1, limit=10, sort="id"
```

### **🎭 GraphQL**

Acesse http://localhost:8080 e use:

**Criar Order:**
```graphql
mutation {
  createOrder(input: {
    id: "order-001"
    Price: 100.5
    Tax: 10.0
  }) {
    id
    Price
    Tax
    FinalPrice
  }
}
```

**Listar Orders:**
```graphql
query {
  listOrders(page: 1, limit: 10, sort: "id") {
    orders {
      id
      Price
      Tax
      FinalPrice
    }
    page
    limit
    total
    totalPages
  }
}
```

## 🗄️ **Estrutura do Banco de Dados**

### **Tabela: orders**
| Campo | Tipo | Descrição |
|-------|------|-----------|
| `id` | VARCHAR(36) | UUID único da order |
| `price` | DECIMAL(10,2) | Preço base |
| `tax` | DECIMAL(10,2) | Taxa aplicada |
| `final_price` | DECIMAL(10,2) | Preço final (price + tax) |
| `created_at` | TIMESTAMP | Data de criação |
| `updated_at` | TIMESTAMP | Data de atualização |

## 🔧 **Configurações (.env)**

```env
# Banco de Dados
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders

# Servidores
WEB_SERVER_PORT=8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

## 📊 **Paginação**

Todos os endpoints de listagem suportam paginação:

- **page**: Número da página (padrão: 1)
- **limit**: Registros por página (padrão: 10)
- **sort**: Campo para ordenação (padrão: "id")

**Exemplo de resposta:**
```json
{
  "orders": [...],
  "page": 1,
  "limit": 10,
  "total": 25,
  "total_pages": 3
}
```

## 🏗️ **Clean Architecture Explicada**

### **🟡 Entities (Regras de Negócio)**
- `Order`: Entidade principal com validações
- Métodos: `NewOrder()`, `IsValid()`, `CalculateFinalPrice()`

### **🟢 Use Cases (Casos de Uso)**
- `CreateOrderUseCase`: Criar uma nova order
- `ListOrdersUseCase`: Listar orders com paginação

### **🔵 Interface Adapters**
- **REST**: Handlers HTTP com conversão JSON ↔ DTO
- **gRPC**: Services com conversão Proto ↔ DTO  
- **GraphQL**: Resolvers com conversão GraphQL ↔ DTO

### **🔴 Frameworks & Drivers**
- **Database**: Repository MySQL
- **Events**: RabbitMQ para comunicação assíncrona

## 🛠️ **Tecnologias Utilizadas**

- **Go 1.24+** - Linguagem principal
- **Gin** - Framework web (via chi/v5)
- **gRPC** - Comunicação RPC
- **GraphQL** - API flexível (gqlgen)
- **MySQL** - Banco de dados
- **RabbitMQ** - Mensageria
- **Docker** - Containerização
- **Wire** - Injeção de dependência

## 🧪 **Testes**

```bash
# Executar testes
go test ./...

# Testes com coverage
go test -cover ./...

# Testes específicos
go test ./internal/entity/
go test ./internal/usecase/
```

## 📝 **Scripts Úteis**

```bash
# Gerar código gRPC
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

# Gerar código GraphQL
go run github.com/99designs/gqlgen generate

# Gerar injeção de dependência
wire

# Executar linting
golangci-lint run
```

---

## 🤝 **Contribuindo**

Contribuições são sempre bem-vindas! Para contribuir:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 **Licença**

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👤 **Autor**

**Daniel Encestari**
- GitHub: [@danielencestari](https://github.com/danielencestari)
- Repositório: [pos_03](https://github.com/danielencestari/pos_03)

## 🌟 **Agradecimentos**

- [Full Cycle](https://fullcycle.com.br/) pelo conhecimento em Clean Architecture
- Comunidade Go pelo suporte e recursos
- Todos os contribuidores que ajudaram a tornar este projeto melhor

---

**Desenvolvido seguindo os princípios da Clean Architecture** 🏗️ 

**Tecnologias**: Go 1.24 • Clean Architecture • REST • gRPC • GraphQL • MySQL • RabbitMQ • Docker 