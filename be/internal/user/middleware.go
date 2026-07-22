package user

import (
	"errors"
	"os"
	"strings"

	"icon_exchange/internal/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	// ContextUserIDKey is the key used to store the user ID in the Gin context.
	ContextUserIDKey = "userID"
)

// AuthMiddleware validates the JWT access token in the Authorization header.
// If valid, it stores the userID (subclaim) in the request context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		// Convert sub (which is float64 when decoded from JSON/claims) to uint
		var userID uint
		switch v := sub.(type) {
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		case uint:
			userID = v
		default:
			respondError(c, shared.ErrUnauthenticated)
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set(ContextUserIDKey, userID)
		c.Next()
	}
}

// GetUserIDFromContext retrieves the user ID from the Gin context.
// Returns 0 and false if the user ID is not found or invalid.
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	val, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0, false
	}
	userID, ok := val.(uint)
	return userID, ok
}
