package models

// UserCreateRequest defines the structure for creating a new user.
// It includes the plaintext password which will be hashed by the service.
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserUpdateRequest defines the structure for updating an existing user.
// Fields are optional. If a field is provided, it will be validated.
// For example, if Name is provided, it must not be empty.
// If Email is provided, it must be a valid email format.
type UserUpdateRequest struct {
	Name  *string `json:"name,omitempty" binding:"omitempty,min=1"` // If Name is provided, it must not be empty
	Email *string `json:"email,omitempty" binding:"omitempty,email"` // If Email is provided, it must be a valid email
}
