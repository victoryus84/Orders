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

// TableNameToModel maps database table names to model struct names
func TableNameToModel(tableName string) string {
	tableMap := map[string]string{
		"client_types":       "ClientType",
		"price_types":        "PriceType",
		"users":              "User",
		"clients":            "Client",
		"contracts":          "Contract",
		"contract_addresses": "ContractAddress",
		"products":           "Product",
		"vat_taxes":          "VatTax",
		"income_taxes":       "IncomeTax",
		"units":              "Unit",
		"price_products":     "PriceProduct",
		"orders":             "Order",
		"order_items":        "OrderItem",
	}

	if v, ok := tableMap[tableName]; ok {
		return v
	}
	return tableName
}
