package main

import (
	"fmt"
	"os"

	"english-course-api/config"
	"english-course-api/routes"
	"english-course-api/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title           English Course Registration API
// @version         1.0
// @description     RESTful API backend for English course registration with clean layered architecture and dynamic i18n support.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Novendra
// @contact.email  support@englishcourse.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// 1. Setup Zerolog (Pretty output for console)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 2. Load Configuration from .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load application configuration")
	}

	// Adjust log level from config
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Info().
		Str("app_name", cfg.App.Name).
		Str("env", cfg.App.Env).
		Str("port", cfg.App.Port).
		Msg("Starting application initialization 🚀")

	// 3. Initialize i18n Bundle
	utils.InitI18n()

	// 4. Initialize PostgreSQL Database Connection
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Database connection warning (ensure PostgreSQL service is running)")
	} else {
		defer config.CloseDB(db)

		// 5. Auto-Migration for 6 domain models
		if err := config.AutoMigrate(db); err != nil {
			log.Fatal().Err(err).Msg("Failed to auto-migrate database schema")
		}
	}

	// 6. Initialize Gin Router & Middleware
	router := routes.SetupRouter(db, cfg)

	// 7. Run HTTP Server (Blocking Listener)
	serverAddr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info().Msgf("HTTP Server is active and listening on http://localhost%s 🌐", serverAddr)
	log.Info().Msgf("Swagger UI available at: http://localhost%s/swagger/index.html 📑", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
