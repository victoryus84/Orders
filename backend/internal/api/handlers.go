package api

import (
	"net/http"
	"orders/internal/models"

	//"orders/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Signup(email, password, role string) error
	Login(email, password string) (string, error)
	CreateOrder(userID uint, order *models.Order) error
	FindOrdersByUserID(userID uint) ([]models.Order, error)
	FindOrderByID(id uint) (*models.Order, error)
	CreateClient(client *models.Client) error
	FindClientByID(id uint) (*models.Client, error)
	CreateContract(contract *models.Contract) error
	FindContractByID(id uint) (*models.Contract, error)
	CreateContractAddress(addr *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)
	CreateProduct(product *models.Product) error
	FindProductByID(id uint) (*models.Product, error)
}

func SetupRoutes(router *gin.Engine, service Service) {
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
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := context.ShouldBindJSON(&req); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := service.Login(req.Email, req.Password)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		context.JSON(http.StatusOK, gin.H{"token": token})
	})

	protected := router.Group("/").Use(authMiddleware())
	{
		protected.POST("/orders", func(context *gin.Context) {
			userID := context.GetUint("user_id")
			var order models.Order
			if err := context.ShouldBindJSON(&order); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := service.CreateOrder(userID, &order); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, order)
		})

		protected.GET("/orders", func(context *gin.Context) {
			userID := context.GetUint("user_id")
			orders, err := service.FindOrdersByUserID(userID)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, orders)
		})

		protected.GET("/orders/:id", func(context *gin.Context) {
			userID := context.GetUint("user_id")
			id, _ := strconv.Atoi(context.Param("id"))
			order, err := service.FindOrderByID(uint(id))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
				return
			}
			if order.UserID != userID {
				context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
				return
			}
			context.JSON(http.StatusOK, order)
		})

		// --- Clients ---
		protected.POST("/clients", func(context *gin.Context) {
			var client models.Client
			if err := context.ShouldBindJSON(&client); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := service.CreateClient(&client); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, client)
		})

		protected.GET("/clients/:id", func(context *gin.Context) {
			id, _ := strconv.Atoi(context.Param("id"))
			client, err := service.FindClientByID(uint(id))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
				return
			}
			context.JSON(http.StatusOK, client)
		})

		// --- Contracts ---
		protected.POST("/contracts", func(context *gin.Context) {
			var contract models.Contract
			if err := context.ShouldBindJSON(&contract); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := service.CreateContract(&contract); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, contract)
		})

		protected.GET("/contracts/:id", func(context *gin.Context) {
			id, _ := strconv.Atoi(context.Param("id"))
			contract, err := service.FindContractByID(uint(id))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
				return
			}
			context.JSON(http.StatusOK, contract)
		})

		// --- ContractAddresses ---
		protected.POST("/contract_addresses", func(context *gin.Context) {
			var addr models.ContractAddress
			if err := context.ShouldBindJSON(&addr); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := service.CreateContractAddress(&addr); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, addr)
		})

		protected.GET("/contract_addresses/:id", func(context *gin.Context) {
			id, _ := strconv.Atoi(context.Param("id"))
			addr, err := service.FindContractAddressByID(uint(id))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Contract address not found"})
				return
			}
			context.JSON(http.StatusOK, addr)
		})

		// --- Products ---
		protected.POST("/products", func(context *gin.Context) {
			var product models.Product
			if err := context.ShouldBindJSON(&product); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := service.CreateProduct(&product); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			context.JSON(http.StatusOK, product)
		})

		protected.GET("/products/:id", func(context *gin.Context) {
			id, _ := strconv.Atoi(context.Param("id"))
			product, err := service.FindProductByID(uint(id))
			if err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			context.JSON(http.StatusOK, product)
		})
	}
}
