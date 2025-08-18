package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func IcdRoutes(rg *gin.RouterGroup, h *handler.IcdHandler) {
	icdRoutes := rg.Group("/icd")
	{
		icdRoutes.GET("", h.GetAll)
		icdRoutes.GET("/:id", h.GetByID)

		user := icdRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
