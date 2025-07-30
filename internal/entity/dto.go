package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

// CreateUserRequest representa o request para criar um usuário
type CreateUserRequest struct {
	FullName string   `json:"full_name" validate:"required,min=3,max=100"`
	Document string   `json:"document" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	UserType UserType `json:"user_type" validate:"required,oneof=common merchant"`
}

// CreateUserResponse representa a resposta da criação de usuário
type CreateUserResponse struct {
	ID       string   `json:"id"`
	FullName string   `json:"full_name"`
	Document string   `json:"document"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
	Balance  string   `json:"balance"`
}

// GetUserResponse representa a resposta de consulta de usuário
type GetUserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Document  string    `json:"document"`
	Email     string    `json:"email"`
	UserType  UserType  `json:"user_type"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest representa o request para atualizar um usuário
type UpdateUserRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest representa o request para alterar senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// CreateTransactionRequest representa o request para criar uma transação
type CreateTransactionRequest struct {
	PayeeID string          `json:"payee_id" validate:"required,uuid"`
	Amount  decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

// CreateTransactionResponse representa a resposta da criação de transação
type CreateTransactionResponse struct {
	ID              string            `json:"id"`
	PayerID         string            `json:"payer_id"`
	PayeeID         string            `json:"payee_id"`
	Amount          string            `json:"amount"`
	Status          TransactionStatus `json:"status"`
	StatusDesc      string            `json:"status_description"`
	CreatedAt       time.Time         `json:"created_at"`
}

// GetTransactionResponse representa a resposta de consulta de transação
type GetTransactionResponse struct {
	ID                string            `json:"id"`
	PayerID           string            `json:"payer_id"`
	PayeeID           string            `json:"payee_id"`
	Amount            string            `json:"amount"`
	Status            TransactionStatus `json:"status"`
	StatusDesc        string            `json:"status_description"`
	AuthorizationID   *string           `json:"authorization_id,omitempty"`
	NotificationSent  bool              `json:"notification_sent"`
	FailureReason     *string           `json:"failure_reason,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	CompletedAt       *time.Time        `json:"completed_at,omitempty"`
	
	// Dados dos usuários relacionados
	Payer *UserSummary `json:"payer,omitempty"`
	Payee *UserSummary `json:"payee,omitempty"`
}

// UserSummary representa um resumo do usuário
type UserSummary struct {
	ID       string   `json:"id"`
	FullName string   `json:"full_name"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
}

// ListTransactionsResponse representa a resposta de listagem de transações
type ListTransactionsResponse struct {
	Transactions []GetTransactionResponse `json:"transactions"`
	Total        int                      `json:"total"`
	Page         int                      `json:"page"`
	Limit        int                      `json:"limit"`
	TotalPages   int                      `json:"total_pages"`
}

// ListUsersResponse representa a resposta de listagem de usuários
type ListUsersResponse struct {
	Users      []GetUserResponse `json:"users"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// BalanceResponse representa a resposta de consulta de saldo
type BalanceResponse struct {
	UserID    string `json:"user_id"`
	Balance   string `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details string      `json:"details,omitempty"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

// SuccessResponse representa uma resposta de sucesso padronizada
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginationParams representa os parâmetros de paginação
type PaginationParams struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

// DefaultPagination retorna parâmetros de paginação padrão
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:  1,
		Limit: 20,
	}
}

// Offset calcula o offset para queries de banco
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// TransactionFilters representa filtros para consulta de transações
type TransactionFilters struct {
	PaginationParams
	UserID    string            `json:"user_id,omitempty"`
	Status    TransactionStatus `json:"status,omitempty"`
	DateFrom  *time.Time        `json:"date_from,omitempty"`
	DateTo    *time.Time        `json:"date_to,omitempty"`
	MinAmount *decimal.Decimal  `json:"min_amount,omitempty"`
	MaxAmount *decimal.Decimal  `json:"max_amount,omitempty"`
}

// UserFilters representa filtros para consulta de usuários
type UserFilters struct {
	PaginationParams
	UserType UserType `json:"user_type,omitempty"`
	Email    string   `json:"email,omitempty"`
	Document string   `json:"document,omitempty"`
}
