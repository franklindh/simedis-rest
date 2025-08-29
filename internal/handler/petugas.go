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

type PetugasHandler struct {
	Service *service.PetugasService
}

func NewPetugasHandler(svc *service.PetugasService) *PetugasHandler {
	return &PetugasHandler{Service: svc}
}

func (h *PetugasHandler) Create(c *gin.Context) {
	var newPetugas model.CreatePetugasRequest
	if err := c.ShouldBindJSON(&newPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	if err := utils.ValidatePetugasUsername(newPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	createdPetugas, err := h.Service.CreatePetugas(c.Request.Context(), newPetugas)
	if err != nil {
		if errors.Is(err, service.ErrPetugasConflict) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdPetugas, "data created successfully")
}

func (h *PetugasHandler) GetAll(c *gin.Context) {
	var params repository.ParamsGetAllPetugas

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

	allPetugas, metadata, err := h.Service.GetAllPetugas(c.Request.Context(), params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"metadata": metadata,
		"data":     allPetugas,
	})
}

func (h *PetugasHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	petugas, err := h.Service.GetPetugasByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, petugas, "success")
}

func (h *PetugasHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var updatedPetugas model.UpdatePetugasRequest
	if err := c.ShouldBindJSON(&updatedPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	result, err := h.Service.UpdatePetugas(c.Request.Context(), id, updatedPetugas)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
		case errors.Is(err, service.ErrPetugasConflict):
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *PetugasHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	err = h.Service.DeletePetugas(c.Request.Context(), id)
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

func (h *PetugasHandler) Login(c *gin.Context) {
	var req model.LoginPetugasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "username and password are required", err)
		return
	}

	token, err := h.Service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "an internal error occurred", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
		"token":   token,
	})
}

func (h *PetugasHandler) ChangePassword(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User ID not found in token", nil)
		return
	}

	petugasID, ok := userID.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid user ID format in token", nil)
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	err := h.Service.ChangePassword(c.Request.Context(), petugasID, req)
	if err != nil {

		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Password changed successfully")
}
