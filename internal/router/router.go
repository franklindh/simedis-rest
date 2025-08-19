package router

import (
	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-gonic/gin"
)

func New(app *config.Application) *gin.Engine {
	router := gin.Default()

	db := app.DB
	cfg := app.Config

	poliRepo := repository.NewPoliRepository(db)
	poliService := service.NewPoliService(poliRepo)
	poliHandler := handler.NewPoliHandler(poliService)

	petugasRepo := repository.NewPetugasRepository(db)
	petugasService := service.NewPetugasService(petugasRepo, cfg)
	petugasHandler := handler.NewPetugasHandler(petugasService)

	jadwalRepo := repository.NewJadwalRepository(db)
	jadwalService := service.NewJadwalService(jadwalRepo)
	jadwalHandler := handler.NewJadwalHandler(jadwalService)

	pasienRepo := repository.NewPasienRepository(db)
	pasienService := service.NewPasienService(pasienRepo)
	pasienHandler := handler.NewPasienHandler(pasienService)

	antrianRepo := repository.NewAntrianRepository(db)
	antrianService := service.NewAntrianService(antrianRepo, jadwalRepo)
	antrianHandler := handler.NewAntrianHandler(antrianService)

	icdRepo := repository.NewIcdRepository(db)
	icdService := service.NewIcdService(icdRepo)
	icdHandler := handler.NewIcdHandler(icdService)

	pemeriksaanRepo := repository.NewPemeriksaanRepository(db)
	pemeriksaanService := service.NewPemeriksaanService(pemeriksaanRepo, antrianRepo)
	pemeriksaanHandler := handler.NewPemeriksaanHandler(pemeriksaanService)

	laporanRepo := repository.NewLaporanRepository(db)
	laporanService := service.NewLaporanService(laporanRepo)
	laporanHandler := handler.NewLaporanHandler(laporanService)

	jenisPemeriksaanLabRepo := repository.NewJenisPemeriksaanLabRepository(db)
	jenisPemeriksaanLabService := service.NewJenisPemeriksaanLabService(jenisPemeriksaanLabRepo)
	jenisPemeriksaanLabHandler := handler.NewJenisPemeriksaanLabHandler(jenisPemeriksaanLabService)

	pemeriksaanLabRepo := repository.NewPemeriksaanLabRepository(db)
	pemeriksaanLabService := service.NewPemeriksaanLabService(pemeriksaanLabRepo)
	pemeriksaanLabHandler := handler.NewPemeriksaanLabHandler(pemeriksaanLabService)

	// public
	router.POST("/login/petugas", petugasHandler.Login)

	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware(cfg))
	{
		PoliRoutes(authRoutes, poliHandler)
		PetugasRoutes(authRoutes, petugasHandler)
		JadwalRoutes(authRoutes, jadwalHandler)
		PasienRoutes(authRoutes, pasienHandler)
		AntrianRoutes(authRoutes, antrianHandler)
		IcdRoutes(authRoutes, icdHandler)
		PemeriksaanRoutes(authRoutes, pemeriksaanHandler)
		LaporanRoutes(authRoutes, laporanHandler)
		JenisPemeriksaanLabRoutes(authRoutes, jenisPemeriksaanLabHandler)
		PemeriksaanLabRoutes(authRoutes, pemeriksaanLabHandler)
	}

	return router
}
