package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	External ExternalConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ExternalConfig struct {
	AuthorizerURL   string
	NotificationURL string
	RequestTimeout  int
}

func Load() (*Config, error) {
	// Carrega variáveis de ambiente do arquivo .env se existir
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "payflow"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		External: ExternalConfig{
			AuthorizerURL:   getEnv("AUTHORIZER_URL", "https://util.devi.tools/api/v2/authorize"),
			NotificationURL: getEnv("NOTIFICATION_URL", "https://util.devi.tools/api/v1/notify"),
			RequestTimeout:  getEnvAsInt("REQUEST_TIMEOUT", 10),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
