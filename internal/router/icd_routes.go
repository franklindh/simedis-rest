package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func IcdRoutes(rg *gin.RouterGroup, h *handler.IcdHandler) {
	icdRoutes := rg.Group("/icd")
	{
		icdRoutes.GET("", h.GetAll)
		icdRoutes.POST("", h.Create)
		icdRoutes.GET("/:id", h.GetByID)
		icdRoutes.PUT("/:id", h.Update)
		icdRoutes.DELETE("/:id", h.Delete)
	}
}
