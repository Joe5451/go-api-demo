package main

import (
	"context"
	"go-api-demo/internal/adapter/handlers"
	"go-api-demo/internal/adapter/repositories"
	"go-api-demo/internal/application"
	"go-api-demo/internal/config"
	"go-api-demo/internal/http/routes"
	"go-api-demo/internal/infra"
	"log"

	"github.com/gin-gonic/gin"

	_ "go-api-demo/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API Demo
// @version         1.0
// @description     This is a sample server API Demo.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}
