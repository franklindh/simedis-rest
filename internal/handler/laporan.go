package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-gonic/gin"
)

type LaporanHandler struct {
	Service *service.LaporanService
}

func NewLaporanHandler(svc *service.LaporanService) *LaporanHandler {
	return &LaporanHandler{Service: svc}
}

func (h *LaporanHandler) GetKunjunganPoli(c *gin.Context) {

	today := time.Now()
	firstDayOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	startDate := c.DefaultQuery("startDate", firstDayOfMonth.Format("2006-01-02"))
	endDate := c.DefaultQuery("endDate", today.Format("2006-01-02"))

	laporan, err := h.Service.GetLaporanKunjunganPerPoli(c.Request.Context(), startDate, endDate)
	if err != nil {

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, laporan, "Laporan kunjungan per poli berhasil diambil")
}

func (h *LaporanHandler) GetPenyakitTeratas(c *gin.Context) {
	today := time.Now()
	firstDayOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	limitConv, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	startDate := c.DefaultQuery("startDate", firstDayOfMonth.Format("2006-01-02"))
	endDate := c.DefaultQuery("endDate", today.Format("2006-01-02"))
	limit := limitConv

	laporan, err := h.Service.GetLaporanPenyakitTeratas(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, laporan, "Laporan penyakit teratas berhasil diambil")
}
