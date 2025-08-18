package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PemeriksaanRoutes(rg *gin.RouterGroup, h *handler.PemeriksaanHandler) {
	pemeriksaanRoutes := rg.Group("/pemeriksaan")
	{
		pemeriksaanRoutes.GET("/:id", h.GetByID)

		user := pemeriksaanRoutes.Group("")
		user.Use(middleware.Authorize("Dokter", "Poliklinik"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
		}
	}
}
