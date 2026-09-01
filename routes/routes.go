package routes

import (
	"net/http"

	"english-course-api/config"
	"english-course-api/middleware"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter menginisialisasi Gin Engine, middleware, dan mendaftarkan route endpoint
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Atur mode Gin (Release jika production, Debug jika development)
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Gunakan gin.New() agar kita bisa mengontrol middleware secara manual
	r := gin.New()

	// 1. Pasang Middleware Global
	r.Use(gin.Recovery())            // Menangkap panic agar server tidak crash
	r.Use(middleware.HTTPLogger())    // Structured Logging dengan Zerolog

	// 2. Health Check Endpoint (Root & API v1)
	r.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, "Service is healthy 🚀", gin.H{
			"app_name": cfg.App.Name,
			"env":      cfg.App.Env,
			"status":   "UP",
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			utils.SuccessResponse(c, http.StatusOK, "API v1 is healthy 🚀", gin.H{
				"version": "v1",
				"status":  "UP",
			})
		})

		// Route modul-modul lain (Student, Course, Registration, dll.) akan didaftarkan di sini pada fase berikutnya
	}

	return r
}
