package repository

import (
	"orders/internal/models"

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

// Creează o nouă comandă (Order) în baza de date
func (repository *Repository) CreateOrder(order *models.Order) error {
	return repository.db.Create(order).Error
}

// Găsește toate comenzile pentru un anumit utilizator (după user_id/owner_id)
func (repository *Repository) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := repository.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}

// Găsește o comandă după ID (cu OrderItems preload)
func (repository *Repository) FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := repository.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}

// Creează un nou client în baza de date
func (repository *Repository) CreateClient(client *models.Client) error {
	return repository.db.Create(client).Error
}

// Returnează primii 1000 de clienți din baza de date
func (repository *Repository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := repository.db.Limit(1000).Find(&clients).Error
	return clients, err
}

// Găsește până la 5 clienți după un substring (minim 5 caractere)
func (repository *Repository) FindClientsByQuery(query string) ([]models.Client, error) {
	if len(query) < 3 {
		return []models.Client{}, nil // Return empty if less than 5 chars
	}
	var clients []models.Client
	err := repository.db.
		Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(50).
		Find(&clients).Error
	return clients, err
}

// Găsește un client după ID
func (repository *Repository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := repository.db.First(&client, id).Error
	return &client, err
}

// Creează un nou contract în baza de date
func (repository *Repository) CreateContract(contract *models.Contract) error {
	return repository.db.Create(contract).Error
}

// Găsește un contract după ID
func (repository *Repository) FindContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	err := repository.db.First(&contract, id).Error
	return &contract, err
}

// Creează o nouă adresă de contract în baza de date
func (repository *Repository) CreateContractAddress(addr *models.ContractAddress) error {
	return repository.db.Create(addr).Error
}

// Găsește o adresă de contract după ID
func (repository *Repository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := repository.db.First(&addr, id).Error
	return &addr, err
}

// Creează un nou produs în baza de date
func (repository *Repository) CreateProduct(product *models.Product) error {
	return repository.db.Create(product).Error
}

// Găsește un produs după ID
func (repository *Repository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := repository.db.First(&product, id).Error
	return &product, err
}
