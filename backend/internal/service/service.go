package service

import (
	"fmt"
	"orders/internal/config"
	"orders/internal/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	// User methods
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)

	// Order methods
	CreateOrder(order *models.Order) error
	FindOrdersByUserID(userID uint) ([]models.Order, error)
	FindOrderByID(id uint) (*models.Order, error)

	// Product methods
	FindProductByID(id uint) (*models.Product, error)
	FindVatTaxByID(id uint) (*models.VatTax, error)
	FindUnitByID(id uint) (*models.Unit, error)

	// Client methods
	CreateClient(client *models.Client) error
	GetFirst1000Clients() ([]models.Client, error)
	FindClientsByQuery(query string) ([]models.Client, error)
	FindClientByID(id uint) (*models.Client, error)
	FindClientByFiscalID(fiscalID string) (*models.Client, error)

	// Contract methods
	CreateContract(contract *models.Contract) error
	FindContractByID(id uint) (*models.Contract, error)
	CreateContractAddress(addr *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)

	// Product methods
	CreateProduct(product *models.Product) error
}

type Service struct {
	repository Repository
	jwtSecret  string
	cfg        *config.Config // Добавляем конфигурацию
}

func NewService(repository Repository, jwtSecret string) *Service {
	cfg := config.Load()
	return &Service{repository: repository, jwtSecret: jwtSecret, cfg: &cfg}
}

func (service *Service) Signup(email, password, role string) error {
	if !service.cfg.Allowsignup {
		return fmt.Errorf("user registration is disabled")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     map[bool]string{true: "admin", false: "user"}[strings.ToLower(role) == "trueadmin"],
	}
	return service.repository.CreateUser(user)
}

func (service *Service) Login(email, password string) (string, error) {
	user, err := service.repository.FindUserByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	return token.SignedString([]byte(service.jwtSecret))
}

func (service *Service) CreateOrder(userID uint, order *models.Order) error {
	total := 0.0
	for i := range order.OrderItems {
		product, err := service.repository.FindProductByID(order.OrderItems[i].ProductID)
		if err != nil {
			return err
		}
		order.OrderItems[i].Price = product.Price
		total += product.Price * float64(order.OrderItems[i].Quantity)
	}
	order.OwnerID = userID
	order.TotalPrice = total
	order.Status = "pending"
	return service.repository.CreateOrder(order)
}

func (service *Service) FindOrderByID(id uint) (*models.Order, error) {
	return service.repository.FindOrderByID(id)
}

func (service *Service) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	return service.repository.FindOrdersByUserID(userID)
}

// Clients methods
func (service *Service) CreateClient(client *models.Client) error {
	return service.repository.CreateClient(client)
}
func (service *Service) GetFirst1000Clients() ([]models.Client, error) {
	return service.repository.GetFirst1000Clients()
}

func (service *Service) FindClientsByQuery(query string) ([]models.Client, error) {
	return service.repository.FindClientsByQuery(query)
}

func (service *Service) FindClientByID(id uint) (*models.Client, error) {
	return service.repository.FindClientByID(id)
}

func (service *Service) FindClientByFiscalID(fiscalID string) (*models.Client, error) {
	return service.repository.FindClientByFiscalID(fiscalID)
}

// Contract methods
func (service *Service) CreateContract(contract *models.Contract) error {
	return service.repository.CreateContract(contract)
}

func (service *Service) FindContractByID(id uint) (*models.Contract, error) {
	return service.repository.FindContractByID(id)
}

// ContractAddress methods
func (service *Service) CreateContractAddress(addr *models.ContractAddress) error {
	return service.repository.CreateContractAddress(addr)
}

func (service *Service) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	return service.repository.FindContractAddressByID(id)
}

// Product methods
func (service *Service) CreateProduct(product *models.Product) error {
	return service.repository.CreateProduct(product)
}

func (service *Service) FindProductByID(id uint) (*models.Product, error) {
	return service.repository.FindProductByID(id)
}

func (service *Service) FindVatTaxByID(id uint) (*models.VatTax, error) {
	return service.repository.FindVatTaxByID(id)
}

func (service *Service) FindUnitByID(id uint) (*models.Unit, error) {
	return service.repository.FindUnitByID(id)
}
