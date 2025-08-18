package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PasienRoutes(rg *gin.RouterGroup, h *handler.PasienHandler) {
	pasienRoutes := rg.Group("/pasien")
	{
		pasienRoutes.GET("", h.GetAll)
		pasienRoutes.GET("/:id", h.GetByID)

		user := pasienRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
