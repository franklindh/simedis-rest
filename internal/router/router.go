package router

import (
	"net/http"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/gin-gonic/gin"
)

func New(app *config.Application) *gin.Engine {
	router := gin.Default()

	db := app.DB
	cfg := app.Config

	poliRepo := repository.NewPoliRepository(db)
	poliHandler := handler.NewPoliHandler(poliRepo)

	petugasRepo := repository.NewPetugasRepository(db)
	petugasHandler := handler.NewPetugasHandler(petugasRepo, cfg)

	router.GET("/poli", poliHandler.GetAll)
	router.GET("/poli/:id", poliHandler.GetByID)
	router.POST("/poli", poliHandler.Create)
	router.PUT("/poli/:id", poliHandler.Update)
	router.DELETE("/poli/:id", poliHandler.Delete)

	router.GET("/petugas", petugasHandler.GetAll)
	router.GET("/petugas/:id", petugasHandler.GetByID)
	router.POST("/petugas", petugasHandler.Create)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return router
}
