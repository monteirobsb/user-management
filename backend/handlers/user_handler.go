package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/models"
	"github.com/monteirobsb/user-management/backend/services"
	"gorm.io/gorm"
)

// CreateUserHandler lida com a criação de um novo usuário.
func CreateUserHandler(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	// A senha é passada separadamente para o serviço CreateUser.
	// A validação de senha (ex: min length) é feita via tags em UserCreateRequest.
	if err := services.CreateUser(&user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}
	// Retornar o usuário criado (PasswordHash tem json:"-")
	c.JSON(http.StatusCreated, user)
}

// GetUsersHandler lida com a listagem de todos os usuários.
func GetUsersHandler(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserHandler lida com a busca de um usuário por ID.
func GetUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUserHandler lida com a atualização de um usuário.
func UpdateUserHandler(c *gin.Context) {
	userIDParam := c.Param("id")
	id, err := uuid.Parse(userIDParam)
	if err != nil {
		log.Printf("WARN: Tentativa de atualizar usuário com ID inválido: %s, erro: %v. IP: %s", userIDParam, err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	var req models.UserUpdateRequest
	// BindJSON usará as tags de validação em UserUpdateRequest.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar o usuário existente
	userToUpdate, err := services.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		} else {
			// Para outros erros ao buscar o usuário, log já deve ter ocorrido no service
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar sua solicitação"})
		}
		return
	}

	// Aplicar atualizações se os campos foram fornecidos e validados
	if req.Name != nil {
		userToUpdate.Name = *req.Name
	}
	if req.Email != nil {
		userToUpdate.Email = *req.Email
	}

	// Chamar o serviço para atualizar o usuário.
	// A senha não é atualizada por este handler.
	if err := services.UpdateUser(&userToUpdate, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar usuário"})
		return
	}
	// Retornar o usuário atualizado (PasswordHash tem json:"-")
	c.JSON(http.StatusOK, userToUpdate)
}

// DeleteUserHandler lida com a remoção de um usuário.
func DeleteUserHandler(c *gin.Context) {
	userIDParam := c.Param("id")
	id, err := uuid.Parse(userIDParam)
	if err != nil {
		log.Printf("WARN: Tentativa de acessar usuário com ID inválido: %s, erro: %v. IP: %s", userIDParam, err, c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover usuário"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário removido com sucesso"})
}
