package main

import (
	"log"
	"orders/internal/api"
	"orders/internal/config"
	"orders/internal/models"
	"orders/internal/repository"
	"orders/internal/seeds"
	"orders/internal/service"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("🚀 Starting application...")

	// Configuration
	cfg := config.Load()
	log.Println("✅ Config loaded")

	// Подключение к БД
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	log.Println("✅ Database connected")

	// Migrate schema
	err = db.AutoMigrate(
		&models.ClientType{},
		&models.User{},
		&models.Client{},
		&models.Contract{},
		&models.ContractAddress{},
		&models.Product{},
		&models.VatTax{},
		&models.IncomeTax{},
		&models.Unit{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("migration failed:", err)
	}
	log.Println("✅ Migration completed successfully")

	// Seed initial data
	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Goroutine for populating ClientType
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := seeds.SeedClientTypes(db); err != nil {
			log.Printf("Error seeding ClientTypes: %v", err)
		}
	}()

	// Goroutine for populating VatTax
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := seeds.SeedVatTaxes(db); err != nil {
			log.Printf("Error seeding VatTaxes: %v", err)
		}
	}()

	// Goroutine for populating IncomeTax
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := seeds.SeedIncomeTaxes(db); err != nil {
			log.Printf("Error seeding IncomeTaxes: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := seeds.SeedUnits(db); err != nil {
			log.Printf("Error seeding Units: %v", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("✅ Seeding completed")

	// Repository and Service
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, cfg.JWTSecret)
	log.Println("✅ Services initialized")

	// Router
	r := gin.Default()
	api.SetupRoutes(r, svc)

	log.Println("🎯 Server starting on :8080")
	r.Run(":8080")
}
