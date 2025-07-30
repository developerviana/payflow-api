.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down

# Variáveis
APP_NAME=payflow-api
DOCKER_COMPOSE=docker-compose

help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compila a aplicação
	go build -o bin/$(APP_NAME) cmd/server/main.go

run: ## Executa a aplicação
	go run cmd/server/main.go

test: ## Executa todos os testes
	go test -v ./...

test-coverage: ## Executa testes com coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Remove arquivos gerados
	rm -rf bin/
	rm -f coverage.out coverage.html

install: ## Instala dependências
	go mod download
	go mod tidy

docker-up: ## Sobe os containers (banco de dados)
	$(DOCKER_COMPOSE) up -d postgres

docker-down: ## Para os containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Mostra logs dos containers
	$(DOCKER_COMPOSE) logs -f

migrate-up: ## Executa migrações do banco
	@echo "Executando migrações..."
	# TODO: Implementar comando de migração

migrate-down: ## Reverte migrações do banco
	@echo "Revertendo migrações..."
	# TODO: Implementar comando de rollback

dev: ## Executa a aplicação em modo desenvolvimento
	air

format: ## Formata o código
	go fmt ./...

lint: ## Executa linter
	golangci-lint run

deps: ## Verifica dependências
	go mod verify
	go mod download

setup: ## Configuração inicial do projeto
	cp .env.example .env
	make install
	make docker-up
