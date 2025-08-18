package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AntrianRoutes(rg *gin.RouterGroup, h *handler.AntrianHandler) {
	antrianRoutes := rg.Group("/antrian")
	{
		antrianRoutes.GET("", h.GetAll)
		antrianRoutes.GET("/:id", h.GetByID)

		userAdmin := antrianRoutes.Group("")
		userAdmin.Use(middleware.Authorize("Administrasi"))
		{
			userAdmin.POST("", h.Create)
			userAdmin.DELETE("/:id", h.Delete)
		}

		userPoli := antrianRoutes.Group("")
		userPoli.Use(middleware.Authorize("Administrasi", "Poliklnik"))
		{
			userPoli.PUT("/:id", h.Update)
		}
	}
}
