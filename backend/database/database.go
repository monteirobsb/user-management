package database

import (
	"fmt"
	"log"

	"github.com/monteirobsb/user-management/backend/config" // Import the new config package
	"github.com/monteirobsb/user-management/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase inicializa a conexão com o banco de dados e executa as migrações.
// Em caso de falha crítica na configuração ou conexão, esta função chamará log.Fatal.
func InitDatabase() { // Alterado para não retornar erro, pois usará log.Fatal
	// Configuration is now loaded from config.Config
	// Checks for empty essential variables are done in config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Config.DatabaseHost,
		config.Config.PostgresUser,
		config.Config.PostgresPassword,
		config.Config.PostgresDB,
		config.Config.DatabasePort,
		config.Config.DatabaseSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("CRITICAL: Falha ao conectar ao banco de dados (%s). Verifique as credenciais e a disponibilidade do BD. Erro: %v. A aplicação não pode iniciar.", dsn, err)
	}

	log.Print("INFO: Conexão com o banco de dados estabelecida com sucesso.")

	log.Print("INFO: Iniciando migração do schema do banco de dados...")
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("CRITICAL: Falha ao migrar o schema do banco de dados: %v. A aplicação não pode iniciar.", err)
	}
	log.Print("INFO: Schema do banco de dados migrado com sucesso.")
}
