package migrations

import (
	"orders/internal/models"
)

// GetAllModels returns all models that need to be migrated
func GetAllModels() []interface{} {
	return []interface{}{
		// Reference data
		&models.ClientType{},
		&models.PriceType{},
		// Main entities
		&models.User{},
		&models.Client{},
		&models.Contract{},
		&models.ContractAddress{},
		&models.Product{},
		&models.VatTax{},
		&models.IncomeTax{},
		&models.Unit{},
		&models.PriceProduct{},
		// Documents
		&models.Order{},
		&models.OrderItem{},
	}
}
