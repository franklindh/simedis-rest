package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func PemeriksaanRoutes(rg *gin.RouterGroup, h *handler.PemeriksaanHandler) {
	pemeriksaanRoutes := rg.Group("/pemeriksaan")
	{
		pemeriksaanRoutes.POST("", h.Create)
		pemeriksaanRoutes.GET("/:id", h.GetByID)
		pemeriksaanRoutes.PUT("/:id", h.Update)
	}
}
