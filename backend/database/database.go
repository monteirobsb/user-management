package database

import (
	"fmt"
	"log"
	"os"

	"github.com/monteirobsb/user-management/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase inicializa a conexão com o banco de dados e executa as migrações.
// Em caso de falha crítica na configuração ou conexão, esta função chamará log.Fatal.
func InitDatabase() { // Alterado para não retornar erro, pois usará log.Fatal
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")
	sslmode := os.Getenv("DATABASE_SSLMODE")

	if sslmode == "" {
		sslmode = "disable" // Padrão para desenvolvimento
		log.Print("INFO: DATABASE_SSLMODE não definido, usando 'disable' como padrão.")
	}

	// Verifica se as variáveis essenciais não estão vazias
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("CRITICAL: Variáveis de ambiente do banco de dados não estão completamente definidas. A aplicação não pode iniciar.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
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
