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

func (h *PetugasHandler) Create(c *gin.Context) {
	var newPetugas model.Petugas

	if err := c.ShouldBindJSON(&newPetugas); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	if err := utils.ValidatePetugasUsername(newPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
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
			utils.ErrorResponse(c, http.StatusConflict, "data already exists", nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create data", err)
		return
	}

	createdPetugas.Password = ""
	utils.SuccessResponse(c, http.StatusCreated, createdPetugas, "data created successfully")
}

func (h *PetugasHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	sort := c.DefaultQuery("sort", "created_at_desc")
	nameFilter := c.Query("name")
	roleFilter := c.Query("role")
	statusFilter := c.Query("status")

	params := repository.ParamsGetAllPetugas{
		NameOrUsernameFilter: nameFilter,
		RoleFilter:           roleFilter,
		StatusFilter:         statusFilter,
		SortBy:               sort,
		Page:                 page,
		PageSize:             pageSize,
	}

	allPetugas, metadata, err := h.Repo.GetAll(params)
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

func (h *PetugasHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var updatedPetugas model.Petugas
	if err := c.ShouldBindJSON(&updatedPetugas); err != nil {
		errorMessage := utils.FormatValidationError(err)
		utils.ErrorResponse(c, http.StatusBadRequest, errorMessage, err)
		return
	}

	if err := utils.ValidatePetugasUsername(updatedPetugas); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	result, err := h.Repo.Update(id, updatedPetugas)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			utils.ErrorResponse(c, http.StatusConflict, "username already exists", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result, "data updated successfully")
}

func (h *PetugasHandler) Delete(c *gin.Context) {
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
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to delete data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "data deleted successfully")
}

func (h *PetugasHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "username and password are required", err)
		return
	}

	user, err := h.Repo.GetByUsername(loginData.Username)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid username or password", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "database error", err)
		return
	}

	err = utils.VerifyPassword(loginData.Password, user.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "invalid username or password", nil)
		return
	}

	jwtSecret := []byte(h.Config.JWTSecret)

	token, err := utils.SignToken(strconv.Itoa(user.ID), user.Username, user.Role, jwtSecret)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to generate token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "login successful",
		"token":   token,
	})
}
