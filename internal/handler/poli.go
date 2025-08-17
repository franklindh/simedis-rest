package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository" // For ErrNotFound
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-gonic/gin"
)

type PoliHandler struct {
	Service *service.PoliService
}

func NewPoliHandler(svc *service.PoliService) *PoliHandler {
	return &PoliHandler{Service: svc}
}

func (h *PoliHandler) Create(c *gin.Context) {
	var newPoli model.Poli
	if err := c.ShouldBindJSON(&newPoli); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	createdPoli, err := h.Service.CreateOrRestorePoli(c.Request.Context(), newPoli)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPoliConflict):
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrPoliRestored):
			// This is a "success" case communicated via an error.
			utils.SuccessResponse(c, http.StatusOK, createdPoli, err.Error())
		default:
			// For any other error (DB error, etc.)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create data", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdPoli, "Data created successfully")
}

func (h *PoliHandler) GetAll(c *gin.Context) {
	polis, err := h.Service.GetAllPolis(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(polis),
		"data":   polis,
	})
}

func (h *PoliHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	poli, err := h.Service.GetPoliByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, poli, "Success")
}

func (h *PoliHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var updatedPoli model.Poli
	if err := c.ShouldBindJSON(&updatedPoli); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	result, err := h.Service.UpdatePoli(c.Request.Context(), id, updatedPoli)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			utils.ErrorResponse(c, http.StatusNotFound, "Data not found", nil)
		case errors.Is(err, service.ErrPoliConflict):
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update data", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "Data updated successfully")
}

func (h *PoliHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	err = h.Service.DeletePoli(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Data deleted successfully")
}
