package api

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// ParseBody este pentru parsarea cererilor API.
// T este un tip generic (Contract, Client, Order etc.)
func ParseBody[T any](c *gin.Context) ([]T, error) {
	var result []T
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	contentType := c.GetHeader("Content-Type")

	// Detectăm dacă e XML
	if strings.Contains(contentType, "xml") {
		// 1. Încercăm varianta de listă cu wrapper <items><item>...</item></items>
		// Folosim un struct anonim pentru a nu polua codul
		var wrapper struct {
			Items []T `xml:"item"`
		}
		if err := xml.Unmarshal(data, &wrapper); err == nil && len(wrapper.Items) > 0 {
			return wrapper.Items, nil
		}
		// 2. Încercăm obiect singur <item>...</item>
		var single T
		if err := xml.Unmarshal(data, &single); err == nil {
			return []T{single}, nil
		}
	}
	// Detectăm dacă e JSON (default)
	if strings.Contains(contentType, "json") {
		// Detectăm dacă e JSON (default)
		// 1. Încercăm listă [{}, {}]
		if err := json.Unmarshal(data, &result); err == nil {
			return result, nil
		}
		// 2. Încercăm obiect singur {}
		var single T
		if err := json.Unmarshal(data, &single); err == nil {
			return []T{single}, nil
		}
	}

	return nil, err
}
