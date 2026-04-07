package repository

import (
	"orders/internal/models"
	"strings"
	"gorm.io/gorm"
)

// Repository - structura principală pentru acces la baza de date
type Repository struct {
	db *gorm.DB
}

// Creează o nouă instanță de Repository cu conexiunea la DB
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Creează un nou utilizator în baza de date
func (repository *Repository) CreateUser(user *models.User) error {
	return repository.db.Create(user).Error
}

// Găsește un utilizator după email
func (repository *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repository.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Client methods
func (repository *Repository) CreateClient(client *models.Client) error {
    if client.Email != nil {
        em := strings.TrimSpace(*client.Email)
        el := strings.ToLower(em)

        // Dacă e mizerie (n/a, none, gol), îl facem NIL
        if em == "" || el == "not inserted" || el == "n/a" || el == "none" {
            client.Email = nil // În DB se va duce NULL (și e unic!)
        }
    }

    return repository.db.Create(client).Error
}
func (repository *Repository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := repository.db.Preload("ClientType").Limit(1000).Find(&clients).Error
	return clients, err
}
func (repository *Repository) FindClientsByQuery(query string) ([]models.Client, error) {
	if len(query) < 3 {
		return []models.Client{}, nil // Return empty if less than 3 chars
	}
	var clients []models.Client
	err := repository.db.
		Where("name ILIKE ? OR email ILIKE ? OR fiscal_id ILIKE ? OR phone ILIKE ?",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%").
    Limit(50).
    Find(&clients).Error
	return clients, err
}
func (repository *Repository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := repository.db.Preload("ClientType").First(&client, id).Error
	return &client, err
}
func (repository *Repository) FindClientByFiscalID(fiscalID string) (*models.Client, error) {
	var client models.Client
	err := repository.db.Where("fiscal_id = ?", fiscalID).First(&client).Error
	return &client, err
}
func (repository *Repository) CreateClientAddress(addr *models.ClientAddress) error {
	return repository.db.Create(addr).Error
}

// Contract methods
func (repository *Repository) CreateContract(contract *models.Contract) error {
	return repository.db.Create(contract).Error
}
func (repository *Repository) FindContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	err := repository.db.First(&contract, id).Error
	return &contract, err
}
func (repository *Repository) CreateContractAddress(addr *models.ContractAddress) error {
	return repository.db.Create(addr).Error
}
func (repository *Repository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := repository.db.First(&addr, id).Error
	return &addr, err
}

func (repository *Repository) FindContractByClientID(clientID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	err := repository.db.Where("client_id = ?", clientID).Find(&contracts).Error
	return contracts, err
}

// Product methods
func (repository *Repository) CreateProduct(product *models.Product) error {
	return repository.db.Create(product).Error
}
func (repository *Repository) FindProductGroupByID(id uint) (*models.ProductGroup, error) {
	var group models.ProductGroup
	err := repository.db.First(&group, id).Error
	return &group, err
}
func (repository *Repository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := repository.db.First(&product, id).Error
	return &product, err
}
func (repository *Repository) FindVatTaxByID(id uint) (*models.VatTax, error) {
	var vatTax models.VatTax
	err := repository.db.First(&vatTax, id).Error
	return &vatTax, err
}
func (repository *Repository) FindUnitByID(id uint) (*models.Unit, error) {
	var unit models.Unit
	err := repository.db.First(&unit, id).Error
	return &unit, err
}

// Order methods
func (repository *Repository) CreateOrder(order *models.Order) error {
	return repository.db.Create(order).Error
}
func (repository *Repository) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := repository.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
func (repository *Repository) FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := repository.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}
