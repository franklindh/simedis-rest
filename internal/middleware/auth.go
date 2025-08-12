package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}", nil)
			c.Abort()
			return
		}
		tokenString := parts[1]

		jwtSecret := []byte(cfg.JWTSecret)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["sub"])
			c.Set("role", claims["role"])
			c.Set("username", claims["username"])
		} else {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
