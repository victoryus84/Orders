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

func SeedVatTaxes(db *gorm.DB) error {
	vatTaxes := []models.VatTax{
		{Name: "VAT 20%", Rate: 20.0, Description: "20-00"},
		{Name: "VAT 10%", Rate: 10.0, Description: "10-00"},
		{Name: "VAT 6%", Rate: 6.0, Description: "6-00"},
		{Name: "VAT 5%", Rate: 5.0, Description: "5-00"},
		{Name: "VAT 0%", Rate: 0.0, Description: "0-00"},
		{Name: "VAT Exempt", Rate: 0.0, Description: "exempt"},
	}

	for _, vatTax := range vatTaxes {
		if err := db.Create(&vatTax).Error; err != nil {
			log.Printf("❌ Failed to seed VatTax '%s': %v\n", vatTax.Name, err)
			return err
		}
		log.Printf("✅ Seeded VatTax: %s\n", vatTax.Name)
	}
	return nil
}

func SeedIncomeTaxes(db *gorm.DB) error {
	incomeTaxes := []models.IncomeTax{
		{Name: "Income 12%", Rate: 12.0, Description: "12-00"},
		{Name: "Income Exempt", Rate: 0.0, Description: "exempt"},
	}

	for _, incomeTax := range incomeTaxes {
		if err := db.Create(&incomeTax).Error; err != nil {
			log.Printf("❌ Failed to seed IncomeTax '%s': %v\n", incomeTax.Name, err)
			return err
		}
		log.Printf("✅ Seeded IncomeTax: %s\n", incomeTax.Name)
	}
	return nil
}

func SeedUnits(db *gorm.DB) error {
	units := []models.Unit{
		{Name: "buc", Description: "bucăți"},
		{Name: "kg", Description: "kilograme"},
		{Name: "m", Description: "metri"},
		{Name: "l", Description: "litri"},
		{Name: "set", Description: "seturi"},
		{Name: "pet", Description: "plastic"},
	}

	for _, unit := range units {
		if err := db.Create(&unit).Error; err != nil {
			log.Printf("❌ Failed to seed Unit '%s': %v\n", unit.Name, err)
			return err
		}
		log.Printf("✅ Seeded Unit: %s\n", unit.Name)
	}
	return nil
}
