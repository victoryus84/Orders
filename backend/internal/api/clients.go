package api

import (
	"errors"
	"net/http"
	"orders/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler pentru crearea clientului
type ClientCreateReq struct {
	ClientTypeID uint   `json:"client_type" xml:"client_type" binding:"required"`
	Name         string `json:"name" xml:"name" binding:"required"`
	FiscalID     string `json:"fiscal_code" xml:"fiscal_code" binding:"required"`
	// Email is optional for now; accept empty or placeholder values until the DB holds actual emails
	Email   string `json:"email" xml:"email" binding:"omitempty"`
	Phone   string `json:"phone" xml:"phone"`
	Address string `json:"address" xml:"address"`
}

func CreateClientHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASUL 1: Folosim procedura "ParseBody" pentru parsare
		// Această singură linie înlocuiește tot blocul tău mare de IF-uri (JSON/XML/Array)
		requests, err := ParseBody[ClientCreateReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
			return
		}

		created := make([]*models.Client, 0)
		skipped := make([]map[string]string, 0)

		// PASUL 2: Logica ta specifică de Business
		for _, req := range requests {
			// A. Validare de bază
			if req.ClientTypeID == 0 || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.FiscalID) == "" {
				skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "missing_required_fields"})
				continue
			}

			// B. Verificare duplicate (Logica ta de aur)
			existing, err := s.FindClientByFiscalID(req.FiscalID)
			if err == nil && existing != nil {
				skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "duplicate"})
				continue
			}
			// Dacă eroarea nu e "Not Found", înseamnă că e ceva grav la baza de date
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
				return
			}

			// C. Sanitizarea email-ului (Asta e foarte deșteaptă!)
			email := strings.TrimSpace(req.Email)
			lowEmail := strings.ToLower(email)
			if lowEmail == "not inserted" || lowEmail == "n/a" || lowEmail == "none" || lowEmail == "" {
				email = "" // O lăsăm goală dacă e placeholder
			}

			// D. Conversia de la REQ la MODEL
			client := &models.Client{
				ClientTypeID: req.ClientTypeID,
				Name:         req.Name,
				FiscalID:     req.FiscalID,
				Email:        email,
				Phone:        req.Phone,
				Address:      req.Address,
			}

			// E. Salvarea efectivă
			if err := s.CreateClient(client); err != nil {
				skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": err.Error()})
				continue
			}
			created = append(created, client)
		}

		// PASUL 3: Răspunsul final
		c.JSON(http.StatusCreated, gin.H{
			"created": created,
			"skipped": skipped,
			"count":   len(created),
		})
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
