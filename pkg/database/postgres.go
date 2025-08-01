package database

import (
	"database/sql"
	"fmt"
	"time"

	"payflow-api/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

// NewPostgresConnection cria uma nova conexão com PostgreSQL
func NewPostgresConnection(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Log da conexão (sem mostrar a senha)
	fmt.Printf("🔗 Conectando ao PostgreSQL: host=%s port=%s user=%s dbname=%s\n",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com PostgreSQL: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão com PostgreSQL: %w", err)
	}

	return &Database{DB: db}, nil
}

// Close fecha a conexão com o banco
func (d *Database) Close() error {
	return d.DB.Close()
}

// BeginTx inicia uma transação
func (d *Database) BeginTx() (*sql.Tx, error) {
	return d.DB.Begin()
}
