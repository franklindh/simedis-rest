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

type JadwalHandler struct {
	Service *service.JadwalService
}

func NewJadwalHandler(svc *service.JadwalService) *JadwalHandler {
	return &JadwalHandler{Service: svc}
}

func (h *JadwalHandler) Create(c *gin.Context) {
	var req model.JadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	createdJadwal, err := h.Service.CreateJadwal(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrWaktuSelesaiInvalid) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrJadwalConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdJadwal, "data created successfully")
}

func (h *JadwalHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	poliID, _ := strconv.Atoi(c.Query("poli_id"))
	petugasID, _ := strconv.Atoi(c.Query("petugas_id"))

	params := repository.ParamsGetAllJadwal{
		Page:            page,
		PageSize:        pageSize,
		SortBy:          c.DefaultQuery("sort", "tanggal_desc"),
		PoliIDFilter:    poliID,
		PetugasIDFilter: petugasID,
		StartDateFilter: c.Query("start_date"),
		EndDateFilter:   c.Query("end_date"),
	}

	responseData, metadata, err := h.Service.GetAllJadwal(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"metadata": metadata,
		"data":     responseData,
	})
}

func (h *JadwalHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	jadwal, err := h.Service.GetJadwalByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, jadwal, "success")
}

func (h *JadwalHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	var req model.JadwalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	result, err := h.Service.UpdateJadwal(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		if errors.Is(err, service.ErrJadwalConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *JadwalHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	err = h.Service.DeleteJadwal(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "data deleted successfully")
}
