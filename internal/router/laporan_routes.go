package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func LaporanRoutes(rg *gin.RouterGroup, h *handler.LaporanHandler) {
	laporanRoutes := rg.Group("/laporan")
	{
		laporanRoutes.GET("/kunjungan-poli", h.GetKunjunganPoli)
		laporanRoutes.GET("/penyakit-teratas", h.GetPenyakitTeratas)
	}
}
