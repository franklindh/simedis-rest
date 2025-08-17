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

type IcdHandler struct {
	Repo *repository.IcdRepository
}

func NewIcdHandler(repo *repository.IcdRepository) *IcdHandler {
	return &IcdHandler{Repo: repo}
}

func (h *IcdHandler) Create(c *gin.Context) {
	var newIcd model.Icd
	if err := c.ShouldBindJSON(&newIcd); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	createdIcd, err := h.Repo.Create(newIcd)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "data already exists", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdIcd, "data created successfully")
}

func (h *IcdHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	sort := c.DefaultQuery("sort", "kode_asc")
	kodeFilter := c.Query("kode")
	namaFilter := c.Query("nama")
	statusFilter := c.Query("status")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	params := repository.ParamsGetAllIcd{
		KodeFilter:   kodeFilter,
		NamaFilter:   namaFilter,
		StatusFilter: statusFilter,
		SortBy:       sort,
		Page:         page,
		PageSize:     pageSize,
	}

	allIcd, metadata, err := h.Repo.GetAll(params)
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
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	icd, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
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
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	var updatedIcd model.Icd
	if err := c.ShouldBindJSON(&updatedIcd); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	result, err := h.Repo.Update(id, updatedIcd)
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

func (h *IcdHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "data deleted successfully")
}
