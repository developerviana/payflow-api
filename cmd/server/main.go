package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"payflow-api/internal/config"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Configurar Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware básico
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS simples para desenvolvimento
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Rotas básicas para teste
	v1 := router.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "PayFlow API está rodando!",
				"version": "1.0.0",
			})
		})

		// Info da API
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name":         "PayFlow API",
				"description":  "API para transferências financeiras",
				"version":      "1.0.0",
				"environment":  cfg.Server.Env,
				"database":     "PostgreSQL",
				"port":         cfg.Server.Port,
			})
		})

		// Rotas de usuários (placeholder)
		users := v1.Group("/users")
		{
			users.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Endpoint de usuários - Em desenvolvimento",
					"routes": []string{
						"GET /api/v1/users - Listar usuários",
						"POST /api/v1/users - Criar usuário",
						"GET /api/v1/users/:id - Buscar usuário",
						"PUT /api/v1/users/:id - Atualizar usuário",
					},
				})
			})
		}

		// Rotas de transações (placeholder)
		transactions := v1.Group("/transactions")
		{
			transactions.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Endpoint de transações - Em desenvolvimento",
					"routes": []string{
						"GET /api/v1/transactions - Listar transações",
						"POST /api/v1/transactions - Criar transação",
						"GET /api/v1/transactions/:id - Buscar transação",
					},
				})
			})
		}
	}

	// Iniciar servidor
	log.Printf("🚀 Servidor iniciando na porta %s", cfg.Server.Port)
	log.Printf("📚 Documentação disponível em: http://localhost:%s/api/v1/info", cfg.Server.Port)
	log.Printf("💚 Health check em: http://localhost:%s/api/v1/health", cfg.Server.Port)
	
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
