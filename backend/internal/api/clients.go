package api

import (
    "net/http"
    "strconv"
    "orders/internal/models"
    "github.com/gin-gonic/gin"
)

// Handler pentru crearea clientului
type ClientCreateRequest struct {
    Name     string `json:"name" binding:"required"`
    FiscalID string `json:"fiscal_code"`
    Email    string `json:"email" binding:"required,email"`
    Phone    string `json:"phone"`
    Address  string `json:"address"`
    UserID   uint   `json:"user_id" binding:"required"`
}

func CreateClientHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req ClientCreateRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        ownerID := c.GetUint("user_id")

        client := &models.Client{
            Name:     req.Name,
            FiscalID: req.FiscalID,
            Email:    req.Email,
            Phone:    req.Phone,
            Address:  req.Address,
            OwnerID:  ownerID,
        }

        if err := s.CreateClient(client); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, client)
    }
}

// Handler pentru obținerea primilor 1000 de clienți
func GetFirst1000Clients(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        clients, err := s.GetFirst1000Clients()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, clients)
    }
}   

// Handler pentru căutarea clienților după query
func SearchClientsHandler(s Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        query := c.Query("q")
        if query == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
            return
        }

        clients, err := s.FindClientsByQuery(query)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, clients)
    }
}

// Handler pentru obținerea clientului după id
func GetClientsHandler(s Service) gin.HandlerFunc {
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
// Handler pentru obținerea clientului după id
func GetClientByIDHandler(s Service) gin.HandlerFunc {
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

// Request pentru crearea contractului
type ContractCreateRequest struct {
    Number string  `json:"number" binding:"required"`
    Date   string  `json:"date" binding:"required"`
    Amount float64 `json:"amount" binding:"required"`
    Status string  `json:"status" binding:"required"`
}

// Crearea contractului pentru client (client id în path)
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

// Obține contractul după id
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

