package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/monteirobsb/user-management/backend/auth"
)

// AuthMiddleware verifica o token JWT na requisição.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("WARN: Tentativa de acesso não autorizado à rota %s (IP: %s): Cabeçalho de autorização não encontrado.", c.FullPath(), c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho de autorização não encontrado"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Significa que o prefixo "Bearer " não estava lá
			log.Printf("WARN: Tentativa de acesso não autorizado à rota %s (IP: %s): Formato de token inválido (sem prefixo 'Bearer ').", c.FullPath(), c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}

		// JWT_SECRET_KEY é verificado no init() do pacote auth. Se estiver vazio, o app não inicia.
		// Portanto, os.Getenv("JWT_SECRET_KEY") aqui deve ser seguro.
		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		if len(jwtKey) == 0 {
			// Esta situação não deveria ocorrer se o init() do pacote auth funcionou.
			// Mas, por segurança, logamos um erro crítico se, por algum motivo, jwtKey for vazio aqui.
			log.Printf("CRITICAL: JWT_SECRET_KEY resultou em chave vazia no AuthMiddleware. Rota: %s, IP: %s", c.FullPath(), c.ClientIP())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro de configuração interna do servidor"})
			return
		}


		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validação do método de assinatura é importante aqui se múltiplos são possíveis.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("WARN: Método de assinatura inesperado no token: %v. Rota: %s, IP: %s", token.Header["alg"], c.FullPath(), c.ClientIP())
				return nil, jwt.ErrSignatureInvalid // Ou um erro mais específico
			}
			return jwtKey, nil
		})

		if err != nil {
			log.Printf("WARN: Tentativa de acesso não autorizado à rota %s (IP: %s): Token inválido ou expirado. Erro: %v", c.FullPath(), c.ClientIP(), err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		if !token.Valid {
			// Este caso pode ser redundante se err != nil já o cobre, mas é uma checagem explícita.
			log.Printf("WARN: Tentativa de acesso não autorizado à rota %s (IP: %s): Token marcado como inválido (sem erro explícito na parse).", c.FullPath(), c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
