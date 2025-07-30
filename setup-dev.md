# Guia de Configuração do Ambiente de Desenvolvimento - PayFlow API

## ✅ Situação Atual

Seu projeto está configurado com:
- ✅ Clean Architecture implementada
- ✅ Entidades de domínio (User e Transaction) com validação completa
- ✅ Testes unitários funcionais
- ✅ DTOs e conversores prontos
- ✅ Servidor HTTP básico configurado

## 🔧 Configuração Atual de Portas

**Problema Identificado:** Conflitos de porta no ambiente local

**Solução Aplicada:**
- PostgreSQL: Movido para porta 5433 (evita conflito com instalação local na 5432)
- API Server: Configurado para porta 3000 (evita conflitos com 8080/8081)

## 📋 Próximos Passos para Ambiente Solo

### 1. Opção Simples - Desenvolvimento sem Docker

```bash
# 1. Configure um banco PostgreSQL local
# Crie o banco: payflow
# Use as credenciais: postgres/postgres

# 2. Execute as migrações manualmente
# Copie o conteúdo dos arquivos em migrations/ e execute no banco

# 3. Rode a aplicação
go run cmd/server/main.go
```

### 2. Opção Docker - Quando resolver problemas de espaço

```bash
# 1. Limpe o Docker para liberar espaço
docker system prune -a

# 2. Suba apenas o banco
docker-compose up postgres -d

# 3. Rode a aplicação localmente
go run cmd/server/main.go
```

### 3. Configuração Atual de Arquivos

**`.env` configurado para:**
```env
SERVER_PORT=3000
DB_HOST=localhost
DB_PORT=5432  # ou 5433 se usar Docker
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=payflow
```

## 🚀 Desenvolvimento Solo - Workflow Recomendado

1. **Para desenvolvimento rápido:** Use PostgreSQL local + API local
2. **Para environment isolation:** Resolva problemas de Docker primeiro
3. **Para testes:** Execute `go test ./...`

## 📊 Endpoints Disponíveis

Uma vez funcionando, sua API terá:
- `GET /api/v1/health` - Health check
- `GET /api/v1/info` - Informações da API
- `GET /api/v1/users/` - Placeholder para usuários
- `GET /api/v1/transactions/` - Placeholder para transações

## 🔍 Debug de Problemas

**Espaço em disco:** 
- Limpe arquivos temporários do Go: `go clean -cache`
- Verifique espaço disponível

**Conflitos de porta:**
- Use `netstat -an | grep LISTENING` para verificar portas ocupadas
- Ajuste SERVER_PORT no .env conforme necessário

**Problemas de Docker:**
- Verifique se Docker Desktop está rodando
- Execute `docker system df` para ver uso de espaço

## ✨ Próximas Implementações

Após resolver o ambiente, implemente na ordem:
1. Repository layer (interfaces + PostgreSQL)
2. UseCase layer (business logic)
3. Handler layer (HTTP endpoints completos)
4. Middleware de autenticação
5. Testes de integração
