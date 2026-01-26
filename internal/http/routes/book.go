package routes

import (
	"go-api-boilerplate/internal/adapter/handlers"

	"github.com/gin-gonic/gin"
)

func SetupBookRoutes(router *gin.Engine, bookHandler *handlers.BookHandler) {
	router.POST("/books", bookHandler.CreateBook)
	router.GET("/books/:id", bookHandler.GetBook)
	router.GET("/books", bookHandler.GetBooks)
	router.PUT("/books/:id", bookHandler.UpdateBook)
	router.DELETE("/books/:id", bookHandler.DeleteBook)
}
