package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func JenisPemeriksaanLabRoutes(rg *gin.RouterGroup, h *handler.JenisPemeriksaanLabHandler) {
	jenisPemeriksaanLab := rg.Group("/jenis-pemeriksaan-lab")
	{
		jenisPemeriksaanLab.GET("", h.GetAll)
		jenisPemeriksaanLab.GET("/:id", h.GetByID)

		user := jenisPemeriksaanLab.Group("")
		user.Use(middleware.Authorize("Administrasi", "Lab"))
		{
			user.POST("", h.Create)
			user.PUT("/:id", h.Update)
			user.DELETE("/:id", h.Delete)
		}
	}
}
