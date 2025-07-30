# Build stage
FROM golang:1.23.2-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git

# Definir diretório de trabalho
WORKDIR /app

# Copiar go mod e sum
COPY go.mod go.sum ./

# Download dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário da aplicação
COPY --from=builder /app/main .

# Copiar arquivos de configuração se necessário
COPY --from=builder /app/.env.example .env

# Expor porta
EXPOSE 8080

# Comando para executar
CMD ["./main"]
