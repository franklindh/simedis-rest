package router

import (
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PemeriksaanLabRoutes(rg *gin.RouterGroup, h *handler.PemeriksaanLabHandler) {

	hasilLabGroup := rg.Group("/pemeriksaan/:id/hasil-lab")
	{

		hasilLabGroup.GET("", h.GetAll, middleware.Authorize("Dokter", "Lab", "Poliklinik"))
		hasilLabGroup.POST("", h.Create, middleware.Authorize("Dokter", "Lab", "Poliklinik"))
	}

	rg.PUT("/hasil-lab/:hasil_id", h.Update, middleware.Authorize("Dokter", "Lab", "Poliklinik"))
	rg.DELETE("/hasil-lab/:hasil_id", h.Delete, middleware.Authorize("Dokter", "Lab", "Poliklinik"))
}
