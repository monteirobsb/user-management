package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/monteirobsb/user-management/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase inicializa a conexão com o banco de dados e executa as migrações.
func InitDatabase() error {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("DATABASE_PORT")
	sslmode := os.Getenv("DATABASE_SSLMODE") // Permitir configuração do sslmode

	if sslmode == "" {
		sslmode = "disable" // Padrão para desenvolvimento; considere "require" ou "verify-full" para produção
	}

	// Verifica se as variáveis essenciais não estão vazias
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		return errors.New("variáveis de ambiente do banco de dados não estão completamente definidas")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso.")

	// AutoMigrate
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("falha ao migrar o schema do banco de dados: %w", err)
	}
	log.Println("Schema do banco de dados migrado com sucesso.")
	return nil
}
