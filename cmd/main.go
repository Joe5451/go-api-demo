package main

import (
	"context"
	"go-api-demo/internal/adapter/handlers"
	"go-api-demo/internal/adapter/repositories"
	"go-api-demo/internal/application"
	"go-api-demo/internal/config"
	"go-api-demo/internal/infra"
	"go-api-demo/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := infra.NewPostgresPool(context.Background(), config.Database.Postgres, config.Debug)
	if err != nil {
		log.Fatalf("failed to create postgres pool: %v", err)
	}
	defer db.Close()

	bookRepo := repositories.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	router := gin.New()
	router.Use(gin.Recovery())
	if config.Debug {
		router.Use(gin.Logger())
	}
	routes.SetupRoutes(router, bookHandler)
	router.Run(":8080")
}
