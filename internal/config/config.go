package config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DSN                    string
	DefaultPetugasPassword string
}

type Application struct {
	Config *Config
	DB     *sql.DB
	Logger *log.Logger
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found", err)
	}

	if os.Getenv("DB_USER") == "" || os.Getenv("DB_NAME") == "" {
		return nil, errors.New("DB_USER and DB_NAME must be set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))

	return &Config{Port: os.Getenv("API_PORT"), DSN: dsn, DefaultPetugasPassword: os.Getenv("DEFAULT_PETUGAS_PASSWORD")}, nil
}
