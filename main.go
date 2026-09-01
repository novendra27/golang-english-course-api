package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Setup Zerolog (Console Writer untuk Pretty Output di mode development)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 2. Load Environment Variables dari file .env
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "english-course-api"
	}
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Info().
		Str("app_name", appName).
		Str("port", appPort).
		Msg(fmt.Sprintf("%s skeleton siap dijalankan 🚀", appName))
}
