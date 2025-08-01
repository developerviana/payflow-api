package main

import (
	"log"
	"net/http"

	"payflow-api/internal/config"
	"payflow-api/internal/handler"
	"payflow-api/internal/repository"
	"payflow-api/internal/usecase"
	"payflow-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Erro ao conectar com banco de dados: %v", err)
	}
	defer db.Close()

	// Inicializar camadas
	userRepo := repository.NewUserPostgresRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

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
				"name":        "PayFlow API",
				"description": "API para transferências financeiras",
				"version":     "1.0.0",
				"environment": cfg.Server.Env,
				"database":    "PostgreSQL",
				"port":        cfg.Server.Port,
			})
		})

		// Rotas de usuários
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/:id/balance", userHandler.GetBalance)
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
