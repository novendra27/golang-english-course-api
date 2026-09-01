package routes

import (
	"net/http"

	"english-course-api/config"
	"english-course-api/handlers"
	"english-course-api/middleware"
	"english-course-api/repositories"
	"english-course-api/services"
	"english-course-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter menginisialisasi Gin Engine, middleware, dependency injection, dan mendaftarkan route endpoint
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Atur mode Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	// 1. Middleware Global
	r.Use(gin.Recovery())
	r.Use(middleware.HTTPLogger())

	// 2. Health Check Endpoint
	r.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, "Service is healthy 🚀", gin.H{
			"app_name": cfg.App.Name,
			"env":      cfg.App.Env,
			"status":   "UP",
		})
	})

	// 3. Inisialisasi Repositories, Services, dan Handlers (Dependency Injection)
	studentRepo := repositories.NewStudentRepository(db)
	studentService := services.NewStudentService(studentRepo)
	studentHandler := handlers.NewStudentHandler(studentService)

	courseRepo := repositories.NewCourseRepository(db)
	courseService := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	// 4. API v1 Router Group
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			utils.SuccessResponse(c, http.StatusOK, "API v1 is healthy 🚀", gin.H{
				"version": "v1",
				"status":  "UP",
			})
		})

		// --- Student Endpoints ---
		students := apiV1.Group("/students")
		{
			students.POST("", studentHandler.Create)
			students.GET("", studentHandler.GetAll)
			students.GET("/:id", studentHandler.GetByID)
			students.PUT("/:id", studentHandler.Update)
			students.DELETE("/:id", studentHandler.Delete)
		}

		// --- Course Endpoints ---
		courses := apiV1.Group("/courses")
		{
			courses.POST("", courseHandler.Create)
			courses.GET("", courseHandler.GetAll)
			courses.GET("/:id", courseHandler.GetByID)
			courses.PUT("/:id", courseHandler.Update)
			courses.DELETE("/:id", courseHandler.Delete)
		}
	}

	return r
}
