package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/router"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("could not load config: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {

		logger.Fatalf("could not open sql connection: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatalf("could not ping database: %v", err)
	}

	logger.Println("Database connection pool established")

	app := &config.Application{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}

	r := router.New(app)

	logger.Printf("Starting server on port %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		logger.Fatalf("could not start server: %v", err)
	}
}
