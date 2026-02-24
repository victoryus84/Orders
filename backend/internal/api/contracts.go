package api

import (
	"net/http"
	"orders/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateContractHandler gestionează POST /contracts
func CreateContractHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var contract models.Contract
		if err := c.ShouldBindJSON(&contract); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.CreateContract(&contract); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, contract)
	}
}

// GetContractByIDHandler gestionează GET /contracts/:id
func GetContractByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		contract, err := s.FindContractByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}

		c.JSON(http.StatusOK, contract)
	}
}

// CreateContractAddressHandler gestionează POST /contract_addresses
func CreateContractAddressHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var addr models.ContractAddress
		if err := c.ShouldBindJSON(&addr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.CreateContractAddress(&addr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, addr)
	}
}

// GetContractAddressByIDHandler gestionează GET /contract_addresses/:id
func GetContractAddressByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		addr, err := s.FindContractAddressByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract address not found"})
			return
		}

		c.JSON(http.StatusOK, addr)
	}
}
