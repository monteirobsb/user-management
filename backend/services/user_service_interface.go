package services

import (
	"github.com/google/uuid"
	"github.com/monteirobsb/user-management/backend/models"
)

// UserServiceInterface define as operações do serviço de usuário.
// Esta interface será usada para mocking nos testes de handler.
type UserServiceInterface interface {
	CreateUser(user *models.User, plainPassword string) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uuid.UUID) (models.User, error)
	UpdateUser(user *models.User, id uuid.UUID) error
	DeleteUser(id uuid.UUID) error
}
