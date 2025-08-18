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

		adminOnly := petugasRoutes.Group("")
		adminOnly.Use(middleware.Authorize("Administrasi"))
		{
			adminOnly.POST("", h.Create)
			adminOnly.PUT("/:id", h.Update)
			adminOnly.DELETE("/:id", h.Delete)
		}
	}
}
