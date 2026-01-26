package main

import (
	"context"
	"go-api-boilerplate/internal/bootstrap"
	"go-api-boilerplate/internal/config"
	"log"

	_ "go-api-boilerplate/docs"

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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app, err := bootstrap.NewApp(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer app.Close()

	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Router.Run(":8080")
}
