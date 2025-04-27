package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract "sub" as a string
		userIDStr, okID := (*claims)["sub"].(string)
		if !okID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token data, could not get ID"})
			c.Abort()
			return
		}

		// Convert the string to uint
		userID, err := strconv.ParseUint(userIDStr, 10, 32) // 32-bit uint
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// Store the user ID as uint in context
		c.Set("user_id", uint(userID)) // Store as uint
		c.Next()
	}
}
