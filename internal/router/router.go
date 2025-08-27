package router

import (
	"time"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/handler"
	"github.com/franklindh/simedis-api/internal/middleware"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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

	router.Use(secure.New(secure.Config{
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
		// ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy: "no-referrer",
		IsDevelopment:  false,
	}))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://icikiwir.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	rate := limiter.Rate{Period: 1 * time.Hour, Limit: 100}
	store := memory.NewStore()
	limitermiddleware := ginmiddleware.NewMiddleware(limiter.New(store, rate))
	router.Use(limitermiddleware)

	router.Use(gzip.Gzip(gzip.DefaultCompression))

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
