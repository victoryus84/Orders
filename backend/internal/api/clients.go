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

type ClientReq struct {
	ClientTypeID  uint   `json:"client_type" xml:"client_type" binding:"required"`
	Name          string `json:"name" xml:"name" binding:"required"`
	FiscalID      string `json:"fiscal_id" xml:"fiscal_id" binding:"required"`
	Email         string `json:"email" xml:"email" binding:"omitempty"`
	Phone         string `json:"phone" xml:"phone"`
	FiscalAddress string `json:"fiscal_address" xml:"fiscal_address"`
	PostalAddress string `json:"postal_address" xml:"postal_address"`
}

type ClientAddressReq struct {
	FiscalID string `json:"fiscal_id" xml:"fiscal_id" binding:"required"`
	Name     string `json:"name" xml:"name" binding:"required"`
	Address  string `json:"address" xml:"address" binding:"required"`
	Type     string `json:"type" xml:"type"` // billing, shipping
}

func CreateClientHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASUL 1: Folosim procedura "ParseBody" pentru parsare
		// Această singură linie înlocuiește tot blocul tău mare de IF-uri (JSON/XML/Array)
		requests, err := ParseBody[ClientReq](c)
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
			rawEmail := strings.TrimSpace(req.Email)
			lowEmail := strings.ToLower(rawEmail)
			var emailPtr *string
			// Dacă e email pe bune, salvăm adresa lui
			if lowEmail != "" && lowEmail != "not inserted" && lowEmail != "n/a" && lowEmail != "none" {
				emailPtr = &rawEmail
			}

			// D. Conversia de la REQ la MODEL
			client := &models.Client{
				ClientTypeID:  req.ClientTypeID,
				Name:          req.Name,
				FiscalID:      req.FiscalID,
				Email:         emailPtr,
				Phone:         req.Phone,
				FiscalAddress: req.FiscalAddress,
				PostalAddress: req.PostalAddress,
			}

			// E. Salvarea efectivă
			if err := s.CreateClient(client); err != nil {
				skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": err.Error()})
				continue
			}
			created = append(created, client)
		}

		// Pregătim o listă scurtă cu erorile (doar primele 10, să nu omorâm 1C-ul)
		shortSkipped := skipped
		if len(skipped) > 20 {
			shortSkipped = skipped[:20]
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":         "success",
			"total_created":  len(created),
			"total_skipped":  len(skipped),
			"errors_preview": shortSkipped, // Trimitem doar o mostră de erori
			"message":        "Import clients finalizat cu succes",
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

// Client address handlers
func CreateClientAddressHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		requests, err := ParseBody[ClientAddressReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
			return
		}
		// Luăm ID-ul utilizatorului din token (pentru coloana owner_id din DB)
		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.ClientAddress, 0)
		skipped := make([]map[string]string, 0)

		// PASUL 2: Logica de business pentru fiecare contract
		for _, req := range requests {
			// A. Validare detaliată
			// Verificăm dacă lipsește Codul Fiscal
			if strings.TrimSpace(req.FiscalID) == "" {
				skipped = append(skipped, map[string]string{
					"number": req.Address,
					"reason": "EROARE: Campul 'fiscal_id' a ajuns GOL. Clientul nu poate fi gasit!",
				})
				continue
			}

			// B. Căutarea clientului (rămâne la fel)
			client, err := s.FindClientByFiscalID(req.FiscalID)
			if err != nil {
				skipped = append(skipped, map[string]string{
					"number": req.Address,
					"reason": "Clientul cu FiscalID " + req.FiscalID + " nu exista in baza de date!",
				})
				continue
			}

			// --- LOGICA PENTRU NUMBER (NUMAR SAU NULL) ---
			rawaddress := strings.TrimSpace(req.Address)
			var addressPtr *string
			// Dacă e număr pe bune, salvăm numărul
			if rawaddress != "" {
				addressPtr = &rawaddress
			}

			// C. Conversia de la REQ (ce vine din 1C) la MODEL (ce pleacă în Postgres)
			addr := &models.ClientAddress{
				ClientID: client.ID,
				Name:  	  req.Name,
				Address:  addressPtr,
				Type:     req.Type,
				OwnerID:  ownerID,
			}
			if err := s.CreateClientAddress(addr); err == nil {
				created = append(created, addr)
			}
		}
		c.JSON(http.StatusCreated, created)
	}
}
