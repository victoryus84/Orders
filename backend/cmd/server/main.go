package main

import (
	"log"
	"orders/internal/api"
	"orders/internal/config"
	"orders/internal/migrations"
	"orders/internal/repository"
	"orders/internal/seeds"
	"orders/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("🚀 Starting application...")

	// Configuration
	cfg := config.Load()
	log.Println("✅ Config loaded")

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	log.Println("✅ Database connected")

	// Migrate schema
	if err := db.AutoMigrate(migrations.GetAllModels()...); err != nil {
		log.Fatal("migration failed:", err)
	}
	log.Println("✅ Migration completed successfully")

	// Analyze schema differences
	migrations.AnalyzeSchemaSync(db)
	migrations.PrintSyncCommands(db)

	// Clean up orphaned columns
	if err := migrations.DropUnusedColumns(db); err != nil {
		log.Fatal("cleanup failed:", err)
	}
	log.Println("✅ Cleanup completed successfully")

	// Seed initial data
	seeds.RunAllSeeds(db)
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
