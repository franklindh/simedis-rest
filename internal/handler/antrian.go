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

type AntrianHandler struct {
	Service *service.AntrianService
}

func NewAntrianHandler(svc *service.AntrianService) *AntrianHandler {
	return &AntrianHandler{Service: svc}
}

func (h *AntrianHandler) Create(c *gin.Context) {
	var req model.CreateAntrianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	createdAntrian, err := h.Service.CreateAntrian(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrScheduleOverlap) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrAntrianExists) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrForeignKey) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdAntrian, "data created successfully")
}

func (h *AntrianHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	poliID, _ := strconv.Atoi(c.Query("poli_id"))

	params := repository.ParamsGetAllAntrian{
		Page:          page,
		PageSize:      pageSize,
		SortBy:        c.DefaultQuery("sort", "created_at_asc"),
		StatusFilter:  c.Query("status"),
		TanggalFilter: c.Query("tanggal"),
		PoliIDFilter:  poliID,
	}

	responseData, metadata, err := h.Service.GetAllAntrianDetails(c.Request.Context(), params)
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

func (h *AntrianHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	antrian, err := h.Service.GetAntrianByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, antrian, "success")
}

func (h *AntrianHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var req model.UpdateAntrianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	result, err := h.Service.UpdateAntrian(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *AntrianHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	err = h.Service.DeleteAntrian(c.Request.Context(), id)
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
