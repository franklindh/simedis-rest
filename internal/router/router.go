package router

import (
	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/gin-gonic/gin"
)

func New(app *config.Application) *gin.Engine {
	router := gin.Default()

	db := app.DB
	cfg := app.Config

	poliRepo := repository.NewPoliRepository(db)
	poliHandler := handler.NewPoliHandler(poliRepo)

	// petugasRepo := repository.NewPetugasRepository(db)
	// petugasHandler := handler.NewPetugasHandler(petugasRepo, cfg)

	// jadwalRepo := repository.NewJadwalRepository(db)
	// jadwalHandler := handler.NewJadwalHandler(jadwalRepo)

	// publc
	// router.POST("/login/petugas", petugasHandler.Login)
	// router.POST("/petugas", petugasHandler.Create)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(cfg))
	{

		// PetugasRoutes(authRoutes, petugasHandler)
		// JadwalRoutes(authRoutes, jadwalHandler)
	}
	PoliRoutes(nil, poliHandler)

	return router
}
