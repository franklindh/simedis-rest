package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func JadwalRoutes(rg *gin.RouterGroup, h *handler.JadwalHandler) {
	jadwalRoutes := rg.Group("/jadwal")
	{
		jadwalRoutes.GET("", h.GetAll)
		jadwalRoutes.POST("", h.Create)
		jadwalRoutes.GET("/:id", h.GetByID)
		jadwalRoutes.PUT("/:id", h.Update)
		jadwalRoutes.DELETE("/:id", h.Delete)
	}
}
