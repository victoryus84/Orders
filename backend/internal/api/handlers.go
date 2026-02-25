package api

import (
	"net/http"
	"orders/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	// Authentication methods
	Signup(email, password, role string) error
	Login(email, password string) (string, error)

	// Order methods
	CreateOrder(userID uint, order *models.Order) error
	FindOrdersByUserID(userID uint) ([]models.Order, error)
	FindOrderByID(id uint) (*models.Order, error)

	// Client methods
	CreateClient(client *models.Client) error
	FindClientByID(id uint) (*models.Client, error)
	FindClientByFiscalID(fiscalID string) (*models.Client, error)
	GetFirst1000Clients() ([]models.Client, error)
	FindClientsByQuery(query string) ([]models.Client, error)

	// Contract methods
	CreateContract(contract *models.Contract) error
	FindContractByID(id uint) (*models.Contract, error)
	CreateContractAddress(addr *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)

	// Product methods
	CreateProduct(product *models.Product) error
	FindProductByID(id uint) (*models.Product, error)
	FindVatTaxByID(id uint) (*models.VatTax, error)
	FindUnitByID(id uint) (*models.Unit, error)
	FindProductGroupByID(id uint) (*models.ProductGroup, error)
}

func SetupRoutes(router *gin.Engine, service Service) {

	// --- Health-check ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	// --- Auth ---
	router.POST("/signup", func(context *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := service.Signup(req.Email, req.Password, req.Role); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"message": "User created"})
	})

	router.POST("/login", func(context *gin.Context) {
		// 1. Folosim procedura universală care știe JSON și XML
		// Definim o structură locală sau folosim una din modele
		type LoginReq struct {
			Email    string `json:"email" xml:"email"`
			Password string `json:"password" xml:"password"`
		}

		requests, err := ParseBody[LoginReq](context)
		if err != nil || len(requests) == 0 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Date invalide sau body gol"})
			return
		}

		// Luăm prima cerere din listă
		req := requests[0]

		// 2. Logica de login
		token, err := service.Login(req.Email, req.Password)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		context.JSON(http.StatusOK, gin.H{"token": token})
	})

	// API v1 routes with prefix
	api := router.Group("/api/v1")
	protected := api.Group("/").Use(authMiddleware())
	{

		// --- Orders ---
		protected.POST("/orders", CreateOrderHandler(service))
		protected.GET("/orders", GetOrdersHandler(service))
		protected.GET("/orders/:id", GetOrderHandler(service))

		// --- Clients ---
		protected.POST("/clients", CreateClientHandler(service))
		protected.GET("/clients", GetFirst1000Clients(service))
		protected.GET("/clients/search", SearchClientsHandler(service))
		protected.GET("/clients/:id", GetClientByIDHandler(service))

		// --- Contracts ---
		protected.POST("/contracts", CreateContractHandler(service))
		protected.GET("/contracts/:id", GetContractByIDHandler(service))

		// --- ContractAddresses ---
		protected.POST("/contract_addresses", CreateContractAddressHandler(service))
		protected.GET("/contract_addresses/:id", GetContractAddressByIDHandler(service))

		// --- Products ---
		protected.POST("/products", CreateProductHandler(service))
		protected.GET("/products/:id", GetProductByIDHandler(service))

	}
}
