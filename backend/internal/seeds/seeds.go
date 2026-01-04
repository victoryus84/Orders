package seeds

import (
	"log"
	"orders/internal/models"

	"gorm.io/gorm"
)

// SeedClientTypes populates the ClientType table with initial data
func SeedClientTypes(db *gorm.DB) error {
	clientTypes := []models.ClientType{
		{Name: "individual"},
		{Name: "company"},
		{Name: "government"},
		{Name: "ngo"},
		{Name: "other"},
	}

	for _, ct := range clientTypes {
		// Check if it already exists
		var existing models.ClientType
		if err := db.Where("name = ?", ct.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			// Insert if not found
			if err := db.Create(&ct).Error; err != nil {
				log.Printf("❌ Failed to seed ClientType '%s': %v\n", ct.Name, err)
				return err
			}
			log.Printf("✅ Seeded ClientType: %s\n", ct.Name)
		} else if err != nil {
			return err
		} else {
			log.Printf("⏭️ ClientType '%s' already exists\n", ct.Name)
		}
	}

	return nil
}
