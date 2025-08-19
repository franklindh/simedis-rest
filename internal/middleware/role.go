package middleware

import (
	"net/http"

	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userRole, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "User role not found in token", nil)
			c.Abort()
			return
		}

		isAllowed := false
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to access this resource", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
