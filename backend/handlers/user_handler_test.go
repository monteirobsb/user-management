package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/handlers"
	"github.com/monteirobsb/user-management/backend/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
		api.PUT("/:id", handlers.UpdateUserHandler)
		api.GET("/:id", handlers.GetUserHandler)
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
	assert.Empty(t, user.Password) // Garante que a senha não foi retornada
	// PasswordHash não é mais retornado na resposta JSON
}

func TestUpdateUserPasswordHandler(t *testing.T) {
	router := setupRouter()

	// 1. Criar um novo usuário
	initialPassword := "password123"
	createUserPayload := map[string]string{
		"name":     "Update Test User",
		"email":    "update@example.com",
		"password": initialPassword,
	}
	jsonPayload, _ := json.Marshal(createUserPayload)
	reqCreate, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonPayload))
	reqCreate.Header.Set("Content-Type", "application/json")

	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, wCreate.Code)

	var createdUser models.User
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)

	userID := createdUser.ID

	// 2. Atualizar a senha do usuário
	newPassword := "newPassword456"
	updatePasswordPayload := map[string]string{"password": newPassword}
	jsonUpdatePayload, _ := json.Marshal(updatePasswordPayload)

	reqUpdate, _ := http.NewRequest("PUT", "/api/users/"+userID.String(), bytes.NewBuffer(jsonUpdatePayload))
	reqUpdate.Header.Set("Content-Type", "application/json")

	wUpdate := httptest.NewRecorder()
	router.ServeHTTP(wUpdate, reqUpdate)
	assert.Equal(t, http.StatusOK, wUpdate.Code)

	// 3. Verificar a atualização da senha no banco de dados
	var updatedUser models.User
	err = database.DB.First(&updatedUser, "id = ?", userID).Error
	assert.NoError(t, err)

	// Verificar se a senha antiga não corresponde mais
	errOldPassword := bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte(initialPassword))
	assert.Error(t, errOldPassword, "A senha antiga ainda corresponde, o que não deveria acontecer")

	// Verificar se a nova senha corresponde
	errNewPassword := bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte(newPassword))
	assert.NoError(t, errNewPassword, "A nova senha não corresponde")
}
