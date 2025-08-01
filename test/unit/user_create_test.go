package entity_test

import (
	"payflow-api/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewUser(t *testing.T) *entity.User {
	fullName := "João Silva"
	document := "11144477735"
	email := "joao.silva@example.com"
	password := "SenhaForte123"
	userType := entity.UserTypeCommon

	user, err := entity.NewUser(fullName, document, email, password, userType)
	if err != nil {
		t.Fatalf("erro ao criar usuário válido: %v", err)
	}
	return user
}

func TestCreateValidUser(t *testing.T) {
	user := NewUser(t)

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "João Silva", user.FullName)
	assert.Equal(t, "11144477735", user.Document)
	assert.Equal(t, "joao.silva@example.com", user.Email)
	assert.Equal(t, entity.UserTypeCommon, user.UserType)
	assert.True(t, user.Balance.IsZero())
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestCreateUserWithInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		document string
		email    string
		password string
		userType entity.UserType
	}{
		{
			name:     "Nome vazio",
			fullName: "",
			document: "11144477735",
			email:    "test@example.com",
			password: "123456",
			userType: entity.UserTypeCommon,
		},
		{
			name:     "CPF inválido",
			fullName: "João Silva",
			document: "11111111111",
			email:    "test@example.com",
			password: "123456",
			userType: entity.UserTypeCommon,
		},
		{
			name:     "Email inválido",
			fullName: "João Silva",
			document: "11144477735",
			email:    "email-inválido",
			password: "123456",
			userType: entity.UserTypeCommon,
		},
		{
			name:     "Senha muito curta",
			fullName: "João Silva",
			document: "11144477735",
			email:    "test@example.com",
			password: "123",
			userType: entity.UserTypeCommon,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.fullName, tt.document, tt.email, tt.password, tt.userType)

			assert.Error(t, err)
			assert.Nil(t, user)
		})
	}
}

func TestCreateMerchantUser(t *testing.T) {
	fullName := "Loja do João LTDA"
	document := "11222333000181"
	email := "loja@example.com"
	password := "SenhaForte123"
	userType := entity.UserTypeMerchant

	user, err := entity.NewUser(fullName, document, email, password, userType)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, entity.UserTypeMerchant, user.UserType)
	assert.True(t, user.IsMerchant())
	assert.False(t, user.CanSendMoney())
}
