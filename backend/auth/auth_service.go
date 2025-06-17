package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/monteirobsb/user-management/backend/database"
	"github.com/monteirobsb/user-management/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey []byte

// init é chamada automaticamente quando o pacote é inicializado.
// Verifica a configuração da JWT_SECRET_KEY.
func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("CRÍTICO: JWT_SECRET_KEY não está configurada. A aplicação não pode iniciar sem esta chave.")
	}
	jwtKey = []byte(secret)
}

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
	// A verificação de jwtKey vazia foi movida para a função init().
	// Se a chave não estiver configurada, a aplicação já terá sido encerrada.

	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Não logar "record not found" como erro, pois é um caso de login esperado (usuário não existe).
			return "", errors.New(errorInvalidCredentials)
		}
		log.Printf("ERROR: Falha ao buscar usuário com email %s: %v", email, result.Error)
		return "", errors.New("erro ao processar login") // Mensagem genérica para outros erros de DB
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// O erro aqui pode ser bcrypt.ErrMismatchedHashAndPassword (que não é um erro do sistema, mas falha de autenticação)
		// ou outro erro (problema com o hash armazenado, etc., que seria um erro do sistema).
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			// Não logar como erro do sistema, é uma falha de login esperada.
			// Poderia ter um log de DEBUG/INFO aqui se houvesse níveis e necessidade de rastrear tentativas.
		} else {
			// Logar outros erros de bcrypt como erro do sistema.
			log.Printf("ERROR: Falha ao comparar hash para usuário com email %s: %v", email, err)
		}
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
		log.Printf("ERROR: Falha ao assinar token para usuário ID %s: %v", user.ID.String(), err)
		return "", errors.New("erro ao gerar token de autenticação")
	}

	return tokenString, nil
}
