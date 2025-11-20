package api

import (
    "net/http"
    "strconv"

    "orders/internal/models"

    "github.com/gin-gonic/gin"
)

// Запрос для создания клиента
type ClientCreateRequest struct {
    Name    string `json:"name" binding:"required"`
    Email   string `json:"email" binding:"required,email"`
    Phone   string `json:"phone"`
    Address string `json:"address"`
    UserID  uint   `json:"user_id" binding:"required"`
}

// Обработчик создания клиента
func CreateClientHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req ClientCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        client := &models.Client{
            Name:    req.Name,
            Email:   req.Email,
            Phone:   req.Phone,
            Address: req.Address,
            UserID:  req.UserID,
        }

        if err := s.CreateClient(client); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, client)
    }
}

// Обработчик получения клиента по id
func GetClientHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }

        client, err := s.FindClientByID(uint(id))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, client)
    }
}

// Запрос для создания договора
type ContractCreateRequest struct {
    Number string  `json:"number" binding:"required"`
    Date   string  `json:"date" binding:"required"`
    Amount float64 `json:"amount" binding:"required"`
    Status string  `json:"status" binding:"required"`
}

// Создать договор для клиента (client id в path)
func CreateContractHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        clientID, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
            return
        }

        var req ContractCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        contract := &models.Contract{
            ClientID: uint(clientID),
            Number:   req.Number,
            Date:     req.Date,
            Amount:   req.Amount,
            Status:   req.Status,
        }

        if err := s.CreateContract(contract); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, contract)
    }
}

// Получить договор по id
func GetContractHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := strconv.ParseUint(idStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }

        contract, err := s.FindContractByID(uint(id))
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, contract)
    }
}