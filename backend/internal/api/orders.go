package api

import (
    "net/http"
    "strconv"
    "orders/internal/models"
    "github.com/gin-gonic/gin"
)

// Request pentru crearea comenzii (fără OwnerID, acesta vine din context)
type OrderCreateRequest struct {
    ClientID   uint    `json:"client_id"`
    ContractID uint    `json:"contract_id"`
    TotalPrice float64 `json:"total_price"`
    Status     string  `json:"status"`
}

// Handler pentru crearea comenzii (POST /orders)
func CreateOrderHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req OrderCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        userID := c.GetUint("user_id") // user_id din context, nu din JSON

        order := &models.Order{
            OwnerID:    userID,
            ClientID:   req.ClientID,
            ContractID: req.ContractID,
            TotalPrice: req.TotalPrice,
            Status:     req.Status,
        }

        if err := s.CreateOrder(userID, order); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, order)
    }
}

// Handler pentru lista comenzilor (GET /orders)
func GetOrdersHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")
        orders, err := s.FindOrdersByUserID(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, orders)
    }
}

// Handler pentru obținerea unei comenzi după id (GET /orders/:id)
func GetOrderHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }

        order, err := s.FindOrderByID(uint(id))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
            return
        }

        // Verifică dacă utilizatorul este owner-ul comenzii
        if order.OwnerID != userID {
            c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
            return
        }

        c.JSON(http.StatusOK, order)
    }
}