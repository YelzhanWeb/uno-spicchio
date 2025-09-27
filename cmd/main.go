package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/YelzhanWeb/uno-spicchio/internal/app"
	"github.com/YelzhanWeb/uno-spicchio/pkg/db"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	port := os.Getenv("APP_PORT")

	db, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("Error when opening the database: %v", err)
	}
	defer db.Close()

	addr := fmt.Sprintf(":%s", port)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := app.Routes(db)

	app.StartServerWithShutdown(addr, mux)
}
