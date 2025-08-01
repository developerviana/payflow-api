package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type CreateUserRequest struct {
	FullName string   `json:"full_name" validate:"required,min=3,max=100"`
	Document string   `json:"document" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	UserType UserType `json:"user_type" validate:"required,oneof=common merchant"`
}

type CreateUserResponse struct {
	ID       string   `json:"id"`
	FullName string   `json:"full_name"`
	Document string   `json:"document"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
	Balance  string   `json:"balance"`
}

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

type UpdateUserRequest struct {
	FullName string `json:"full_name,omitempty" validate:"omitempty,min=3,max=100"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type CreateTransactionRequest struct {
	PayeeID string          `json:"payee_id" validate:"required,uuid"`
	Amount  decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

type CreateTransactionResponse struct {
	ID         string            `json:"id"`
	PayerID    string            `json:"payer_id"`
	PayeeID    string            `json:"payee_id"`
	Amount     string            `json:"amount"`
	Status     TransactionStatus `json:"status"`
	StatusDesc string            `json:"status_description"`
	CreatedAt  time.Time         `json:"created_at"`
}

type GetTransactionResponse struct {
	ID               string            `json:"id"`
	PayerID          string            `json:"payer_id"`
	PayeeID          string            `json:"payee_id"`
	Amount           string            `json:"amount"`
	Status           TransactionStatus `json:"status"`
	StatusDesc       string            `json:"status_description"`
	AuthorizationID  *string           `json:"authorization_id,omitempty"`
	NotificationSent bool              `json:"notification_sent"`
	FailureReason    *string           `json:"failure_reason,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	CompletedAt      *time.Time        `json:"completed_at,omitempty"`

	Payer *UserSummary `json:"payer,omitempty"`
	Payee *UserSummary `json:"payee,omitempty"`
}

type UserSummary struct {
	ID       string   `json:"id"`
	FullName string   `json:"full_name"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
}

type ListTransactionsResponse struct {
	Transactions []GetTransactionResponse `json:"transactions"`
	Total        int                      `json:"total"`
	Page         int                      `json:"page"`
	Limit        int                      `json:"limit"`
	TotalPages   int                      `json:"total_pages"`
}

type ListUsersResponse struct {
	Users      []GetUserResponse `json:"users"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

type BalanceResponse struct {
	UserID    string    `json:"user_id"`
	Balance   string    `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details string      `json:"details,omitempty"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationParams struct {
	Page  int `json:"page" validate:"min=1"`
	Limit int `json:"limit" validate:"min=1,max=100"`
}

func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:  1,
		Limit: 20,
	}
}

func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

type TransactionFilters struct {
	PaginationParams
	UserID    string            `json:"user_id,omitempty"`
	Status    TransactionStatus `json:"status,omitempty"`
	DateFrom  *time.Time        `json:"date_from,omitempty"`
	DateTo    *time.Time        `json:"date_to,omitempty"`
	MinAmount *decimal.Decimal  `json:"min_amount,omitempty"`
	MaxAmount *decimal.Decimal  `json:"max_amount,omitempty"`
}

type UserFilters struct {
	PaginationParams
	UserType UserType `json:"user_type,omitempty"`
	Email    string   `json:"email,omitempty"`
	Document string   `json:"document,omitempty"`
}
