package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PoliHandler struct {
	Repo *repository.PoliRepository
}

func NewPoliHandler(repo *repository.PoliRepository) *PoliHandler {
	return &PoliHandler{Repo: repo}
}

func (h *PoliHandler) GetAll(c *gin.Context) {
	poli, err := h.Repo.GetAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "failed to retrieve data", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(poli),
		"data":   poli,
	})
}

func (h *PoliHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	poli, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repository.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "data not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to retrieve data", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, poli, "success")
}

func (h *PoliHandler) Create(c *gin.Context) {
	var newPoli model.Poli
	if err := c.ShouldBindJSON(&newPoli); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	existingPoli, err := h.Repo.FindByNameIncludingDeleted(newPoli.Name)

	if err != nil && err != sql.ErrNoRows {
		utils.ErrorResponse(c, http.StatusInternalServerError, "database error", err)
		return
	}

	if err == nil {
		if existingPoli.DeletedAt.Valid {
			restoredPoli, restoreErr := h.Repo.Restore(existingPoli.ID)
			if restoreErr != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, "failed to restore data", restoreErr)
				return
			}
			utils.SuccessResponse(c, http.StatusOK, restoredPoli, "data restored successfully")
			return
		} else {
			utils.ErrorResponse(c, http.StatusConflict, "data with that name already exists", nil)
			return
		}
	}

	if err == sql.ErrNoRows {
		createdPoli, createErr := h.Repo.Create(newPoli)
		if createErr != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "failed to save data to database", createErr)
			return
		}
		utils.SuccessResponse(c, http.StatusCreated, createdPoli, "data created successfully")
		return
	}
}

// func (h *PoliHandler) Update(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
// 		return
// 	}

// 	var updates map[string]any
// 	if err := c.ShouldBindJSON(&updates); err != nil {
// 		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
// 		return
// 	}

// 	updatedPoli, err := h.Repo.Update(id, updates)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			utils.ErrorResponse(c, http.StatusNotFound, "data not found", err)
// 			return
// 		}
// 		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
// 			utils.ErrorResponse(c, http.StatusConflict, "Polyclinic name already exists", nil)
// 			return
// 		}
// 		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to update data", err)
// 		return
// 	}

// 	utils.SuccessResponse(c, http.StatusOK, updatedPoli, "data updated successfully")
// }

func (h *PoliHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id format", err)
		return
	}

	var updatedPoli model.Poli
	if err := c.ShouldBindJSON(&updatedPoli); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body", err)
		return
	}

	result, err := h.Repo.Update(id, updatedPoli)
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

func (h *PoliHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	fmt.Println(idStr)
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

	utils.SuccessResponse(c, http.StatusOK, nil, "data deleted successfulyy")
}
