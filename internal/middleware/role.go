package middleware

import (
	"fmt"
	"net/http"

	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("\n--- ROLE MIDDLEWARE CHECK ---")
		fmt.Printf("Endpoint: %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Printf("Role yang diizinkan (roles): %v\n", roles)

		userRole, exists := c.Get("role")
		if !exists {
			fmt.Println("HASIL: GAGAL! Role tidak ditemukan di context.")
			fmt.Println("---------------------------")
			utils.ErrorResponse(c, http.StatusForbidden, "User role not found in token", nil)
			c.Abort()
			return
		}

		fmt.Printf("Role user dari token (userRole): '%v'\n", userRole)

		isAllowed := false
		for _, allowedRole := range roles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		fmt.Printf("Apakah diizinkan? (isAllowed): %t\n", isAllowed)
		fmt.Println("---------------------------")
		// --- AKHIR DARI ALAT SADAP ---

		if !isAllowed {
			utils.ErrorResponse(c, http.StatusForbidden, "You are not authorized to access this resource", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
