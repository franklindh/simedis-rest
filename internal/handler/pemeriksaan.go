package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-gonic/gin"
)

type PemeriksaanHandler struct {
	Service *service.PemeriksaanService
}

func NewPemeriksaanHandler(svc *service.PemeriksaanService) *PemeriksaanHandler {
	return &PemeriksaanHandler{Service: svc}
}

func (h *PemeriksaanHandler) Create(c *gin.Context) {
	var req model.CreatePemeriksaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	created, err := h.Service.CreatePemeriksaan(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrPemeriksaanExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, created, "data created successfully")
}

func (h *PemeriksaanHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	pemeriksaan, err := h.Service.GetPemeriksaanByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, pemeriksaan, "success")
}

func (h *PemeriksaanHandler) GetRiwayatByPasienID(c *gin.Context) {
	pasienID, err := strconv.Atoi(c.Param("pasien_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid pasien_id format", err)
		return
	}

	riwayat, err := h.Service.GetRiwayatPemeriksaanPasien(c.Request.Context(), pasienID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   riwayat,
	})
}

func (h *PemeriksaanHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var req model.UpdatePemeriksaanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	updated, err := h.Service.UpdatePemeriksaan(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updated, "data updated successfully")
}
