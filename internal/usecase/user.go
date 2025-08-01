package usecase

import (
	"context"
	"fmt"
	"payflow-api/internal/entity"
	"payflow-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// UserUseCase define as operações de negócio para usuários
type UserUseCase interface {
	CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error)
	GetUser(ctx context.Context, id string) (*entity.GetUserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.GetUserResponse, error)
	UpdateUser(ctx context.Context, id string, req *entity.UpdateUserRequest) (*entity.GetUserResponse, error)
	ListUsers(ctx context.Context, filters *entity.UserFilters) (*entity.ListUsersResponse, error)
	DeleteUser(ctx context.Context, id string) error
	GetBalance(ctx context.Context, id string) (*entity.BalanceResponse, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase cria uma nova instância do use case
func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

// CreateUser cria um novo usuário
func (uc *userUseCase) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.CreateUserResponse, error) {
	// Validar se já existe usuário com mesmo email ou documento
	exists, err := uc.userRepo.ExistsByEmailOrDocument(ctx, req.Email, req.Document)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar usuário existente: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("já existe um usuário com este email ou documento")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar senha: %w", err)
	}

	// Criar usuário usando a entidade (com todas as validações)
	user, err := entity.FromCreateUserRequest(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao validar dados do usuário: %w", err)
	}

	// Atualizar senha criptografada
	user.Password = string(hashedPassword)

	// Salvar no banco
	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("erro ao salvar usuário: %w", err)
	}

	// Retornar resposta
	return user.ToCreateUserResponse(), nil
}

// GetUser busca um usuário por ID
func (uc *userUseCase) GetUser(ctx context.Context, id string) (*entity.GetUserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user.ToGetUserResponse(), nil
}

// GetUserByEmail busca um usuário por email
func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.GetUserResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user.ToGetUserResponse(), nil
}

// UpdateUser atualiza os dados de um usuário
func (uc *userUseCase) UpdateUser(ctx context.Context, id string, req *entity.UpdateUserRequest) (*entity.GetUserResponse, error) {
	// Buscar usuário existente
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aplicar alterações
	err = user.ApplyUpdateUserRequest(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao validar dados de atualização: %w", err)
	}

	// Salvar alterações
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user.ToGetUserResponse(), nil
}

// ListUsers lista usuários com paginação e filtros
func (uc *userUseCase) ListUsers(ctx context.Context, filters *entity.UserFilters) (*entity.ListUsersResponse, error) {
	// Validar paginação
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.Limit == 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	users, total, err := uc.userRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Converter para responses
	var userResponses []entity.GetUserResponse
	for _, user := range users {
		userResponses = append(userResponses, *user.ToGetUserResponse())
	}

	// Calcular total de páginas
	totalPages := (total + filters.Limit - 1) / filters.Limit

	return &entity.ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteUser remove um usuário
func (uc *userUseCase) DeleteUser(ctx context.Context, id string) error {
	return uc.userRepo.Delete(ctx, id)
}

// GetBalance retorna o saldo de um usuário
func (uc *userUseCase) GetBalance(ctx context.Context, id string) (*entity.BalanceResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user.ToBalanceResponse(), nil
}
