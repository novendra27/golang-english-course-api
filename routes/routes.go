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

	// 3. Inisialisasi Repositories
	studentRepo := repositories.NewStudentRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	classRepo := repositories.NewClassRepository(db)

	// 4. Inisialisasi Services
	studentService := services.NewStudentService(studentRepo)
	courseService := services.NewCourseService(courseRepo)
	classService := services.NewClassService(classRepo, courseRepo)

	// 5. Inisialisasi Handlers
	studentHandler := handlers.NewStudentHandler(studentService)
	courseHandler := handlers.NewCourseHandler(courseService)
	classHandler := handlers.NewClassHandler(classService)

	// 6. API v1 Router Group
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

		// --- Class Endpoints ---
		classes := apiV1.Group("/classes")
		{
			classes.POST("", classHandler.Create)
			classes.GET("", classHandler.GetAll)
			classes.GET("/:id", classHandler.GetByID)
			classes.PUT("/:id", classHandler.Update)
			classes.DELETE("/:id", classHandler.Delete)
			classes.GET("/:id/students", classHandler.GetStudents)
		}
	}

	return r
}
