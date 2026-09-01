package main

import (
	"os"

	"english-course-api/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 1. Setup Zerolog (Pretty output untuk console)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// 2. Load Konfigurasi murni dari .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Gagal memuat konfigurasi aplikasi")
	}

	// Sesuaikan log level dari config
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Info().
		Str("app_name", cfg.App.Name).
		Str("env", cfg.App.Env).
		Str("port", cfg.App.Port).
		Msg("Memulai inisialisasi aplikasi 🚀")

	// 3. Inisialisasi Koneksi Database PostgreSQL
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Database connection warning (pastikan PostgreSQL service sedang berjalan)")
	} else {
		defer config.CloseDB(db)

		// 4. Auto-Migration skema 6 domain model
		if err := config.AutoMigrate(db); err != nil {
			log.Fatal().Err(err).Msg("Gagal melakukan auto-migration database")
		}
	}

	log.Info().Msg("Aplikasi siap melanjutkan ke fase berikutnya! ✨")
}
