package routes

import (
	"go-api-boilerplate/internal/adapter/handlers"
	"go-api-boilerplate/internal/http/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, bookHandler *handlers.BookHandler) {
	// Set up middlewares
	router.Use(middlewares.ErrorHandler())

	// Set up routes
	SetupBookRoutes(router, bookHandler)
}
