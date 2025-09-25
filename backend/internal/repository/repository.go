package repository

import (
	"orders/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) CreateUser(user *models.User) error {
	return repository.db.Create(user).Error
}

func (repository *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repository.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

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

func (repository *Repository) CreateClient(client *models.Client) error {
	return repository.db.Create(client).Error
}

func (repository *Repository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := repository.db.First(&client, id).Error
	return &client, err
}

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

func (repository *Repository) CreateProduct(product *models.Product) error {
	return repository.db.Create(product).Error
}

func (repository *Repository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := repository.db.First(&product, id).Error
	return &product, err
}
