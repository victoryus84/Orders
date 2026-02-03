package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"orders/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler pentru crearea clientului
type ClientCreateRequest struct {
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
		contentType := c.ContentType()
		var requests []ClientCreateRequest

		// Read raw body so we can parse deterministically and return helpful errors
		data, _ := c.GetRawData()
		// Restore body for potential downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		if contentType == "application/xml" || contentType == "text/xml" {
			// Try to unmarshal wrapper <clients><client>...</client></clients>
			var wrapper struct {
				XMLName xml.Name              `xml:"clients"`
				Clients []ClientCreateRequest `xml:"client"`
			}
			if err := xml.Unmarshal(data, &wrapper); err == nil && len(wrapper.Clients) > 0 {
				requests = wrapper.Clients
			} else {
				// Fallback to single <client>...</client>
				var single ClientCreateRequest
				if err := xml.Unmarshal(data, &single); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML: " + err.Error(), "body": string(data)})
					return
				}
				requests = append(requests, single)
			}
		} else {
			// JSON: try array, wrapper {"clients": [...]}, then single object
			var arr []ClientCreateRequest
			if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
				requests = arr
			} else {
				var wrapper struct {
					Clients []ClientCreateRequest `json:"clients"`
				}
				if err := json.Unmarshal(data, &wrapper); err == nil && len(wrapper.Clients) > 0 {
					requests = wrapper.Clients
				} else {
					var single ClientCreateRequest
					if err := json.Unmarshal(data, &single); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error(), "body": string(data)})
						return
					}
					requests = append(requests, single)
				}
			}
		}

		created := make([]*models.Client, 0, len(requests))
		for _, req := range requests {
			// Sanitize placeholder emails so validation won't reject requests like <email>not inserted</email>
			em := strings.TrimSpace(req.Email)
			if em != "" {
				el := strings.ToLower(em)
				if el == "not inserted" || el == "not_inserted" || el == "n/a" || el == "none" {
					req.Email = ""
				}
			}

			client := &models.Client{
				ClientTypeID: req.ClientTypeID,
				Name:         req.Name,
				FiscalID:     req.FiscalID,
				Email:        req.Email,
				Phone:        req.Phone,
				Address:      req.Address,
			}
			if err := s.CreateClient(client); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			created = append(created, client)
		}

		c.JSON(http.StatusCreated, created)
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
	Number string  `json:"number" xml:"number" binding:"required"`
	Date   string  `json:"date" xml:"date" binding:"required"`
	Amount float64 `json:"amount" xml:"amount" binding:"required"`
	Status string  `json:"status" xml:"status" binding:"required"`
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
		contentType := c.ContentType()

		// Handle both JSON and XML based on Content-Type header
		if contentType == "application/xml" || contentType == "text/xml" {
			if err := c.ShouldBindXML(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML: " + err.Error()})
				return
			}
		} else {
			// Default to JSON
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
				return
			}
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
