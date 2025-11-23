package helpers

import (
	"context"
	"fmt"
	"go-api-demo/internal/bootstrap"
	"go-api-demo/internal/config"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer *postgres.PostgresContainer
	dbPool      *pgxpool.Pool
	cfg         *config.Config
	once        sync.Once
	initErr     error
)

// InitTestContainer initializes the shared postgres container
// Call this in TestMain to ensure container is ready before any tests run
func InitTestContainer() error {
	once.Do(func() {
		ctx := context.Background()

		pgContainer, initErr = postgres.Run(ctx,
			"postgres:15.3-alpine",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
			postgres.WithInitScripts(getInitSQLPath()),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second),
			),
		)
		if initErr != nil {
			return
		}

		// Get dynamic connection info from container
		host, err := pgContainer.Host(ctx)
		if err != nil {
			initErr = fmt.Errorf("failed to get container host: %w", err)
			return
		}

		port, err := pgContainer.MappedPort(ctx, "5432")
		if err != nil {
			initErr = fmt.Errorf("failed to get container port: %w", err)
			return
		}

		// Create config with actual container connection info
		cfg = &config.Config{
			Debug: false,
			Database: config.Database{
				Postgres: config.Postgres{
					Host:     host,
					Port:     port.Port(),
					User:     "testuser",
					Password: "testpass",
					DBName:   "testdb",
					Schema:   "public",
				},
			},
		}

		// Create a dedicated DB pool for cleanup operations
		connStr := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.Postgres.User,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.DBName,
		)

		dbPool, initErr = pgxpool.New(ctx, connStr)
		if initErr != nil {
			initErr = fmt.Errorf("failed to create db pool: %w", initErr)
			return
		}
	})

	return initErr
}

// TeardownTestContainer terminates the shared container
// Call this at the end of TestMain
func TeardownTestContainer() error {
	if dbPool != nil {
		dbPool.Close()
	}
	if pgContainer != nil {
		return pgContainer.Terminate(context.Background())
	}
	return nil
}

// SetupTestApp creates a new app instance using the shared container
func SetupTestApp(t *testing.T) *bootstrap.App {
	if err := InitTestContainer(); err != nil {
		t.Fatalf("failed to initialize test container: %v", err)
	}

	gin.SetMode(gin.TestMode)

	app, err := bootstrap.NewApp(context.Background(), cfg)
	if err != nil {
		t.Fatalf("failed to setup app: %v", err)
	}

	t.Cleanup(func() {
		app.Close()
	})

	return app
}

// CleanupDatabase truncates all tables and resets sequences
func CleanupDatabase(t *testing.T) {
	if dbPool == nil {
		t.Fatal("database pool not initialized")
	}

	_, err := dbPool.Exec(
		context.Background(),
		"TRUNCATE books RESTART IDENTITY CASCADE",
	)
	if err != nil {
		t.Logf("warning: failed to truncate: %v", err)
	}
}

// returns the shared database pool for direct database queries in tests
func DB() *pgxpool.Pool {
	return dbPool
}

func getInitSQLPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "init.sql")
}
