package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-gonic/gin"
)

type PemeriksaanLabHandler struct {
	Service *service.PemeriksaanLabService
}

func NewPemeriksaanLabHandler(svc *service.PemeriksaanLabService) *PemeriksaanLabHandler {
	return &PemeriksaanLabHandler{Service: svc}
}

func (h *PemeriksaanLabHandler) Create(c *gin.Context) {
	pemeriksaanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	var reqs []model.CreateHasilLabRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	created, err := h.Service.CreateBatch(c.Request.Context(), pemeriksaanID, reqs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, created, "Lab results created successfully")
}

func (h *PemeriksaanLabHandler) GetAll(c *gin.Context) {
	pemeriksaanID, err := strconv.Atoi(c.Param("id"))
	fmt.Println(pemeriksaanID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid id format", err)
		return
	}

	results, err := h.Service.GetAllByPemeriksaanID(c.Request.Context(), pemeriksaanID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, results, "Success")
}

func (h *PemeriksaanLabHandler) Update(c *gin.Context) {
	hasilID, err := strconv.Atoi(c.Param("hasil_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid hasil_id format", err)
		return
	}

	var req model.UpdateHasilLabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.FormatValidationError(err), err)
		return
	}

	updated, err := h.Service.Update(c.Request.Context(), hasilID, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update data", err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, updated, "Data updated successfully")
}

func (h *PemeriksaanLabHandler) Delete(c *gin.Context) {
	hasilID, err := strconv.Atoi(c.Param("hasil_id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid hasil_id format", err)
		return
	}

	err = h.Service.Delete(c.Request.Context(), hasilID)
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
