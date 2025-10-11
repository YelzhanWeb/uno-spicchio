// file: cmd/app/main.go

package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/YelzhanWeb/uno-spicchio/internal/adapters/postgre"
	"github.com/YelzhanWeb/uno-spicchio/internal/app"
	"github.com/YelzhanWeb/uno-spicchio/internal/config"
	"github.com/YelzhanWeb/uno-spicchio/internal/controller/http"
	"github.com/YelzhanWeb/uno-spicchio/internal/usecase"
	"github.com/YelzhanWeb/uno-spicchio/pkg/db"
	"github.com/jmoiron/sqlx" // <-- Добавляем импорт sqlx
)

func main() {
	// 1. Инициализация конфигурации
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Configuration initialization error: %v", err)
	}

	// 2. Инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	slog.Info("Logger initialized")

	// 3. Инициализация подключения к БД
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgre.Host,
		cfg.Postgre.Port,
		cfg.Postgre.UserName,
		cfg.Postgre.Password,
		cfg.Postgre.DBName,
	)

	// Предполагаем, что db.InitDB() возвращает *sql.DB.
	// Оборачиваем его в *sqlx.DB для работы с нашими репозиториями.
	sqlDB, err := db.InitDB(dsn)
	if err != nil {
		slog.Error("Error when opening the database", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	// Используем sqlx
	db := sqlx.NewDb(sqlDB, "postgres")
	slog.Info("Database connection pool established")

	// --- СБОРКА ВСЕХ ЗАВИСИМОСТЕЙ (DEPENDENCY INJECTION) ---

	// 4. Инициализация Репозиториев (Слой данных)
	// Вам нужно будет создать NewUserRepository и NewCategoryRepository
	// по аналогии с тем, как мы сделали NewOrderRepository и NewDishRepository
	userRepo := postgre.NewUserRepository(db)         // <-- Вам нужно создать этот конструктор
	categoryRepo := postgre.NewCategoryRepository(db) // <-- и этот
	dishRepo := postgre.NewDishRepository(db)
	orderRepo := postgre.NewOrderRepository(db)
	slog.Info("Repositories initialized")

	// 5. Инициализация Сервисов (Слой бизнес-логики)
	deps := usecase.Dependencies{
		UserRepo:     userRepo,
		CategoryRepo: categoryRepo,
		OrderRepo:    orderRepo,
		DishRepo:     dishRepo,
	}
	useCases := usecase.NewUseCases(deps)
	slog.Info("Use cases initialized")

	// 6. Инициализация Обработчиков (Слой контроллеров)
	handler := http.NewHandler(useCases, categoryRepo)
	slog.Info("Handler initialized")

	// 7. Инициализация Роутера
	router := http.NewRouter(handler)
	slog.Info("Router initialized")

	// 8. Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.HTTP.Port) // Используем порт из конфига
	slog.Info("Starting server", "address", addr)
	app.StartServerWithShutdown(addr, router)
}
