package bootstrap

import (
	"context"
	"go-api-demo/internal/adapter/handlers"
	"go-api-demo/internal/adapter/repositories"
	"go-api-demo/internal/application"
	"go-api-demo/internal/config"
	"go-api-demo/internal/http/routes"
	"go-api-demo/internal/infra"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router *gin.Engine
	db     *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := infra.NewPostgresPool(ctx, cfg.Database.Postgres, cfg.Debug)
	if err != nil {
		return nil, err
	}

	// Dependency Injection
	bookRepo := repositories.NewPostgresBookRepo(db)
	bookService := application.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	// Setup Router
	router := gin.New()
	router.Use(gin.Recovery())
	if cfg.Debug {
		router.Use(gin.Logger())
	}
	routes.SetupRoutes(router, bookHandler)

	return &App{Router: router, db: db}, nil
}

func (a *App) Close() {
	a.db.Close()
}
