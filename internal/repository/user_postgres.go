package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"payflow-api/internal/entity"
	"payflow-api/pkg/database"
)

type userPostgresRepository struct {
	db *database.Database
}

func NewUserPostgresRepository(db *database.Database) UserRepository {
	return &userPostgresRepository{
		db: db,
	}
}

func (r *userPostgresRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, full_name, document, email, password, user_type, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		user.ID,
		user.FullName,
		user.Document,
		user.Email,
		user.Password,
		user.UserType,
		user.Balance,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return fmt.Errorf("email já cadastrado")
			}
			if strings.Contains(err.Error(), "document") {
				return fmt.Errorf("documento já cadastrado")
			}
		}
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}

func (r *userPostgresRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, full_name, document, email, password, user_type, balance, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Document,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userPostgresRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, full_name, document, email, password, user_type, balance, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	err := r.db.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Document,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userPostgresRepository) GetByDocument(ctx context.Context, document string) (*entity.User, error) {
	query := `
		SELECT id, full_name, document, email, password, user_type, balance, created_at, updated_at
		FROM users
		WHERE document = $1
	`

	user := &entity.User{}
	err := r.db.DB.QueryRowContext(ctx, query, document).Scan(
		&user.ID,
		&user.FullName,
		&user.Document,
		&user.Email,
		&user.Password,
		&user.UserType,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuário não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userPostgresRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users 
		SET full_name = $2, email = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		user.ID,
		user.FullName,
		user.Email,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

func (r *userPostgresRepository) List(ctx context.Context, filters *entity.UserFilters) ([]*entity.User, int, error) {
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	query := `
		SELECT id, full_name, document, email, password, user_type, balance, created_at, updated_at
		FROM users
		WHERE 1=1
	`

	if filters.UserType != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND user_type = $%d", argCount)
		query += fmt.Sprintf(" AND user_type = $%d", argCount)
		args = append(args, filters.UserType)
	}

	if filters.Email != "" {
		argCount++
		countQuery += fmt.Sprintf(" AND email ILIKE $%d", argCount)
		query += fmt.Sprintf(" AND email ILIKE $%d", argCount)
		args = append(args, "%"+filters.Email+"%")
	}

	var total int
	err := r.db.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar usuários: %w", err)
	}

	query += " ORDER BY created_at DESC"
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filters.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filters.Offset())

	rows, err := r.db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao listar usuários: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Document,
			&user.Email,
			&user.Password,
			&user.UserType,
			&user.Balance,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao fazer scan do usuário: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userPostgresRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

func (r *userPostgresRepository) ExistsByEmailOrDocument(ctx context.Context, email, document string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1 OR document = $2"

	var count int
	err := r.db.DB.QueryRowContext(ctx, query, email, document).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar usuário existente: %w", err)
	}

	return count > 0, nil
}
