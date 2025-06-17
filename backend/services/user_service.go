package services

import (
	"log"

	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser cria um novo usuário no banco de dados com senha hasheada.
// Aceita o usuário a ser criado e a senha em texto plano.
func CreateUser(user *models.User, plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: Falha ao gerar hash de senha para novo usuário (email: %s): %v", user.Email, err)
		return err
	}
	user.PasswordHash = string(hashedPassword)

	result := database.DB.Create(user)
	if result.Error != nil {
		log.Printf("ERROR: Falha ao criar usuário (email: %s) no banco de dados: %v", user.Email, result.Error)
		return result.Error
	}
	return nil
}

// GetAllUsers retorna todos os usuários do banco de dados.
func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		log.Printf("ERROR: Falha ao buscar todos os usuários: %v", result.Error)
		return nil, result.Error
	}
	return users, nil
}

// GetUserByID retorna um usuário pelo seu ID.
func GetUserByID(id uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		// Nota: gorm.ErrRecordNotFound será logado aqui como ERROR.
		// O handler é responsável por traduzir isso para uma resposta 404 se apropriado.
		log.Printf("ERROR: Falha ao buscar usuário com ID %s: %v", id, result.Error)
		return user, result.Error
	}
	return user, nil
}

// UpdateUser atualiza os dados de um usuário existente.
// O ID é usado para identificar o usuário, e o user *models.User contém os campos a serem atualizados.
func UpdateUser(user *models.User, id uuid.UUID) error {
	// A lógica de hashing de senha em UpdateUser é mantida conforme original,
	// mas o UpdateUserHandler agora não preenche user.Password.
	// Esta lógica permaneceria para outros usos potenciais ou refatorações futuras.
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("ERROR: Falha ao gerar hash de nova senha para usuário ID %s: %v", id, err)
			return err
		}
		user.PasswordHash = string(hashedPassword)
		user.Password = "" // Limpa a senha em texto plano
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		log.Printf("ERROR: Falha ao atualizar usuário ID %s no banco de dados: %v", id, result.Error)
		return result.Error
	}
	return nil
}

// DeleteUser remove um usuário do banco de dados.
func DeleteUser(id uuid.UUID) error {
	result := database.DB.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("ERROR: Falha ao deletar usuário ID %s do banco de dados: %v", id, result.Error)
		return result.Error
	}
	return nil
}
