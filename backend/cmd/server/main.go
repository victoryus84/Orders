package main

import (
	"log"
	"orders/internal/api"
	"orders/internal/config"
	"orders/internal/models"
	"orders/internal/repository"
	"orders/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Configuration
	cfg := config.Load()

	//log.Println("DSN:", cfg.DSN)

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Migrate
	// Migrate schema
    err = db.AutoMigrate(
        &models.User{},
        &models.Client{},
        &models.Contract{},
        &models.ContractAddress{},
        &models.Product{},
        &models.Order{},
        &models.OrderItem{},
    )
    if err != nil {
        log.Fatal("migration failed:", err)
    }

    log.Println("Migration completed successfully")

	// Repository and Service
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg.JWTSecret)

	// Router
	r := gin.Default()
	api.SetupRoutes(r, svc)

	r.Run(":8080")
}