package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/monteirobsb/user-management/backend/auth"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/handlers"
	"github.com/monteirobsb/user-management/backend/middleware"
)

// ... (definição do handler de login)
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Payload inválido"})
		return
	}

	token, err := auth.LoginUser(payload.Email, payload.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Aviso: Não foi possível encontrar o arquivo .env, usando variáveis de ambiente do sistema.")
	}

	// Inicializa o banco de dados
	database.InitDatabase()

	// Inicia o roteador Gin
	router := gin.Default()

	// Adiciona o Error Handler globalmente
	router.Use(middleware.ErrorHandler())

	// Agrupa as rotas da API sob o prefixo /api
	api := router.Group("/api")
	{
		// Rotas públicas
		api.POST("/login", LoginHandler)
		api.POST("/users", handlers.CreateUserHandler) // <-- ROTA DE CRIAÇÃO AGORA É PÚBLICA

		// Rotas protegidas
		protected := api.Group("/users").Use(middleware.AuthMiddleware())
		{
			// A rota POST foi removida daqui
			protected.GET("", handlers.GetUsersHandler)
			protected.GET("/:id", handlers.GetUserHandler)
			protected.PUT("/:id", handlers.UpdateUserHandler)
			protected.DELETE("/:id", handlers.DeleteUserHandler)
		}
	}

	// Inicia o servidor na porta definida
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}
	router.Run(":" + port)
}
