package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	minio "github.com/YelzhanWeb/uno-spicchio/internal/adapters/minIO"
	"github.com/YelzhanWeb/uno-spicchio/internal/adapters/postgre"
	"github.com/YelzhanWeb/uno-spicchio/internal/config"
	httpAdapter "github.com/YelzhanWeb/uno-spicchio/internal/controller/http"
	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
	"github.com/YelzhanWeb/uno-spicchio/pkg/jwt"
	"github.com/YelzhanWeb/uno-spicchio/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// добавь ЭТО:
func init() {
	// убираем стандартную дату/время, чтобы не было дубля
	log.SetFlags(0)
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	logger.Startup("===========================================")
	logger.Startup("🍕 Starting UNO Spicchio Server")
	logger.Startup("===========================================")
	logger.Info("Host=%s, Port=%s, Env=%s", cfg.Server.Host, cfg.Server.Port, cfg.Env)

	// Connect to database
	logger.Info("Connecting to database: %s", cfg.Database.DSN())
	db, err := connectDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Success("✓ Connected to database successfully")

	// Initialize MinIO storage
	logger.Info("Initializing MinIO storage...")
	storage, err := minio.NewFileStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		logger.Fatal("Failed to initialize MinIO: %v", err)
	}

	// Ensure buckets exist
	ctx := context.Background()
	if err := storage.EnsureBucket(ctx, cfg.MinIO.BucketDishes); err != nil {
		logger.Warning("Failed to ensure dishes bucket: %v", err)
	}
	if err := storage.EnsureBucket(ctx, cfg.MinIO.BucketUsers); err != nil {
		logger.Warning("Failed to ensure users bucket: %v", err)
	}

	logger.Success("✓ MinIO storage initialized successfully")

	// Initialize repositories
	logger.Info("Initializing repositories...")
	userRepo := postgre.NewUserRepository(db)
	tableRepo := postgre.NewTableRepository(db)
	categoryRepo := postgre.NewCategoryRepository(db)
	dishRepo := postgre.NewDishRepository(db)
	ingredientRepo := postgre.NewIngredientRepository(db)
	orderRepo := postgre.NewOrderRepository(db)
	supplyRepo := postgre.NewSupplyRepository(db)
	analyticsRepo := postgre.NewAnalyticsRepository(db)
	logger.Success("✓ Repositories initialized")

	// Initialize JWT token manager
	logger.Info("Initializing JWT token manager...")
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, cfg.JWT.ExpirationDuration())
	logger.Success("✓ JWT token manager initialized")

	// Initialize services
	logger.Info("Initializing services...")
	authService := usecase.NewAuthService(userRepo, tokenManager)
	userService := usecase.NewUserService(userRepo)
	orderService := usecase.NewOrderService(orderRepo, dishRepo, ingredientRepo, tableRepo)
	dishService := usecase.NewDishService(dishRepo)
	ingredientService := usecase.NewIngredientService(ingredientRepo)
	supplyService := usecase.NewSupplyService(supplyRepo)
	tableService := usecase.NewTableService(tableRepo)
	analyticsService := usecase.NewAnalyticsService(analyticsRepo)
	logger.Success("✓ Services initialized")

	// Setup router
	logger.Info("Setting up router...")
	router := httpAdapter.NewRouter(
		authService,
		userService,
		orderService,
		dishService,
		ingredientService,
		supplyService,
		tableService,
		categoryRepo,
		analyticsService,
		tokenManager,
	)

	// Get base router
	r := router.Setup()
	logger.Success("✓ Router configured")

	// Serve static files from ./static directory
	staticDir := "./static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		logger.Warning("Static directory '%s' does not exist", staticDir)
	} else {
		logger.Info("Serving static files from: %s", staticDir)
		fs := http.FileServer(http.Dir(staticDir))
		r.Handle("/static/*", http.StripPrefix("/static/", fs))

		// Redirect root to login page
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/static/login.html", http.StatusFound)
		})
		logger.Success("✓ Static files configured")
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Startup("===========================================")
		logger.Startup("🚀 Server is running!")
		logger.Startup("===========================================")
		logger.Info("Address: http://localhost:%s", cfg.Server.Port)
		logger.Info("Frontend: http://localhost:%s/static/login.html", cfg.Server.Port)
		logger.Info("Environment: %s", cfg.Env)
		logger.Startup("===========================================")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warning("⚠ Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Success("✓ Server exited gracefully")
}

func connectDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	logger.Database("Opening database connection...")
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	logger.Database("Connection pool configured (Max: 25, Idle: 5)")

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Database("Pinging database...")
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
