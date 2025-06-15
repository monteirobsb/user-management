package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/handlers"
	"github.com/monteirobsb/user-management/backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupRouter configura um roteador Gin com um banco de dados de teste.
func setupRouter() *gin.Engine {
	// Usa SQLite em memória para testes
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("Falha ao conectar ao banco de dados de teste")
	}

	// AutoMigrate o schema E VERIFICA O ERRO
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		// Se a migração falhar, o teste vai parar aqui com uma mensagem clara.
		panic("Falha ao migrar o schema do banco de teste: " + err.Error())
	}

	database.DB = db // Sobrescreve a conexão global do DB

	r := gin.Default()
	api := r.Group("/api/users")
	{
		api.POST("", handlers.CreateUserHandler)
	}
	return r
}

func TestCreateUserHandler(t *testing.T) {
	router := setupRouter()

	newUser := `{"name":"Test User", "email":"test@example.com", "password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(newUser))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	json.Unmarshal(w.Body.Bytes(), &user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Empty(t, user.Password)        // Garante que a senha não foi retornada
	assert.NotEmpty(t, user.PasswordHash) // Garante que o hash foi criado
}
