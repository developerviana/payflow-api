package entity_test

import (
	"payflow-api/internal/entity"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		fullName    string
		document    string
		email       string
		password    string
		userType    entity.UserType
		expectError bool
	}{
		{
			name:        "Valid common user with CPF",
			fullName:    "João Silva",
			document:    "11144477735", // CPF válido
			email:       "joao@teste.com",
			password:    "123456",
			userType:    entity.UserTypeCommon,
			expectError: false,
		},
		{
			name:        "Valid merchant user with CNPJ",
			fullName:    "Empresa LTDA",
			document:    "11222333000181", // CNPJ válido
			email:       "empresa@teste.com",
			password:    "123456",
			userType:    entity.UserTypeMerchant,
			expectError: false,
		},
		{
			name:        "Invalid user - empty name",
			fullName:    "",
			document:    "11144477735",
			email:       "joao@teste.com",
			password:    "123456",
			userType:    entity.UserTypeCommon,
			expectError: true,
		},
		{
			name:        "Invalid user - short name",
			fullName:    "Jo",
			document:    "11144477735",
			email:       "joao@teste.com",
			password:    "123456",
			userType:    entity.UserTypeCommon,
			expectError: true,
		},
		{
			name:        "Invalid user - invalid CPF",
			fullName:    "João Silva",
			document:    "11111111111",
			email:       "joao@teste.com",
			password:    "123456",
			userType:    entity.UserTypeCommon,
			expectError: true,
		},
		{
			name:        "Invalid user - invalid email",
			fullName:    "João Silva",
			document:    "11144477735",
			email:       "invalid-email",
			password:    "123456",
			userType:    entity.UserTypeCommon,
			expectError: true,
		},
		{
			name:        "Invalid user - short password",
			fullName:    "João Silva",
			document:    "11144477735",
			email:       "joao@teste.com",
			password:    "123",
			userType:    entity.UserTypeCommon,
			expectError: true,
		},
		{
			name:        "Invalid user - invalid user type",
			fullName:    "João Silva",
			document:    "11144477735",
			email:       "joao@teste.com",
			password:    "123456",
			userType:    "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.fullName, tt.document, tt.email, tt.password, tt.userType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEmpty(t, user.ID)
				assert.Equal(t, tt.fullName, user.FullName)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.userType, user.UserType)
				assert.True(t, user.Balance.IsZero())
			}
		})
	}
}

func TestUser_CanSendMoney(t *testing.T) {
	commonUser, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)
	merchantUser, _ := entity.NewUser("Empresa LTDA", "11222333000181", "empresa@teste.com", "123456", entity.UserTypeMerchant)

	assert.True(t, commonUser.CanSendMoney())
	assert.False(t, merchantUser.CanSendMoney())
}

func TestUser_HasSufficientBalance(t *testing.T) {
	user, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)

	// Usuário sem saldo
	assert.False(t, user.HasSufficientBalance(decimal.NewFromFloat(100.00)))

	// Adicionar saldo
	user.CreditBalance(decimal.NewFromFloat(500.00))
	assert.True(t, user.HasSufficientBalance(decimal.NewFromFloat(100.00)))
	assert.True(t, user.HasSufficientBalance(decimal.NewFromFloat(500.00)))
	assert.False(t, user.HasSufficientBalance(decimal.NewFromFloat(600.00)))
}

func TestUser_DebitBalance(t *testing.T) {
	user, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)
	user.CreditBalance(decimal.NewFromFloat(500.00))

	// Débito válido
	err := user.DebitBalance(decimal.NewFromFloat(100.00))
	assert.NoError(t, err)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(400.00)))

	// Débito inválido (saldo insuficiente)
	err = user.DebitBalance(decimal.NewFromFloat(500.00))
	assert.Error(t, err)
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(400.00))) // Saldo não alterado
}

func TestUser_CreditBalance(t *testing.T) {
	user, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)

	user.CreditBalance(decimal.NewFromFloat(250.00))
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(250.00)))

	user.CreditBalance(decimal.NewFromFloat(150.00))
	assert.True(t, user.Balance.Equal(decimal.NewFromFloat(400.00)))
}

func TestUser_ValidateCPF(t *testing.T) {
	tests := []struct {
		cpf   string
		valid bool
	}{
		{"11144477735", true},  // CPF válido
		{"11111111111", false}, // CPF inválido (todos os dígitos iguais)
		{"12345678901", false}, // CPF inválido
		{"", false},            // CPF vazio
		{"123", false},         // CPF muito curto
	}

	for _, tt := range tests {
		t.Run(tt.cpf, func(t *testing.T) {
			user := &entity.User{Document: tt.cpf}
			result := user.IsValidDocument()
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestUser_ValidateCNPJ(t *testing.T) {
	tests := []struct {
		cnpj  string
		valid bool
	}{
		{"11222333000181", true},  // CNPJ válido
		{"11111111111111", false}, // CNPJ inválido (todos os dígitos iguais)
		{"12345678000195", false}, // CNPJ inválido
		{"", false},               // CNPJ vazio
		{"123", false},            // CNPJ muito curto
	}

	for _, tt := range tests {
		t.Run(tt.cnpj, func(t *testing.T) {
			user := &entity.User{Document: tt.cnpj}
			result := user.IsValidDocument()
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestUser_IsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"joao@teste.com", true},
		{"usuario@dominio.com.br", true},
		{"test@example.org", true},
		{"invalid-email", false},
		{"@dominio.com", false},
		{"usuario@", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			user := &entity.User{Email: tt.email}
			result := user.IsValidEmail()
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)

	// Senha válida
	err := user.UpdatePassword("novaSenha123")
	assert.NoError(t, err)
	assert.Equal(t, "novaSenha123", user.Password)

	// Senha muito curta
	err = user.UpdatePassword("123")
	assert.Error(t, err)
}

func TestUser_ToCreateUserResponse(t *testing.T) {
	user, _ := entity.NewUser("João Silva", "11144477735", "joao@teste.com", "123456", entity.UserTypeCommon)
	response := user.ToCreateUserResponse()

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.FullName, response.FullName)
	assert.Equal(t, "111.***.***-35", response.Document) // Documento mascarado
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.UserType, response.UserType)
	assert.Equal(t, "0.00", response.Balance)
}
