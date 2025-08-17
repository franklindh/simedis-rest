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

type JadwalHandler struct {
	Repo *repository.JadwalRepository
}

func NewJadwalHandler(repo *repository.JadwalRepository) *JadwalHandler {
	return &JadwalHandler{Repo: repo}
}

func (h *JadwalHandler) Create(c *gin.Context) {
	var newJadwal model.Jadwal

	if err := c.ShouldBindJSON(&newJadwal); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	createdJadwal, err := h.Repo.Create(newJadwal)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "data already exists", nil)
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
	sort := c.DefaultQuery("sort", "tanggal_desc")

	poliID, _ := strconv.Atoi(c.Query("poli_id"))
	petugasID, _ := strconv.Atoi(c.Query("petugas_id"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}

	params := repository.ParamsGetAllJadwal{
		PoliIDFilter:    poliID,
		PetugasIDFilter: petugasID,
		StartDateFilter: startDate,
		EndDateFilter:   endDate,
		SortBy:          sort,
		Page:            page,
		PageSize:        pageSize,
	}

	allJadwal, metadata, err := h.Repo.GetAll(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	var responseData []model.JadwalDetail
	for _, jadwal := range allJadwal {
		detail := model.JadwalDetail{
			ID:           jadwal.ID,
			Tanggal:      jadwal.Tanggal,
			WaktuMulai:   jadwal.WaktuMulai,
			WaktuSelesai: jadwal.WaktuSelesai,
			Keterangan:   jadwal.Keterangan,
			Petugas: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   jadwal.Petugas.ID,
				Name: jadwal.Petugas.Nama,
			},
			Poli: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   jadwal.Poli.ID,
				Name: jadwal.Poli.Nama,
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

func (h *JadwalHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ID format", err)
		return
	}

	jadwal, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
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

	var updatedJadwal model.Jadwal
	if err := c.ShouldBindJSON(&updatedJadwal); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	result, err := h.Repo.Update(id, updatedJadwal)
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

func (h *JadwalHandler) Delete(c *gin.Context) {
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
