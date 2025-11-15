package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	minio "github.com/YelzhanWeb/uno-spicchio/internal/adapters/minIO"
	"github.com/YelzhanWeb/uno-spicchio/internal/adapters/postgre"
	"github.com/YelzhanWeb/uno-spicchio/internal/config"
	httpAdapter "github.com/YelzhanWeb/uno-spicchio/internal/controller/http"
	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
	"github.com/YelzhanWeb/uno-spicchio/pkg/jwt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting server with config: Host=%s, Port=%s, Env=%s", cfg.Server.Host, cfg.Server.Port, cfg.Env)

	// Connect to database
	log.Printf("Connecting to database: %s", cfg.Database.DSN())
	db, err := connectDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize MinIO storage
	storage, err := minio.NewFileStorage(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.UseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Ensure buckets exist
	ctx := context.Background()
	if err := storage.EnsureBucket(ctx, cfg.MinIO.BucketDishes); err != nil {
		log.Printf("Warning: Failed to ensure dishes bucket: %v", err)
	}
	if err := storage.EnsureBucket(ctx, cfg.MinIO.BucketUsers); err != nil {
		log.Printf("Warning: Failed to ensure users bucket: %v", err)
	}

	log.Println("MinIO storage initialized successfully")

	// Initialize repositories
	userRepo := postgre.NewUserRepository(db)
	tableRepo := postgre.NewTableRepository(db)
	categoryRepo := postgre.NewCategoryRepository(db)
	dishRepo := postgre.NewDishRepository(db)
	ingredientRepo := postgre.NewIngredientRepository(db)
	orderRepo := postgre.NewOrderRepository(db)
	supplyRepo := postgre.NewSupplyRepository(db)
	analyticsRepo := postgre.NewAnalyticsRepository(db)

	// Initialize JWT token manager
	tokenManager := jwt.NewTokenManager(cfg.JWT.Secret, cfg.JWT.ExpirationDuration())

	// Initialize services
	authService := usecase.NewAuthService(userRepo, tokenManager)
	userService := usecase.NewUserService(userRepo)
	orderService := usecase.NewOrderService(orderRepo, dishRepo, ingredientRepo, tableRepo)
	dishService := usecase.NewDishService(dishRepo)
	ingredientService := usecase.NewIngredientService(ingredientRepo)
	supplyService := usecase.NewSupplyService(supplyRepo)
	tableService := usecase.NewTableService(tableRepo)
	analyticsService := usecase.NewAnalyticsService(analyticsRepo)

	// Setup router
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

	// Serve static files from ./static directory
	staticDir := "./static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Warning: Static directory '%s' does not exist", staticDir)
	} else {
		log.Printf("Serving static files from: %s", staticDir)
		fs := http.FileServer(http.Dir(staticDir))
		r.Handle("/static/*", http.StripPrefix("/static/", fs))

		// Redirect root to login page
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/static/login.html", http.StatusFound)
		})
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
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Frontend available at: http://localhost:%s/static/login.html", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func connectDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
