package handler

import (
	"net/http"

	"strconv"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type PasienHandler struct {
	Repo *repository.PasienRepository
}

func NewPasienHandler(repo *repository.PasienRepository) *PasienHandler {
	return &PasienHandler{Repo: repo}
}

func (h *PasienHandler) Create(c *gin.Context) {
	var newPasien model.Pasien
	if err := c.ShouldBindJSON(&newPasien); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	createdPasien, err := h.Repo.Create(newPasien)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "data already exists", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdPasien, "data created successfully")
}

func (h *PasienHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sort := c.DefaultQuery("sort", "nama_pasien_asc")
	nameFilter := c.Query("name")
	nikFilter := c.Query("nik")
	noRekamMedis := c.Query("no_rekam_medis")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	params := repository.ParamsGetAllPasien{
		NameFilter:   nameFilter,
		NIKFilter:    nikFilter,
		NoRekamMedis: noRekamMedis,
		SortBy:       sort,
		Page:         page,
		PageSize:     pageSize,
	}

	allPasien, metadata, err := h.Repo.GetAll(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"metadata": metadata,
		"data":     allPasien,
	})
}

func (h *PasienHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	pasien, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, pasien, "success")
}

func (h *PasienHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var updatedPasien model.Pasien
	if err := c.ShouldBindJSON(&updatedPasien); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	result, err := h.Repo.Update(id, updatedPasien)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "data already exists", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *PasienHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "data deleted successfully")
}
