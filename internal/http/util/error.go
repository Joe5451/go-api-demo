package util

import (
	"go-api-demo/internal/constant"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code    constant.ErrorCode `json:"code" example:"ERROR_CODE"`
	Message string             `json:"message" example:"Error message"`
}

func NewError(c *gin.Context, status int, code constant.ErrorCode, err error) {
	c.JSON(status, HTTPError{
		Code:    code,
		Message: err.Error(),
	})
}
