# 💸 payflow-api

`payflow-api` é uma API RESTful desenvolvida em Go como parte de um desafio técnico. O projeto simula uma plataforma de pagamentos simplificada, permitindo transferências financeiras entre usuários com regras específicas de negócio.

Este desafio teve como objetivo me tirar da zona de conforto e aplicar conceitos sólidos de engenharia de software, como arquitetura limpa, responsabilidade única, testes automatizados, integração com serviços externos e uso de containers.

---

## 🧠 Objetivo do Projeto

- Simular um sistema real de pagamentos.
- Trabalhar com diferentes tipos de usuários (comum e lojista).
- Implementar regras de negócio como restrição de envio, validação de saldo e transações atômicas.
- Integrar com serviços externos simulados para autorização e notificação.
- Garantir código limpo, organizado e de fácil manutenção.

---

## ⚙️ Funcionalidades Implementadas

- ✅ Cadastro de usuários com validação de CPF/CNPJ e e-mail únicos
- ✅ Tipos de usuário: comum e lojista
- ✅ Transferência de valores entre usuários comuns e lojistas
- ✅ Lojistas não podem transferir dinheiro (restrição de operação)
- ✅ Verificação de saldo antes da transferência
- ✅ Requisição ao serviço externo de autorização (mock)
- ✅ Envio de notificação ao recebedor (mock de serviço externo)
- ✅ Transações seguras com rollback em caso de falhas
- ✅ API RESTful com tratamento de erros e validações

---

## 📦 Tecnologias Utilizadas

- **Go (Golang)** – Linguagem principal da aplicação
- **Docker** – Containerização do ambiente
- **Clean Architecture** – Separação clara entre camadas
- **Gorilla Mux ou Gin** – Framework para rotas HTTP
- **PostgreSQL / SQLite (opcional)** – Persistência de dados
- **Testify** – Testes unitários e mocks
- **Go Modules** – Gerenciamento de dependências
- **Logs estruturados** – Com `logrus` ou `zap`

---

## 🔐 Regras de Negócio

- Lojistas não podem ser pagadores
- CPF/CNPJ e e-mail devem ser únicos
- O pagador precisa ter saldo suficiente
- Toda transação deve ser atômica (rollback em falhas)
- Notificação e autorização via serviços externos

---
