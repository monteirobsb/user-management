package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User representa a estrutura de um usuário no banco de dados
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name         string    `gorm:"size:255;not null" json:"name" binding:"required"`
	Email        string    `gorm:"size:255;not null;unique" json:"email" binding:"required,email"`
	Password     string    `gorm:"-" json:"password,omitempty" binding:"omitempty,min=8"` // omitempty para edição, min=8 para criação
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at"`
}

// BeforeCreate é um hook do GORM que será chamado antes de um usuário ser criado.
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Gera um novo UUID e o atribui ao ID do usuário
	user.ID = uuid.New()
	return
}
