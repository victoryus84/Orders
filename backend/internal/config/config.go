package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     	string
	DBPort     	string
	DBUser     	string
	DBPassword 	string
	DBName     	string
	DBSSLMode  	string
	JWTSecret  	string
	DSN        	string 
	Allowsignup bool
}

func Load() Config {
	godotenv.Load() // Загружаем .env файл

	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		Allowsignup: os.Getenv("ALLOWSIGNUP") == "true",
	}

	// Формируем DSN из переменных
	cfg.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	return cfg
}