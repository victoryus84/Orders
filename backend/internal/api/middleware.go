package api

import (
	"net/http"
	"strings"
	"orders/internal/config"


	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Конфигурация
		cfg := config.Load()

		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			context.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			context.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			context.Abort()
			return
		}

		context.Set("user_id", uint(claims["user_id"].(float64)))
		context.Set("role", claims["role"].(string))
		context.Next()
	}
}