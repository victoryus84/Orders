package migrations

import (
	"orders/internal/models"
)

// GetAllModels returns all models that need to be migrated
func GetAllModels() []interface{} {
	return []interface{}{
		&models.ClientType{},
		&models.User{},
		&models.Client{},
		&models.Contract{},
		&models.ContractAddress{},
		&models.Product{},
		&models.VatTax{},
		&models.IncomeTax{},
		&models.Unit{},
		&models.PriceType{},
		&models.PriceProduct{},
		&models.Order{},
		&models.OrderItem{},
	}
}
