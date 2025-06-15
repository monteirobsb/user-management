package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler é um middleware para capturar e responder a erros de forma padronizada.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Processa a requisição

		// Após a requisição, verifica se há erros
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Tratamento específico para erros de validação
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				errorMessages := make(map[string]string)
				for _, e := range validationErrs {
					errorMessages[e.Field()] = "Erro de validação no campo " + e.Field()
				}
				c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": errorMessages})
				return
			}

			// Tratamento para outros erros
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
