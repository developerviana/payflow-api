package entity

import (
	"testing"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"payflow-api/internal/entity"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name        string
		payerID     string
		payeeID     string
		amount      decimal.Decimal
		expectError bool
	}{
		{
			name:        "Valid transaction",
			payerID:     "payer-123",
			payeeID:     "payee-456",
			amount:      decimal.NewFromFloat(100.00),
			expectError: false,
		},
		{
			name:        "Invalid transaction - same payer and payee",
			payerID:     "user-123",
			payeeID:     "user-123",
			amount:      decimal.NewFromFloat(100.00),
			expectError: true,
		},
		{
			name:        "Invalid transaction - zero amount",
			payerID:     "payer-123",
			payeeID:     "payee-456",
			amount:      decimal.Zero,
			expectError: true,
		},
		{
			name:        "Invalid transaction - negative amount",
			payerID:     "payer-123",
			payeeID:     "payee-456",
			amount:      decimal.NewFromFloat(-50.00),
			expectError: true,
		},
		{
			name:        "Invalid transaction - amount exceeds limit",
			payerID:     "payer-123",
			payeeID:     "payee-456",
			amount:      decimal.NewFromFloat(15000.00),
			expectError: true,
		},
		{
			name:        "Invalid transaction - empty payer",
			payerID:     "",
			payeeID:     "payee-456",
			amount:      decimal.NewFromFloat(100.00),
			expectError: true,
		},
		{
			name:        "Invalid transaction - empty payee",
			payerID:     "payer-123",
			payeeID:     "",
			amount:      decimal.NewFromFloat(100.00),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := entity.NewTransaction(tt.payerID, tt.payeeID, tt.amount)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, transaction)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transaction)
				assert.NotEmpty(t, transaction.ID)
				assert.Equal(t, tt.payerID, transaction.PayerID)
				assert.Equal(t, tt.payeeID, transaction.PayeeID)
				assert.True(t, transaction.Amount.Equal(tt.amount))
				assert.Equal(t, entity.TransactionStatusPending, transaction.Status)
				assert.False(t, transaction.NotificationSent)
			}
		})
	}
}

