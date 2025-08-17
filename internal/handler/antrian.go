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

type AntrianHandler struct {
	Repo *repository.AntrianRepository
}

func NewAntrianHandler(repo *repository.AntrianRepository) *AntrianHandler {
	return &AntrianHandler{Repo: repo}
}

func (h *AntrianHandler) Create(c *gin.Context) {
	var input model.AntrianCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	newAntrian := model.Antrian{
		JadwalID:  input.JadwalID,
		PasienID:  input.PasienID,
		Prioritas: input.Prioritas,
		Status:    input.Status,
	}

	createdAntrian, err := h.Repo.Create(newAntrian)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			utils.ErrorResponse(c, http.StatusBadRequest, "invalid jadwal id or pasien id", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, createdAntrian, "data created successfully")
}

func (h *AntrianHandler) GetAll(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sort := c.DefaultQuery("sort", "created_at_asc")
	statusFilter := c.Query("status")
	tanggalFilter := c.Query("tanggal")
	poliID, _ := strconv.Atoi(c.Query("poli_id"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	params := repository.ParamsGetAllAntrian{
		StatusFilter:  statusFilter,
		TanggalFilter: tanggalFilter,
		PoliIDFilter:  poliID,
		SortBy:        sort,
		Page:          page,
		PageSize:      pageSize,
	}

	allAntrian, metadata, err := h.Repo.GetAll(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	var responseData []model.AntrianDetail
	for _, antrian := range allAntrian {
		detail := model.AntrianDetail{
			ID:           antrian.ID,
			NomorAntrian: antrian.NomorAntrian,
			Prioritas:    antrian.Prioritas,
			Status:       antrian.Status,
			Jadwal: struct {
				ID      int    `json:"id"`
				Tanggal string `json:"tanggal"`
				Poli    struct {
					Name string `json:"name"`
				} `json:"poli"`
				Dokter struct {
					Name string `json:"name"`
				} `json:"dokter"`
			}{
				ID:      antrian.Jadwal.ID,
				Tanggal: antrian.Jadwal.Tanggal,
				Poli: struct {
					Name string `json:"name"`
				}{Name: antrian.Jadwal.Poli.Nama},
				Dokter: struct {
					Name string `json:"name"`
				}{Name: antrian.Jadwal.Petugas.Name},
			},
			Pasien: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   antrian.Pasien.ID,
				Name: antrian.Pasien.NamaPasien,
			},
		}
		responseData = append(responseData, detail)
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
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	antrian, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
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
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	var updatedAntrian model.Antrian
	if err := c.ShouldBindJSON(&updatedAntrian); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	result, err := h.Repo.Update(id, updatedAntrian)
	if err != nil {
		if err == repository.ErrNotFound {
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
