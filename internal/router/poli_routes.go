package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func PoliRoutes(rg *gin.RouterGroup, h *handler.PoliHandler) {
	poliRoutes := rg.Group("/poli")
	{
		poliRoutes.GET("", h.GetAll)
		poliRoutes.POST("", h.Create)
		poliRoutes.GET("/:id", h.GetByID)
		poliRoutes.PUT("/:id", h.Update)
		poliRoutes.DELETE("/:id", h.Delete)
	}
}
