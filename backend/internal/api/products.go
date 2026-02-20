package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"orders/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler pentru crearea produsului
type ProductRequest struct {
	Name        string  `json:"name" xml:"name" binding:"required"`
	Price       float64 `json:"price" xml:"price" binding:"required"`
	Description string  `json:"description" xml:"description"`
	UnitID      uint    `json:"unit_id" xml:"unit_id" binding:"required"`
	VatTaxID    uint    `json:"vat_tax_id" xml:"vat_tax_id" binding:"required"`
}

func CreateProductHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.ContentType()
		var requests []ProductRequest

		// Read raw body so we can parse deterministically and return helpful errors
		data, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
			return
		}

		// Restore body for potential downstream handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		if contentType == "application/xml" || contentType == "text/xml" {
			// Try to unmarshal wrapper <products><product>...</product></products>
			var wrapper struct {
				XMLName  xml.Name         `xml:"products"`
				Products []ProductRequest `xml:"product"`
			}
			if err := xml.Unmarshal(data, &wrapper); err == nil && len(wrapper.Products) > 0 {
				requests = wrapper.Products
			} else {
				// Fallback to single <product>...</product>
				var single ProductRequest
				if err := xml.Unmarshal(data, &single); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML: " + err.Error(), "body": string(data)})
					return
				}
				requests = append(requests, single)
			}
		} else {
			// JSON: try array, wrapper {"products": [...]}, then single object
			var arr []ProductRequest
			if err := json.Unmarshal(data, &arr); err == nil && len(arr) > 0 {
				requests = arr
			} else {
				var wrapper struct {
					Products []ProductRequest `json:"products"`
				}
				if err := json.Unmarshal(data, &wrapper); err == nil && len(wrapper.Products) > 0 {
					requests = wrapper.Products
				} else {
					var single ProductRequest
					if err := json.Unmarshal(data, &single); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error(), "body": string(data)})
						return
					}
					requests = append(requests, single)
				}
			}
		}

		created := make([]*models.Product, 0)
		skipped := make([]map[string]string, 0)

		for _, req := range requests {
			// Basic validation (these fields are required)
			if strings.TrimSpace(req.Name) == "" || req.Price <= 0 || req.UnitID == 0 || req.VatTaxID == 0 {
				skipped = append(skipped, map[string]string{"name": req.Name, "reason": "missing_required_fields_or_invalid_values"})
				continue
			}

			if _, err := s.FindUnitByID(req.UnitID); err != nil {
				skipped = append(skipped, map[string]string{"name": req.Name, "reason": "invalid_unit_id"})
				continue
			}
			if _, err := s.FindVatTaxByID(req.VatTaxID); err != nil {
				skipped = append(skipped, map[string]string{"name": req.Name, "reason": "invalid_vat_tax_id"})
				continue
			}
			product := &models.Product{
				Name:        req.Name,
				Price:       req.Price,
				Description: req.Description,
				UnitID:      req.UnitID,
				VatTaxID:    req.VatTaxID,
			}

			if err := s.CreateProduct(product); err != nil {
				skipped = append(skipped, map[string]string{"name": req.Name, "reason": err.Error()})
				continue
			}
			created = append(created, product)
		}

		c.JSON(http.StatusCreated, gin.H{"created": created, "skipped": skipped})
	}
}
