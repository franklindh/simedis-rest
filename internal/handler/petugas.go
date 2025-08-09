package handler

import (
	"net/http"
	"strconv"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type PetugasHandler struct {
	Repo   *repository.PetugasRepository
	Config *config.Config
}

func NewPetugasHandler(repo *repository.PetugasRepository, cfg *config.Config) *PetugasHandler {
	return &PetugasHandler{Repo: repo, Config: cfg}
}

func (h *PetugasHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	sort := c.DefaultQuery("sort", "created_at_desc")
	nameFilter := c.Query("name")

	params := repository.ParamsGetAllParams{
		NameFilter: nameFilter,
		SortBy:     sort,
		Page:       page,
		PageSize:   pageSize,
	}

	allPetugas, metadata, err := h.Repo.GetAll(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"metadata": metadata,
		"data":     allPetugas,
	})
}

func (h *PetugasHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	petugas, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, petugas, "success")
}

func (h *PetugasHandler) Create(c *gin.Context) {
	var newPetugas model.Petugas

	if err := c.ShouldBindJSON(&newPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if newPetugas.Username == "" || newPetugas.Name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "username and name are required fields", nil)
		return
	}

	defaultPassword := h.Config.DefaultPetugasPassword

	encodedHash, err := utils.HashPassword(defaultPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to hash password", err)
		return
	}
	newPetugas.Password = encodedHash

	createdPetugas, err := h.Repo.Create(newPetugas)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "username already exists", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	createdPetugas.Password = ""
	utils.SuccessResponse(c, http.StatusCreated, createdPetugas, "data created successfully")
}
