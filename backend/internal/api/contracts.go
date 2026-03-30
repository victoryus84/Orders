package api

import (
	"net/http"
	"orders/internal/models"
	"strconv"
	"strings"
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
	FiscalID string `json:"fiscal_id" xml:"fiscal_id" binding:"required"` // Pentru că 1C trimite Codul Fiscal, nu ID-ul de Client din Postgres
}

type AddressReq struct {
	ContractID uint   `json:"contract_id" xml:"contract_id" binding:"required"`
	Address    string `json:"address" xml:"address" binding:"required"`
	Type       string `json:"type" xml:"type"` // billing, shipping
}

// --- HANDLERS ---

func CreateContractHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASUL 1: Parsarea corpului cererii (JSON sau XML)
		// Folosim DTO-ul specific pentru contracte
		requests, err := ParseBody[ContractReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
			return
		}

		// Luăm ID-ul utilizatorului din token (pentru coloana owner_id din DB)
		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.Contract, 0)
		skipped := make([]map[string]string, 0)

		// PASUL 2: Logica de business pentru fiecare contract
		for _, req := range requests {
			// A. Validare de bază (Câmpurile obligatorii din 1C)
			if strings.TrimSpace(req.Number) == "" || strings.TrimSpace(req.FiscalID) == "" {
				skipped = append(skipped, map[string]string{
					"number": req.Number, 
					"reason": "missing_required_fields (number or fiscal_id)",
				})
				continue
			}

			// B. CĂUTAREA TATĂLUI (Găsim Clientul după Codul Fiscal din 1C)
			// Aici folosim funcția pe care o ai deja în Service
			client, err := s.FindClientByFiscalID(req.FiscalID)
			if err != nil {
				// Dacă nu găsim clientul, nu putem crea contractul!
				skipped = append(skipped, map[string]string{
					"number":     req.Number,
					"fiscal_id":  req.FiscalID,
					"reason":     "client_not_found",
				})
				continue
			}

			// C. Conversia de la REQ (ce vine din 1C) la MODEL (ce pleacă în Postgres)
			contract := &models.Contract{
				Number:   req.Number,
				Name:     req.Name,
				Date:     req.Date,
				Amount:   req.Amount,
				Status:   req.Status,
				ClientID: client.ID, // <--- Aici e "magia": ID-ul de Postgres al clientului
				OwnerID:  ownerID,   // Coloana owner_id pe care o aveai în DB
			}

			// D. Salvarea efectivă prin Service
			if err := s.CreateContract(contract); err != nil {
				skipped = append(skipped, map[string]string{
					"number": req.Number, 
					"reason": "db_save_error: " + err.Error(),
				})
				continue
			}
			created = append(created, contract)
		}

		// PASUL 3: Răspunsul final către 1C / Frontend
		c.JSON(http.StatusCreated, gin.H{
			"created": created,
			"skipped": skipped,
			"count":   len(created),
		})
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
