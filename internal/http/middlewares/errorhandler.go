package middlewares

import (
	"errors"
	"go-api-demo/internal/constant"
	"go-api-demo/internal/http/util"
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
			util.NewError(
				c,
				http.StatusBadRequest,
				constant.ErrValidationCode,
				lastErr.Err,
			)
			return
		}

		log.Printf("[INTERNAL_ERROR]: %v\n", lastErr.Err)
		util.NewError(
			c,
			http.StatusInternalServerError,
			constant.ErrInternalServerError,
			errors.New("an unexpected error occurred"),
		)
	}
}
