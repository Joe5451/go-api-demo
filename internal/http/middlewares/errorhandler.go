package middlewares

import (
	"errors"
	"go-api-demo/internal/constant"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()
		if lastErr == nil {
			return
		}

		if errors.Is(lastErr.Err, constant.ErrValidation) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    constant.ErrValidationCode,
				"message": lastErr.Err.Error(),
			})
			return
		}

		log.Printf("[INTERNAL_ERROR]: %v\n", lastErr.Err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred.",
		})
	}
}
