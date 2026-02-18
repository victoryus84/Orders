package api

import (
	"fmt"
)

func echoProductID(id uint) string {
	return fmt.Sprintf("Product ID: %d", id)
}

// Handler pentru căutarea produselor după query
