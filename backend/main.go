package main

import (
	"log"

	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv" // Removed, as config now handles .env loading
	"github.com/monteirobsb/user-management/backend/auth"
	"github.com/monteirobsb/user-management/backend/config" // Import the new config package
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/handlers"
	"github.com/monteirobsb/user-management/backend/middleware"
)

// LoginPayload define a estrutura esperada para o corpo da requisição de login.
type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler processa as requisições de login.
func LoginHandler(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		// O ErrorHandler middleware pode capturar isso se c.Error(err) for usado,
		// ou podemos retornar um JSON específico aqui.
		// Para consistência com validações de UserCreateRequest, o erro de ShouldBindJSON é informativo.
		c.JSON(400, gin.H{"error": "Payload inválido ou dados ausentes", "details": err.Error()})
		return
	}

	token, err := auth.LoginUser(payload.Email, payload.Password)
	if err != nil {
		// auth.LoginUser já loga os erros internos.
		// Retorna uma mensagem genérica para o cliente.
		c.JSON(401, gin.H{"error": err.Error()}) // err.Error() de LoginUser é "usuário não encontrado..." ou similar
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func main() {
	// godotenv.Load is now called within config.LoadConfig()
	// if err := godotenv.Load("../.env"); err != nil {
	// 	log.Print("WARN: Não foi possível carregar o arquivo .env. Usando variáveis de ambiente do sistema, se definidas.")
	// } else {
	// 	log.Print("INFO: Arquivo .env carregado com sucesso.")
	// }

	// Inicializa o banco de dados.
	// InitDatabase agora usa log.Fatal em caso de erro, então não precisamos checar erro aqui.
	// Config is loaded by its init function, so database.InitDatabase() can use it.
	database.InitDatabase()

	// Inicia o roteador Gin
	// gin.Default() já vem com os middlewares Logger e Recovery.
	router := gin.Default()

	// Adiciona o Error Handler globalmente. Este deve vir depois de outros middlewares
	// se quisermos que ele capture erros deles, ou no início se for para uso geral.
	// Para capturar erros de rota/handler, esta posição é boa.
	router.Use(middleware.ErrorHandler())

	// Agrupa as rotas da API sob o prefixo /api
	api := router.Group("/api")
	{
		// Rotas públicas
		api.POST("/login", LoginHandler)
		// A rota de criação de usuário deve ser pública para permitir o registro de novos usuários.
		api.POST("/users", handlers.CreateUserHandler)

		// Rotas protegidas
		// O middleware AuthMiddleware() será aplicado a este grupo.
		protected := api.Group("/users")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("", handlers.GetUsersHandler)
			protected.GET("/:id", handlers.GetUserHandler)
			protected.PUT("/:id", handlers.UpdateUserHandler)
			protected.DELETE("/:id", handlers.DeleteUserHandler)
		}
	}

	// Inicia o servidor na porta definida
	// Use APIPort from the config package
	log.Printf("INFO: Servidor Gin iniciando na porta :%s", config.Config.APIPort)
	if err := router.Run(":" + config.Config.APIPort); err != nil {
		log.Fatalf("CRITICAL: Falha ao iniciar o servidor Gin: %v", err)
	}
}
