package main

import (
	"log"
	"os"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/router"
	"github.com/franklindh/simedis-api/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("could not load config: %v", err)
	}

	utils.RegisterSanitizeValidator()

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {

		logger.Fatalf("could not open gorm connection: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("could not get underline sql.DB: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
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
