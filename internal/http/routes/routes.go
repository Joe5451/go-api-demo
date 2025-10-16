package routes

import (
	"go-api-demo/internal/adapter/handlers"
	"go-api-demo/internal/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, bookHandler *handlers.BookHandler) {
	// Set up middlewares
	router.Use(middlewares.ErrorHandler())

	// Set up routes
	SetupBookRoutes(router, bookHandler)
}
