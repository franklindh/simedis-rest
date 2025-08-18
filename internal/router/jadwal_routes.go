package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func JadwalRoutes(rg *gin.RouterGroup, h *handler.JadwalHandler) {
	jadwalRoutes := rg.Group("/jadwal")
	{
		jadwalRoutes.GET("", h.GetAll)
		jadwalRoutes.GET("/:id", h.GetByID)

		user := jadwalRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
