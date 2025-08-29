package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PetugasRoutes(rg *gin.RouterGroup, h *handler.PetugasHandler) {
	petugasRoutes := rg.Group("/petugas")
	{
		petugasRoutes.GET("", h.GetAll)
		petugasRoutes.GET("/:id", h.GetByID)
		petugasRoutes.PUT("/change-password", h.ChangePassword)

		user := petugasRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