func TestTransaction_ValidateBusinessRules(t *testing.T) {
	// Criar usuários para teste
	commonUser, _ := NewUser("João Silva", "11144477735", "joao@teste.com", "123456", UserTypeCommon)
	merchantUser, _ := NewUser("Empresa LTDA", "11222333000181", "empresa@teste.com", "123456", UserTypeMerchant)
	commonUser2, _ := NewUser("Maria Silva", "22244477735", "maria@teste.com", "123456", UserTypeCommon)
	
	// Adicionar saldo ao usuário comum
	commonUser.CreditBalance(decimal.NewFromFloat(500.00))
	
	tests := []struct {
		name        string
		transaction *Transaction
		payer       *User
		payee       *User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid transaction between common users",
			transaction: &Transaction{
				PayerID: commonUser.ID,
				PayeeID: commonUser2.ID,
				Amount:  decimal.NewFromFloat(100.00),
			},
			payer:       commonUser,
			payee:       commonUser2,
			expectError: false,
		},
		{
			name: "Valid transaction from common to merchant",
			transaction: &Transaction{
				PayerID: commonUser.ID,
				PayeeID: merchantUser.ID,
				Amount:  decimal.NewFromFloat(100.00),
			},
			payer:       commonUser,
			payee:       merchantUser,
			expectError: false,
		},
		{
			name: "Invalid transaction - merchant as payer",
			transaction: &Transaction{
				PayerID: merchantUser.ID,
				PayeeID: commonUser.ID,
				Amount:  decimal.NewFromFloat(100.00),
			},
			payer:       merchantUser,
			payee:       commonUser,
			expectError: true,
			errorMsg:    "lojistas não podem realizar transferências",
		},
		{
			name: "Invalid transaction - insufficient balance",
			transaction: &Transaction{
				PayerID: commonUser.ID,
				PayeeID: commonUser2.ID,
				Amount:  decimal.NewFromFloat(1000.00),
			},
			payer:       commonUser,
			payee:       commonUser2,
			expectError: true,
			errorMsg:    "saldo insuficiente",
		},
		{
			name: "Invalid transaction - same user",
			transaction: &Transaction{
				PayerID: commonUser.ID,
				PayeeID: commonUser.ID,
				Amount:  decimal.NewFromFloat(100.00),
			},
			payer:       commonUser,
			payee:       commonUser,
			expectError: true,
			errorMsg:    "não é possível transferir para si mesmo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.ValidateBusinessRules(tt.payer, tt.payee)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransaction_StatusMethods(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	
	// Estado inicial
	assert.True(t, transaction.IsPending())
	assert.False(t, transaction.IsAuthorized())
	assert.False(t, transaction.IsCompleted())
	assert.False(t, transaction.IsFailed())
	assert.False(t, transaction.IsReversed())
	
	// Autorizar transação
	transaction.Authorize("auth-123")
	assert.False(t, transaction.IsPending())
	assert.True(t, transaction.IsAuthorized())
	assert.Equal(t, TransactionStatusAuthorized, transaction.Status)
	assert.Equal(t, "auth-123", *transaction.AuthorizationID)
	
	// Completar transação
	transaction.Complete()
	assert.False(t, transaction.IsAuthorized())
	assert.True(t, transaction.IsCompleted())
	assert.Equal(t, TransactionStatusCompleted, transaction.Status)
	assert.NotNil(t, transaction.CompletedAt)
}

func TestTransaction_FailAndReverse(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	
	// Falhar transação
	transaction.Fail("Autorização negada")
	assert.True(t, transaction.IsFailed())
	assert.Equal(t, TransactionStatusFailed, transaction.Status)
	assert.Equal(t, "Autorização negada", *transaction.FailureReason)
	
	// Criar nova transação para testar reversão
	transaction2, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	transaction2.Authorize("auth-123")
	transaction2.Complete()
	
	// Reverter transação
	transaction2.Reverse("Chargeback solicitado")
	assert.True(t, transaction2.IsReversed())
	assert.Equal(t, TransactionStatusReversed, transaction2.Status)
	assert.Equal(t, "Chargeback solicitado", *transaction2.FailureReason)
}

func TestTransaction_CanBeMethods(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	
	// Estado inicial - pendente
	assert.True(t, transaction.CanBeAuthorized())
	assert.False(t, transaction.CanBeCompleted())
	assert.False(t, transaction.CanBeReversed())
	
	// Autorizada
	transaction.Authorize("auth-123")
	assert.False(t, transaction.CanBeAuthorized())
	assert.True(t, transaction.CanBeCompleted())
	assert.True(t, transaction.CanBeReversed())
	
	// Completada
	transaction.Complete()
	assert.False(t, transaction.CanBeAuthorized())
	assert.False(t, transaction.CanBeCompleted())
	assert.True(t, transaction.CanBeReversed())
	
	// Falhada
	transaction2, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	transaction2.Fail("Erro")
	assert.False(t, transaction2.CanBeAuthorized())
	assert.False(t, transaction2.CanBeCompleted())
	assert.False(t, transaction2.CanBeReversed())
}

func TestTransaction_MarkNotificationSent(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	
	assert.False(t, transaction.NotificationSent)
	
	transaction.MarkNotificationSent()
	assert.True(t, transaction.NotificationSent)
}

func TestTransaction_GetStatusDescription(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	
	tests := []struct {
		status      TransactionStatus
		description string
	}{
		{TransactionStatusPending, "Pendente"},
		{TransactionStatusAuthorized, "Autorizada"},
		{TransactionStatusCompleted, "Concluída"},
		{TransactionStatusFailed, "Falhou"},
		{TransactionStatusReversed, "Revertida"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			transaction.Status = tt.status
			assert.Equal(t, tt.description, transaction.GetStatusDescription())
		})
	}
}

func TestTransaction_GetAmountFormatted(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(123.45))
	assert.Equal(t, "R$ 123.45", transaction.GetAmountFormatted())
	
	transaction2, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(1000.00))
	assert.Equal(t, "R$ 1000.00", transaction2.GetAmountFormatted())
}

func TestTransaction_ToCreateTransactionResponse(t *testing.T) {
	transaction, _ := NewTransaction("payer-123", "payee-456", decimal.NewFromFloat(100.00))
	response := transaction.ToCreateTransactionResponse()
	
	assert.Equal(t, transaction.ID, response.ID)
	assert.Equal(t, transaction.PayerID, response.PayerID)
	assert.Equal(t, transaction.PayeeID, response.PayeeID)
	assert.Equal(t, "100.00", response.Amount)
	assert.Equal(t, transaction.Status, response.Status)
	assert.Equal(t, "Pendente", response.StatusDesc)
	assert.Equal(t, transaction.CreatedAt, response.CreatedAt)
}
