package entity

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// UserType representa o tipo de usuário
type UserType string

const (
	UserTypeCommon   UserType = "common"
	UserTypeMerchant UserType = "merchant"
)

// User representa a entidade usuário do sistema
type User struct {
	ID           string          `json:"id" db:"id"`
	FullName     string          `json:"full_name" db:"full_name"`
	Document     string          `json:"document" db:"document"` // CPF ou CNPJ
	Email        string          `json:"email" db:"email"`
	Password     string          `json:"-" db:"password"` // Não retorna no JSON
	UserType     UserType        `json:"user_type" db:"user_type"`
	Balance      decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// NewUser cria uma nova instância de usuário
func NewUser(fullName, document, email, password string, userType UserType) (*User, error) {
	user := &User{
		ID:        uuid.New().String(),
		FullName:  strings.TrimSpace(fullName),
		Document:  cleanDocument(document),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Password:  password,
		UserType:  userType,
		Balance:   decimal.Zero,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate valida os dados do usuário
func (u *User) Validate() error {
	if u.FullName == "" {
		return errors.New("nome completo é obrigatório")
	}

	if len(u.FullName) < 3 {
		return errors.New("nome completo deve ter pelo menos 3 caracteres")
	}

	if !u.IsValidDocument() {
		return errors.New("documento inválido")
	}

	if !u.IsValidEmail() {
		return errors.New("email inválido")
	}

	if u.Password == "" {
		return errors.New("senha é obrigatória")
	}

	if len(u.Password) < 6 {
		return errors.New("senha deve ter pelo menos 6 caracteres")
	}

	if u.UserType != UserTypeCommon && u.UserType != UserTypeMerchant {
		return errors.New("tipo de usuário inválido")
	}

	return nil
}

// IsValidDocument valida CPF ou CNPJ
func (u *User) IsValidDocument() bool {
	if len(u.Document) == 11 {
		return u.isValidCPF()
	}
	if len(u.Document) == 14 {
		return u.isValidCNPJ()
	}
	return false
}

// IsValidEmail valida o formato do email
func (u *User) IsValidEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

// IsMerchant verifica se o usuário é lojista
func (u *User) IsMerchant() bool {
	return u.UserType == UserTypeMerchant
}

// IsCommon verifica se o usuário é comum
func (u *User) IsCommon() bool {
	return u.UserType == UserTypeCommon
}

// CanSendMoney verifica se o usuário pode enviar dinheiro
func (u *User) CanSendMoney() bool {
	// Lojistas não podem ser pagadores
	return u.UserType == UserTypeCommon
}

// HasSufficientBalance verifica se o usuário tem saldo suficiente
func (u *User) HasSufficientBalance(amount decimal.Decimal) bool {
	return u.Balance.GreaterThanOrEqual(amount)
}

// DebitBalance debita valor do saldo
func (u *User) DebitBalance(amount decimal.Decimal) error {
	if !u.HasSufficientBalance(amount) {
		return errors.New("saldo insuficiente")
	}
	u.Balance = u.Balance.Sub(amount)
	u.UpdatedAt = time.Now()
	return nil
}

// CreditBalance credita valor ao saldo
func (u *User) CreditBalance(amount decimal.Decimal) {
	u.Balance = u.Balance.Add(amount)
	u.UpdatedAt = time.Now()
}

// UpdatePassword atualiza a senha do usuário
func (u *User) UpdatePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("senha deve ter pelo menos 6 caracteres")
	}
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return nil
}

// cleanDocument remove caracteres especiais do documento
func cleanDocument(document string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(document, "")
}

// isValidCPF valida CPF usando algoritmo oficial
func (u *User) isValidCPF() bool {
	cpf := u.Document
	
	// Verifica se tem 11 dígitos
	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if regexp.MustCompile(`^(\d)\1{10}$`).MatchString(cpf) {
		return false
	}

	// Valida primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if int(cpf[9]-'0') != firstDigit {
		return false
	}

	// Valida segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit := int(cpf[i] - '0')
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return int(cpf[10]-'0') == secondDigit
}

// isValidCNPJ valida CNPJ usando algoritmo oficial
func (u *User) isValidCNPJ() bool {
	cnpj := u.Document
	
	// Verifica se tem 14 dígitos
	if len(cnpj) != 14 {
		return false
	}

	// Verifica se todos os dígitos são iguais
	if regexp.MustCompile(`^(\d)\1{13}$`).MatchString(cnpj) {
		return false
	}

	// Pesos para validação do primeiro dígito
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weights1[i]
	}
	remainder := sum % 11
	firstDigit := 0
	if remainder >= 2 {
		firstDigit = 11 - remainder
	}

	if int(cnpj[12]-'0') != firstDigit {
		return false
	}

	// Pesos para validação do segundo dígito
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		digit := int(cnpj[i] - '0')
		sum += digit * weights2[i]
	}
	remainder = sum % 11
	secondDigit := 0
	if remainder >= 2 {
		secondDigit = 11 - remainder
	}

	return int(cnpj[13]-'0') == secondDigit
}
