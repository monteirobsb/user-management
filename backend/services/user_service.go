package services

import (
	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser cria um novo usuário no banco de dados com senha hasheada.
func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	user.Password = "" // Limpa a senha em texto plano

	result := database.DB.Create(user)
	return result.Error
}

// GetAllUsers retorna todos os usuários do banco de dados.
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	return users, result.Error
}

// GetUserByID retorna um usuário pelo seu ID.
func GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	return user, result.Error
}

// UpdateUser atualiza os dados de um usuário existente.
func UpdateUser(user *models.User, id uuid.UUID) error {
	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(user)
	return result.Error
}

// DeleteUser remove um usuário do banco de dados.
func DeleteUser(id uuid.UUID) error {
	result := database.DB.Delete(&models.User{}, "id = ?", id)
	return result.Error
}
