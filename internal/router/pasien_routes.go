package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func PasienRoutes(rg *gin.RouterGroup, h *handler.PasienHandler) {
	pasienRoutes := rg.Group("/pasien")
	{
		pasienRoutes.GET("", h.GetAll)
		pasienRoutes.POST("", h.Create)
		pasienRoutes.GET("/:id", h.GetByID)
		pasienRoutes.PUT("/:id", h.Update)
		pasienRoutes.DELETE("/:id", h.Delete)
	}
}
