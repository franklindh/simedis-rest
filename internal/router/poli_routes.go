package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PoliRoutes(rg *gin.RouterGroup, h *handler.PoliHandler) {
	poliRoutes := rg.Group("/poli")
	{
		poliRoutes.GET("", h.GetAll)
		poliRoutes.GET("/:id", h.GetByID)

		user := poliRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
