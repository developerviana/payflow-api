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

// NewUserPostgresRepository cria uma nova instância do repositório PostgreSQL
func NewUserPostgresRepository(db *database.Database) UserRepository {
	return &userPostgresRepository{
		db: db,
	}
}

// Create cria um novo usuário no banco de dados
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

// GetByID busca um usuário por ID
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

// GetByEmail busca um usuário por email
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

// GetByDocument busca um usuário por documento
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

// Update atualiza os dados de um usuário
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

// List lista usuários com paginação e filtros
func (r *userPostgresRepository) List(ctx context.Context, filters *entity.UserFilters) ([]*entity.User, int, error) {
	// Query para contar total
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argCount := 0

	// Query principal
	query := `
		SELECT id, full_name, document, email, password, user_type, balance, created_at, updated_at
		FROM users
		WHERE 1=1
	`

	// Aplicar filtros
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

	// Contar total
	var total int
	err := r.db.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar usuários: %w", err)
	}

	// Adicionar paginação
	query += " ORDER BY created_at DESC"
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filters.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filters.Offset())

	// Executar query
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

// Delete remove um usuário (soft delete)
func (r *userPostgresRepository) Delete(ctx context.Context, id string) error {
	// Por enquanto, implementar como hard delete
	// Em produção, seria melhor implementar soft delete
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

// ExistsByEmailOrDocument verifica se um usuário existe por email ou documento
func (r *userPostgresRepository) ExistsByEmailOrDocument(ctx context.Context, email, document string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1 OR document = $2"

	var count int
	err := r.db.DB.QueryRowContext(ctx, query, email, document).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar usuário existente: %w", err)
	}

	return count > 0, nil
}
