package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func PetugasRoutes(rg *gin.RouterGroup, h *handler.PetugasHandler) {
	petugasRoutes := rg.Group("/petugas")
	{
		petugasRoutes.GET("", h.GetAll)
		petugasRoutes.POST("", h.Create)
		petugasRoutes.GET("/:id", h.GetByID)
		petugasRoutes.PUT("/:id", h.Update)
		petugasRoutes.DELETE("/:id", h.Delete)
	}
}
