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
	registrationRepo := repositories.NewRegistrationRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	placementRepo := repositories.NewClassPlacementRepository(db)

	// 4. Inisialisasi Services
	studentService := services.NewStudentService(studentRepo)
	courseService := services.NewCourseService(courseRepo)
	classService := services.NewClassService(classRepo, courseRepo)
	registrationService := services.NewRegistrationService(registrationRepo, studentRepo, courseRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	placementService := services.NewClassPlacementService(placementRepo, registrationRepo, classRepo)

	// 5. Inisialisasi Handlers
	studentHandler := handlers.NewStudentHandler(studentService)
	courseHandler := handlers.NewCourseHandler(courseService)
	classHandler := handlers.NewClassHandler(classService)
	registrationHandler := handlers.NewRegistrationHandler(registrationService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	placementHandler := handlers.NewClassPlacementHandler(placementService)

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
			students.GET("/:id/registrations", registrationHandler.GetByStudentID)
		}

		// --- Course Endpoints ---
		courses := apiV1.Group("/courses")
		{
			courses.POST("", courseHandler.Create)
			courses.GET("", courseHandler.GetAll)
			courses.GET("/:id", courseHandler.GetByID)
			courses.PUT("/:id", courseHandler.Update)
			courses.DELETE("/:id", courseHandler.Delete)
			courses.GET("/:id/registrations", registrationHandler.GetByCourseID)
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

		// --- Registration Endpoints ---
		registrations := apiV1.Group("/registrations")
		{
			registrations.POST("", registrationHandler.Register)
			registrations.GET("", registrationHandler.GetAll)
			registrations.GET("/:id", registrationHandler.GetByID)
			registrations.PUT("/:id/cancel", registrationHandler.CancelRegistration)
		}

		// --- Payment Endpoints ---
		payments := apiV1.Group("/payments")
		{
			payments.GET("", paymentHandler.GetAll)
			payments.GET("/:id", paymentHandler.GetByID)
			payments.POST("/:id/pay", paymentHandler.Pay)
		}

		// --- Class Placement Endpoints ---
		placements := apiV1.Group("/class-placements")
		{
			placements.POST("", placementHandler.PlaceStudent)
			placements.GET("", placementHandler.GetAll)
			placements.GET("/:id", placementHandler.GetByID)
		}
	}

	return r
}
