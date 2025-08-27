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

type IcdHandler struct {
	Service *service.IcdService
}

func NewIcdHandler(svc *service.IcdService) *IcdHandler {
	return &IcdHandler{Service: svc}
}

func (h *IcdHandler) Create(c *gin.Context) {
	var req model.CreateIcdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	createdIcd, err := h.Service.CreateIcd(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrIcdConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdIcd, "data created successfully")
}

func (h *IcdHandler) GetAll(c *gin.Context) {
	var params repository.ParamsGetAllIcd

	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err)
		return
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 5
	}
	if params.SortBy == "" {
		params.SortBy = "created_at_desc"
	}

	allIcd, metadata, err := h.Service.GetAllIcd(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"metadata": metadata,
		"data":     allIcd,
	})
}

func (h *IcdHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	icd, err := h.Service.GetIcdByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, icd, "success")
}

func (h *IcdHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var req model.UpdateIcdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	result, err := h.Service.UpdateIcd(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		if errors.Is(err, service.ErrIcdConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *IcdHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	err = h.Service.DeleteIcd(c.Request.Context(), id)
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
