package api

import (
	"net/http"
	"orders/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- DTOs (Data Transfer Objects) ---

type ContractReq struct {
	Number   string  `json:"number" xml:"number" binding:"required"`
	Name     string  `json:"name" xml:"name" binding:"required"`
	Date     string  `json:"date" xml:"date" binding:"required"` // Format YYYY-MM-DD
	Amount   float64 `json:"amount" xml:"amount"`
	ClientID uint    `json:"client_id" xml:"client_id" binding:"required"`
	Status   string  `json:"status" xml:"status"`
}

type AddressReq struct {
	ContractID uint   `json:"contract_id" xml:"contract_id" binding:"required"`
	Address    string `json:"address" xml:"address" binding:"required"`
	Type       string `json:"type" xml:"type"` // billing, shipping
}

// --- HANDLERS ---

func CreateContractHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASUL 1: Folosim procedura "ParseBody" pentru parsare
		// Această singură linie înlocuiește tot blocul tău mare de IF-uri (JSON/XML/Array)
		requests, err := ParseBody[ContractReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format (JSON/XML list or object required)"})
			return
		}

		// Luăm ID-ul celui care face cererea (din token-ul de auth)
		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.Contract, 0)
		errors := make([]map[string]string, 0)

		for _, req := range requests {
			contract := &models.Contract{
				Number:   req.Number,
				Name:     req.Name,
				Date:     req.Date,
				Amount:   req.Amount,
				ClientID: req.ClientID,
				Status:   req.Status,
				OwnerID:  ownerID, // Foarte important pentru baza de date!
			}

			if err := s.CreateContract(contract); err != nil {
				errors = append(errors, map[string]string{"number": req.Number, "error": err.Error()})
				continue
			}
			created = append(created, contract)
		}

		c.JSON(http.StatusCreated, gin.H{"created": created, "errors": errors})
	}
}

func CreateContractAddressHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		requests, err := ParseBody[AddressReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
			return
		}

		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.ContractAddress, 0)
		for _, req := range requests {
			addr := &models.ContractAddress{
				ContractID: req.ContractID,
				Address:    req.Address,
				Type:       req.Type,
				OwnerID:    ownerID,
			}
			if err := s.CreateContractAddress(addr); err == nil {
				created = append(created, addr)
			}
		}
		c.JSON(http.StatusCreated, created)
	}
}

func GetContractByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := s.FindContractByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetContractAddressByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := s.FindContractAddressByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
