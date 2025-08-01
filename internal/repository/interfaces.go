package repository

import (
	"context"
	"payflow-api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByDocument(ctx context.Context, document string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, filters *entity.UserFilters) ([]*entity.User, int, error)
	Delete(ctx context.Context, id string) error
	ExistsByEmailOrDocument(ctx context.Context, email, document string) (bool, error)
}
