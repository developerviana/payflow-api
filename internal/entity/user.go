package entity

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserType string

const (
	UserTypeCommon   UserType = "common"
	UserTypeMerchant UserType = "merchant"
)

type User struct {
	ID        string          `json:"id" db:"id"`
	FullName  string          `json:"full_name" db:"full_name"`
	Document  string          `json:"document" db:"document"`
	Email     string          `json:"email" db:"email"`
	Password  string          `json:"-" db:"password"`
	UserType  UserType        `json:"user_type" db:"user_type"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

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

func (u *User) IsValidDocument() bool {
	if len(u.Document) == 11 {
		return u.isValidCPF()
	}
	if len(u.Document) == 14 {
		return u.isValidCNPJ()
	}
	return false
}

func (u *User) IsValidEmail() bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

func (u *User) IsMerchant() bool {
	return u.UserType == UserTypeMerchant
}

func (u *User) IsCommon() bool {
	return u.UserType == UserTypeCommon
}

func (u *User) CanSendMoney() bool {
	// Lojistas não podem ser pagadores
	return u.UserType == UserTypeCommon
}

func (u *User) HasSufficientBalance(amount decimal.Decimal) bool {
	return u.Balance.GreaterThanOrEqual(amount)
}

func (u *User) DebitBalance(amount decimal.Decimal) error {
	if !u.HasSufficientBalance(amount) {
		return errors.New("saldo insuficiente")
	}
	u.Balance = u.Balance.Sub(amount)
	u.UpdatedAt = time.Now()
	return nil
}
func (u *User) CreditBalance(amount decimal.Decimal) {
	u.Balance = u.Balance.Add(amount)
	u.UpdatedAt = time.Now()
}

func (u *User) UpdatePassword(newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("senha deve ter pelo menos 6 caracteres")
	}
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return nil
}

func cleanDocument(document string) string {
	re := regexp.MustCompile(`[^0-9]`)
	return re.ReplaceAllString(document, "")
}

func (u *User) isValidCPF() bool {
	cpf := u.Document

	if len(cpf) != 11 {
		return false
	}

	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

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

func (u *User) isValidCNPJ() bool {
	cnpj := u.Document

	if len(cnpj) != 14 {
		return false
	}

	allSame := true
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != cnpj[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

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
