package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusAuthorized TransactionStatus = "authorized"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusReversed   TransactionStatus = "reversed"
)

type Transaction struct {
	ID               string            `json:"id" db:"id"`
	PayerID          string            `json:"payer_id" db:"payer_id"`
	PayeeID          string            `json:"payee_id" db:"payee_id"`
	Amount           decimal.Decimal   `json:"amount" db:"amount"`
	Status           TransactionStatus `json:"status" db:"status"`
	AuthorizationID  *string           `json:"authorization_id,omitempty" db:"authorization_id"`
	NotificationSent bool              `json:"notification_sent" db:"notification_sent"`
	FailureReason    *string           `json:"failure_reason,omitempty" db:"failure_reason"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
	CompletedAt      *time.Time        `json:"completed_at,omitempty" db:"completed_at"`

	Payer *User `json:"payer,omitempty" db:"-"`
	Payee *User `json:"payee,omitempty" db:"-"`
}

type TransactionRequest struct {
	PayerID string          `json:"payer_id" validate:"required,uuid"`
	PayeeID string          `json:"payee_id" validate:"required,uuid"`
	Amount  decimal.Decimal `json:"amount" validate:"required,gt=0"`
}

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

	maxAmount := decimal.NewFromFloat(10000.00)
	if t.Amount.GreaterThan(maxAmount) {
		return errors.New("valor da transação excede o limite máximo")
	}

	return nil
}

func (t *Transaction) ValidateBusinessRules(payer, payee *User) error {
	// Lojistas não podem ser pagadores
	if payer.IsMerchant() {
		return errors.New("lojistas não podem realizar transferências")
	}

	if !payer.HasSufficientBalance(t.Amount) {
		return errors.New("saldo insuficiente")
	}

	if payer.ID == payee.ID {
		return errors.New("não é possível transferir para si mesmo")
	}

	return nil
}

func (t *Transaction) Authorize(authorizationID string) {
	t.Status = TransactionStatusAuthorized
	t.AuthorizationID = &authorizationID
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

func (t *Transaction) Fail(reason string) {
	t.Status = TransactionStatusFailed
	t.FailureReason = &reason
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Reverse(reason string) {
	t.Status = TransactionStatusReversed
	t.FailureReason = &reason
	t.UpdatedAt = time.Now()
}

func (t *Transaction) MarkNotificationSent() {
	t.NotificationSent = true
	t.UpdatedAt = time.Now()
}

func (t *Transaction) CanBeAuthorized() bool {
	return t.Status == TransactionStatusPending
}

func (t *Transaction) CanBeCompleted() bool {
	return t.Status == TransactionStatusAuthorized
}

func (t *Transaction) CanBeReversed() bool {
	return t.Status == TransactionStatusCompleted || t.Status == TransactionStatusAuthorized
}

func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

func (t *Transaction) IsAuthorized() bool {
	return t.Status == TransactionStatusAuthorized
}

func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

func (t *Transaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

func (t *Transaction) IsReversed() bool {
	return t.Status == TransactionStatusReversed
}

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

func (t *Transaction) GetAmountFormatted() string {
	return "R$ " + t.Amount.StringFixed(2)
}
