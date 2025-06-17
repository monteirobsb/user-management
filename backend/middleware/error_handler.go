package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler é um middleware para capturar e responder a erros de forma padronizada.
// Ele também loga os erros não tratados que chegam até ele.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Processa a requisição

		// Após a requisição, verifica se há erros registrados no contexto do Gin.
		if len(c.Errors) > 0 {
			// Pega o último erro. Gin permite múltiplos erros, mas geralmente tratamos o mais recente.
			err := c.Errors.Last().Err

			// Log genérico para qualquer erro que chegue aqui.
			// Este log é útil para capturar erros inesperados ou não tratados especificamente em outros lugares.
			// c.Errors.String() pode fornecer uma representação textual de todos os erros no contexto.
			log.Printf("ERROR: Erro durante o processamento da requisição %s %s: %v. Detalhes Gin: %s. IP: %s",
				c.Request.Method,
				c.Request.URL.Path,
				err,
				c.Errors.String(),
				c.ClientIP(),
			)

			// Tratamento específico para erros de validação do go-playground/validator.
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				errorMessages := make(map[string]string)
				for _, e := range validationErrs {
					// Personalize a mensagem de erro como preferir.
					// e.Tag() dá o tipo de validação que falhou (ex: "required", "email", "min").
					// e.Param() pode dar parâmetros da tag (ex: "8" para "min=8").
					errorMessages[e.Field()] = "Erro de validação no campo '" + e.Field() + "': falha na regra '" + e.Tag() + "'"
				}
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos", "details": errorMessages})
				return
			}

			// Tratamento para outros tipos de erros.
			// Se o erro já tiver sido tratado e um status JSON enviado (com AbortWithStatusJSON),
			// esta parte pode não ser executada ou pode tentar escrever no response novamente,
			// o que o Gin geralmente impede.
			// Se o response ainda não foi escrito, podemos definir um erro genérico.
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ocorreu um erro interno no servidor."})
			}
		}
	}
}
