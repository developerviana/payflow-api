# 📁 Estrutura do Projeto

```
payflow-api/
├── cmd/
│   └── server/              # Ponto de entrada da aplicação
│       └── main.go          # Arquivo principal
├── internal/
│   ├── entity/              # Entidades do domínio
│   │   ├── user.go          # Entidade User
│   │   └── transaction.go   # Entidade Transaction
│   ├── usecase/             # Casos de uso (regras de negócio)
│   │   ├── interfaces/      # Interfaces dos casos de uso
│   │   ├── user/            # Casos de uso de usuário
│   │   └── transaction/     # Casos de uso de transação
│   ├── repository/          # Interfaces e implementações de repositório
│   │   ├── interfaces/      # Contratos dos repositórios
│   │   └── postgres/        # Implementação PostgreSQL
│   ├── handler/             # Handlers HTTP (controllers)
│   │   ├── user.go          # Handlers de usuário
│   │   ├── transaction.go   # Handlers de transação
│   │   └── middleware/      # Middlewares
│   ├── service/             # Serviços externos
│   │   ├── authorizer.go    # Serviço de autorização
│   │   └── notification.go  # Serviço de notificação
│   └── config/              # Configurações da aplicação
│       └── config.go        # Configurações gerais
├── pkg/
│   ├── database/            # Configuração do banco de dados
│   │   └── postgres.go      # Conexão PostgreSQL
│   ├── validator/           # Validadores customizados
│   │   ├── cpf.go           # Validador de CPF
│   │   └── cnpj.go          # Validador de CNPJ
│   └── logger/              # Configuração de logs
│       └── logger.go        # Logger estruturado
├── test/
│   ├── unit/                # Testes unitários
│   ├── integration/         # Testes de integração
│   └── mocks/               # Mocks para testes
├── migrations/              # Migrações do banco de dados
│   ├── 001_create_users.sql
│   └── 002_create_transactions.sql
├── docs/                    # Documentação
│   └── api.md               # Documentação da API
├── docker-compose.yml       # Configuração do Docker
├── Dockerfile              # Imagem Docker da aplicação
├── .env.example            # Exemplo de variáveis de ambiente
└── Makefile                # Comandos de automação
```

## 🏗️ Arquitetura Clean Architecture

### Camadas da Aplicação

1. **Entities (internal/entity/)**
   - Contém as regras de negócio mais fundamentais
   - Estruturas de dados centrais do domínio

2. **Use Cases (internal/usecase/)**
   - Contém as regras de negócio específicas da aplicação
   - Orquestra o fluxo de dados entre entities e repositories

3. **Interface Adapters (internal/handler/, internal/repository/)**
   - **Handlers**: Adaptam dados da web para use cases
   - **Repositories**: Adaptam dados do banco para use cases

4. **Frameworks & Drivers (pkg/, cmd/)**
   - Detalhes de implementação (banco, web framework, etc.)

### Fluxo de Dependências

```
Handler → UseCase → Repository → Database
   ↓         ↓         ↓
Entity ← Entity ← Entity
```

- **Dependency Rule**: Dependências apontam sempre para dentro
- **Inversion of Control**: Use cases definem interfaces, implementações ficam nas camadas externas
