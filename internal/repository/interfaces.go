package repository

import (
	"context"
	"payflow-api/internal/entity"
)

// UserRepository define métodos para manipulação de usuários no repositório.
type UserRepository interface {
	// Create insere um novo usuário.
	Create(ctx context.Context, user *entity.User) error
	// GetByID retorna um usuário pelo ID.
	GetByID(ctx context.Context, id string) (*entity.User, error)
	// GetByEmail retorna um usuário pelo e-mail.
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	// GetByDocument retorna um usuário pelo documento.
	GetByDocument(ctx context.Context, document string) (*entity.User, error)
	// Update atualiza os dados de um usuário.
	Update(ctx context.Context, user *entity.User) error
	// List retorna uma lista de usuários com filtros e total.
	List(ctx context.Context, filters *entity.UserFilters) ([]*entity.User, int, error)
	// Delete remove um usuário pelo ID.
	Delete(ctx context.Context, id string) error
	// ExistsByEmailOrDocument verifica se existe usuário por e-mail ou documento.
	ExistsByEmailOrDocument(ctx context.Context, email, document string) (bool, error)
}
