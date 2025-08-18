package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func LaporanRoutes(rg *gin.RouterGroup, h *handler.LaporanHandler) {
	laporanRoutes := rg.Group("/laporan")
	{
		user := laporanRoutes.Group("")
		user.Use(middleware.Authorize("Administrasi"))
		{
			user.GET("/kunjungan-poli", h.GetKunjunganPoli)
			user.GET("/penyakit-teratas", h.GetPenyakitTeratas)
		}
	}
}
