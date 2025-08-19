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

type JenisPemeriksaanLabHandler struct {
	Service *service.JenisPemeriksaanLabService
}

func NewJenisPemeriksaanLabHandler(svc *service.JenisPemeriksaanLabService) *JenisPemeriksaanLabHandler {
	return &JenisPemeriksaanLabHandler{Service: svc}
}

func (h *JenisPemeriksaanLabHandler) Create(c *gin.Context) {
	var req model.CreateJenisPemeriksaanLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	created, err := h.Service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrJenisPemeriksaanConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, created, "data created successfully")
}

func (h *JenisPemeriksaanLabHandler) GetAll(c *gin.Context) {
	allData, err := h.Service.GetAll(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, allData, "success")
}

func (h *JenisPemeriksaanLabHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.Service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, data, "success")
}

func (h *JenisPemeriksaanLabHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req model.UpdateJenisPemeriksaanLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	updated, err := h.Service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		if errors.Is(err, service.ErrJenisPemeriksaanConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updated, "data updated successfully")
}

func (h *JenisPemeriksaanLabHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.Delete(c.Request.Context(), id)
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
