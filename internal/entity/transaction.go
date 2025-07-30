package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TransactionStatus representa o status de uma transação
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// Transaction representa uma transação financeira entre usuários
type Transaction struct {
	ID                string            `json:"id" db:"id"`
	PayerID           string            `json:"payer_id" db:"payer_id"`
	PayeeID           string            `json:"payee_id" db:"payee_id"`
	Amount            decimal.Decimal   `json:"amount" db:"amount"`
	Status            TransactionStatus `json:"status" db:"status"`
	AuthorizationID   *string           `json:"authorization_id,omitempty" db:"authorization_id"`
	NotificationSent  bool              `json:"notification_sent" db:"notification_sent"`
	FailureReason     *string           `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
	CompletedAt       *time.Time        `json:"completed_at,omitempty" db:"completed_at"`
	
	// Campos relacionais (não persistidos no banco)
	Payer *User `json:"payer,omitempty" db:"-"`
	Payee *User `json:"payee,omitempty" db:"-"`
}

// TransactionRequest representa os dados necessários para criar uma transação
type TransactionRequest struct {
	PayerID string          `json:"payer_id" validate:"required,uuid"`
	PayeeID string          `json:"payee_id" validate:"required,uuid"`
	Amount  decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

// NewTransaction cria uma nova transação
func NewTransaction(payerID, payeeID string, amount decimal.Decimal) (*Transaction, error) {
	transaction := &Transaction{
		ID:               uuid.New().String(),
		PayerID:          payerID,
		PayeeID:          payeeID,
		Amount:           amount,
		Status:           TransactionStatusPending,
		NotificationSent: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Validate valida os dados da transação
func (t *Transaction) Validate() error {
	if t.PayerID == "" {
		return errors.New("pagador é obrigatório")
	}

	if t.PayeeID == "" {
		return errors.New("recebedor é obrigatório")
	}

	if t.PayerID == t.PayeeID {
		return errors.New("pagador e recebedor não podem ser o mesmo usuário")
	}

	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("valor da transação deve ser maior que zero")
	}

	// Limite máximo de transação (exemplo: R$ 10.000,00)
	maxAmount := decimal.NewFromFloat(10000.00)
	if t.Amount.GreaterThan(maxAmount) {
		return errors.New("valor da transação excede o limite máximo")
	}

	return nil
}

// ValidateBusinessRules valida as regras de negócio da transação
func (t *Transaction) ValidateBusinessRules(payer, payee *User) error {
	// Lojistas não podem ser pagadores
	if payer.IsMerchant() {
		return errors.New("lojistas não podem realizar transferências")
	}

	// Verificar se o pagador tem saldo suficiente
	if !payer.HasSufficientBalance(t.Amount) {
		return errors.New("saldo insuficiente")
	}

	// Verificar se os usuários existem e são diferentes
	if payer.ID == payee.ID {
		return errors.New("não é possível transferir para si mesmo")
	}

	return nil
}

// Authorize autoriza a transação
func (t *Transaction) Authorize(authorizationID string) {
	t.Status = TransactionStatusAuthorized
	t.AuthorizationID = &authorizationID
	t.UpdatedAt = time.Now()
}

// Complete completa a transação
func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// Fail marca a transação como falha
func (t *Transaction) Fail(reason string) {
	t.Status = TransactionStatusFailed
	t.FailureReason = &reason
	t.UpdatedAt = time.Now()
}

// Reverse reverte a transação
func (t *Transaction) Reverse(reason string) {
	t.Status = TransactionStatusReversed
	t.FailureReason = &reason
	t.UpdatedAt = time.Now()
}

// MarkNotificationSent marca que a notificação foi enviada
func (t *Transaction) MarkNotificationSent() {
	t.NotificationSent = true
	t.UpdatedAt = time.Now()
}

// CanBeAuthorized verifica se a transação pode ser autorizada
func (t *Transaction) CanBeAuthorized() bool {
	return t.Status == TransactionStatusPending
}

// CanBeCompleted verifica se a transação pode ser completada
func (t *Transaction) CanBeCompleted() bool {
	return t.Status == TransactionStatusAuthorized
}

// CanBeReversed verifica se a transação pode ser revertida
func (t *Transaction) CanBeReversed() bool {
	return t.Status == TransactionStatusCompleted || t.Status == TransactionStatusAuthorized
}

// IsPending verifica se a transação está pendente
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsAuthorized verifica se a transação está autorizada
func (t *Transaction) IsAuthorized() bool {
	return t.Status == TransactionStatusAuthorized
}

// IsCompleted verifica se a transação foi completada
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsFailed verifica se a transação falhou
func (t *Transaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

// IsReversed verifica se a transação foi revertida
func (t *Transaction) IsReversed() bool {
	return t.Status == TransactionStatusReversed
}

// GetStatusDescription retorna a descrição do status em português
func (t *Transaction) GetStatusDescription() string {
	switch t.Status {
	case TransactionStatusPending:
		return "Pendente"
	case TransactionStatusAuthorized:
		return "Autorizada"
	case TransactionStatusCompleted:
		return "Concluída"
	case TransactionStatusFailed:
		return "Falhou"
	case TransactionStatusReversed:
		return "Revertida"
	default:
		return "Status desconhecido"
	}
}

// GetAmountFormatted retorna o valor formatado em reais
func (t *Transaction) GetAmountFormatted() string {
	return "R$ " + t.Amount.StringFixed(2)
}
