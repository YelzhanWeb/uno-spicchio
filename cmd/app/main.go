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
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Configuration initialization error")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgre.Host,
		cfg.Postgre.Port,
		cfg.Postgre.UserName,
		cfg.Postgre.Password,
		cfg.Postgre.DBName,
	)

	db, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error when opening the database: %v", err)
	}
	defer db.Close()

	addr := fmt.Sprintf(":%s", 8080)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	postgrePool := postgre.NewPoolDB(db)
	service := usecase.NewService(postgrePool)
	mux := http.Router(service)

	app.StartServerWithShutdown(addr, mux)
}

// TODO: init config: cleanenv
// TODO: init logger: slog
// TODO: init storage: postgresql
// TODO: init router: chi
// TODO: run server
