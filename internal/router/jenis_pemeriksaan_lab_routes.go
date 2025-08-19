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
		jenisPemeriksaanLab.POST("", h.Create)
		jenisPemeriksaanLab.PUT("/:id", h.Update)
		jenisPemeriksaanLab.DELETE("/:id", h.Delete)

		user := jenisPemeriksaanLab.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{

		}
	}
}
