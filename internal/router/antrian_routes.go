package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func AntrianRoutes(rg *gin.RouterGroup, h *handler.AntrianHandler) {
	antrianRoutes := rg.Group("/antrian")
	{
		antrianRoutes.GET("", h.GetAll)
		antrianRoutes.POST("", h.Create)
		antrianRoutes.GET("/:id", h.GetByID)
		antrianRoutes.PUT("/:id", h.Update)
		antrianRoutes.DELETE("/:id", h.Delete)
	}
}
