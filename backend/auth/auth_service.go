package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// tokenDuration define o tempo de expiração do token.
// Idealmente, este valor viria de uma configuração (ex: variável de ambiente).
const tokenDuration = 24 * time.Hour
const errorInvalidCredentials = "usuário não encontrado ou credenciais inválidas"

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// LoginUser verifica as credenciais e retorna um token se forem válidas.
func LoginUser(email, password string) (string, error) {
	if len(jwtKey) == 0 {
		// Em um cenário real, logar este erro criticamente e não expor detalhes.
		// Ex: log.Error("CRÍTICO: JWT_SECRET_KEY não está configurada.")
		return "", errors.New("falha na autenticação devido a erro de configuração interna")
	}

	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New(errorInvalidCredentials)
		}
		// Logar o erro real para diagnóstico interno: fmt.Errorf("erro ao buscar usuário %s: %w", email, result.Error)
		return "", errors.New("erro ao processar login") // Mensagem genérica para outros erros de DB
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// O erro aqui pode ser bcrypt.ErrMismatchedHashAndPassword ou outro.
		return "", errors.New(errorInvalidCredentials) // Mesma mensagem para evitar enumeração de usuários
	}

	expirationTime := time.Now().Add(tokenDuration)
	claims := &Claims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// Logar o erro real para diagnóstico interno: fmt.Errorf("erro ao assinar token para usuário %s: %w", user.ID, err)
		return "", errors.New("erro ao gerar token de autenticação")
	}

	return tokenString, nil
}
