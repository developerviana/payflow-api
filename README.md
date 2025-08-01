# 💸 payflow-api

`payflow-api` é uma API RESTful desenvolvida em Go como parte de um desafio técnico. O projeto simula uma plataforma de pagamentos simplificada, permitindo transferências financeiras entre usuários com regras específicas de negócio.

Este desafio teve como objetivo aplicar conceitos sólidos de engenharia de software, como arquitetura limpa, responsabilidade única, testes automatizados, integração com serviços externos e uso de containers.

---

## 🧠 Objetivo do Projeto

- Simular um sistema real de pagamentos.
- Trabalhar com diferentes tipos de usuários (comum e lojista).
- Implementar regras de negócio como restrição de envio, validação de saldo e transações atômicas.
- Integrar com serviços externos simulados para autorização e notificação.
- Garantir código limpo, organizado e de fácil manutenção.

---

## ⚙️ Funcionalidades Implementadas

### ✅ **Sistema de Usuários**
- Cadastro completo de usuários com validação de CPF/CNPJ e e-mail únicos
- Tipos de usuário: comum e lojista
- Criptografia de senhas com bcrypt
- Operações CRUD completas (Create, Read, Update, Delete)
- Consulta de saldo
- Listagem com paginação e filtros

### ✅ **Validações e Regras de Negócio**
- Validação rigorosa de CPF/CNPJ usando algoritmos oficiais
- Validação de formato de email
- Lojistas não podem ser pagadores (restrição de operação)
- Verificação de saldo antes da transferência
- Documentos e emails únicos no sistema

### ✅ **Arquitetura e Infraestrutura**
- Clean Architecture com separação clara entre camadas
- Handlers HTTP, Use Cases, Repositories e Entities
- Conexão com PostgreSQL
- Migrations de banco de dados
- Testes unitários abrangentes
- Docker e Docker Compose configurados
- API RESTful com tratamento de erros padronizado

### ⏳ **Em Desenvolvimento**
- Sistema completo de transações
- Integração com serviços externos de autorização
- Sistema de notificações
- Transações atômicas com rollback

---

## 📦 Tecnologias Utilizadas

- **Go (Golang) 1.23** – Linguagem principal da aplicação
- **Gin Framework** – Framework web para rotas HTTP
- **PostgreSQL 16** – Banco de dados relacional
- **Docker & Docker Compose** – Containerização do ambiente
- **Clean Architecture** – Padrão arquitetural
- **bcrypt** – Criptografia de senhas
- **Testify** – Framework de testes unitários
- **Go Modules** – Gerenciamento de dependências
- **github.com/lib/pq** – Driver PostgreSQL

---

## 🚀 Como Executar o Projeto

### **Pré-requisitos**
- Go 1.23+
- Docker e Docker Compose
- PostgreSQL (ou usar o container)

### **1. Clonar o Repositório**
```bash
git clone https://github.com/developerviana/payflow-api.git
cd payflow-api
```

### **2. Configurar Variáveis de Ambiente**
Crie um arquivo `.env` na raiz do projeto:
```env
# Configurações do Servidor
SERVER_PORT=8080
ENVIRONMENT=development

# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=payflow
DB_SSL_MODE=disable

# Configurações de Serviços Externos
AUTHORIZER_URL=https://util.devi.tools/api/v2/authorize
NOTIFICATION_URL=https://util.devi.tools/api/v1/notify
REQUEST_TIMEOUT=10
```

### **3. Subir o Banco de Dados**
```bash
# Apenas PostgreSQL
docker-compose up postgres -d

# Ou toda a aplicação
docker-compose --profile full up -d
```

### **4. Executar a Aplicação**
```bash
# Instalar dependências
go mod tidy

```

### **5. Verificar se está funcionando**
```bash
curl http://localhost:8080/api/v1/health
```

---

## 🌐 Endpoints da API

### **Base URL:** `http://localhost:8080`

### **🔧 Sistema**
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/api/v1/health` | Health check da API |
| `GET` | `/api/v1/info` | Informações da API |

### **👥 Usuários**
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/api/v1/users` | Criar novo usuário |
| `GET` | `/api/v1/users` | Listar usuários (com paginação) |
| `GET` | `/api/v1/users/:id` | Buscar usuário por ID |
| `PUT` | `/api/v1/users/:id` | Atualizar usuário |
| `DELETE` | `/api/v1/users/:id` | Deletar usuário |
| `GET` | `/api/v1/users/:id/balance` | Consultar saldo |

### **💸 Transações**
| Método | Endpoint | Descrição | Status |
|--------|----------|-----------|--------|
| `POST` | `/api/v1/transactions` | Criar transação | ⏳ Em desenvolvimento |
| `GET` | `/api/v1/transactions` | Listar transações | ⏳ Em desenvolvimento |

---

## 📝 Exemplos de Uso

### **Criar Usuário Comum**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "João Silva",
    "document": "11144477735",
    "email": "joao@teste.com",
    "password": "123456",
    "user_type": "common"
  }'
```

### **Criar Usuário Lojista**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Loja do João LTDA",
    "document": "11222333000181",
    "email": "loja@teste.com",
    "password": "123456",
    "user_type": "merchant"
  }'
```

### **Listar Usuários com Filtros**
```bash
# Listar todos
curl http://localhost:8080/api/v1/users

# Com paginação
curl "http://localhost:8080/api/v1/users?page=1&limit=10"

# Filtrar por tipo
curl "http://localhost:8080/api/v1/users?user_type=common"

# Busca por email
curl "http://localhost:8080/api/v1/users?email=joao"
```

### **Buscar Usuário por ID**
```bash
curl http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440001
```

### **Atualizar Usuário**
```bash
curl -X PUT http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440001 \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "João Silva Santos",
    "email": "joao.santos@teste.com"
  }'
```

### **Consultar Saldo**
```bash
curl http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440001/balance
```

---

## 🧪 Testes

```bash

go test ./test/unit/... -v



---

## 🔐 Regras de Negócio

- **Tipos de usuário**: comum e lojista
- **Lojistas não podem ser pagadores** (apenas recebedores)
- **CPF/CNPJ e e-mail devem ser únicos** no sistema
- **Validação rigorosa** de CPF/CNPJ usando algoritmos oficiais
- **Senhas criptografadas** com bcrypt
- **Verificação de saldo** antes de qualquer transferência
- **Transações atômicas** com rollback em caso de falhas
- **Limite máximo** de transação (R$ 10.000,00)

---

## 📁 Estrutura do Projeto

```
payflow-api/
├── cmd/server/              # Ponto de entrada da aplicação
├── internal/
│   ├── entity/              # Entidades do domínio
│   ├── usecase/             # Casos de uso (regras de negócio)
│   ├── repository/          # Interfaces e implementações de repositório
│   ├── handler/             # Handlers HTTP (controllers)
│   └── config/              # Configurações da aplicação
├── pkg/database/            # Configuração do banco de dados
├── test/unit/               # Testes unitários
├── migrations/              # Migrações do banco de dados
├── docs/                    # Documentação
├── docker-compose.yml       # Configuração do Docker
└── Dockerfile              # Imagem Docker da aplicação
```

---

## 🐳 Docker

### **Desenvolvimento (apenas banco)**
```bash
docker-compose up postgres -d
```

### **Produção (aplicação completa)**
```bash
docker-compose --profile full up -d
```

### **Parar containers**
```bash
docker-compose down
```

---

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

## 📧 Contato

**Desenvolvedor**: [@developerviana](https://github.com/developerviana)

**Link do Projeto**: [https://github.com/developerviana/payflow-api](https://github.com/developerviana/payflow-api)

