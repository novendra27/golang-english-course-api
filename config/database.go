package config

import (
	"fmt"
	"time"

	"english-course-api/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ConnectDB membuka koneksi ke database PostgreSQL menggunakan GORM
func ConnectDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
		cfg.DB.SSLMode,
		cfg.DB.TimeZone,
	)

	// GORM log mode sesuai APP_ENV
	logLevel := gormlogger.Warn
	if cfg.App.Env == "development" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Dapatkan underlying *sql.DB untuk konfigurasi connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil instance generic database: %w", err)
	}

	// Konfigurasi Connection Pooling
	sqlDB.SetMaxIdleConns(10)                  // Batas koneksi idle di pool
	sqlDB.SetMaxOpenConns(100)                 // Batas maksimal koneksi terbuka
	sqlDB.SetConnMaxLifetime(1 * time.Hour)    // Masa hidup maksimal satu koneksi
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Waktu idle maksimal sebelum koneksi ditutup

	// Uji konektivitas database (Ping)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan ping ke database: %w", err)
	}

	log.Info().
		Str("db_host", cfg.DB.Host).
		Str("db_port", cfg.DB.Port).
		Str("db_name", cfg.DB.Name).
		Msg("Berhasil terhubung ke database PostgreSQL ✅")

	return db, nil
}

// AutoMigrate menjalankan migrasi skema tabel secara otomatis untuk seluruh domain models
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Menjalankan database auto-migration...")

	err := db.AutoMigrate(
		&models.Student{},
		&models.Course{},
		&models.Class{},
		&models.Registration{},
		&models.Payment{},
		&models.ClassPlacement{},
	)
	if err != nil {
		return fmt.Errorf("gagal melakukan auto-migration database: %w", err)
	}

	log.Info().Msg("Database auto-migration berhasil diselesaikan ✅")
	return nil
}

// CloseDB menutup koneksi database secara graceful
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Gagal mengambil instance database saat akan menutup koneksi")
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Error().Err(err).Msg("Gagal menutup koneksi database")
		return
	}
	log.Info().Msg("Koneksi database berhasil ditutup")
}
