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

// ConnectDB establishes a PostgreSQL database connection using GORM.
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

	// GORM log mode based on APP_ENV
	logLevel := gormlogger.Warn
	if cfg.App.Env == "development" {
		logLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Retrieve underlying *sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve generic database instance: %w", err)
	}

	// Connection Pool Configuration
	sqlDB.SetMaxIdleConns(10)                  // Max idle connections in pool
	sqlDB.SetMaxOpenConns(100)                 // Max open connections to database
	sqlDB.SetConnMaxLifetime(1 * time.Hour)    // Max connection lifetime
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Max connection idle timeout

	// Test database connectivity (Ping)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("db_host", cfg.DB.Host).
		Str("db_port", cfg.DB.Port).
		Str("db_name", cfg.DB.Name).
		Msg("Connected to PostgreSQL database successfully ✅")

	return db, nil
}

// AutoMigrate runs automatic database schema migrations for all domain models.
func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Running database auto-migrations...")

	err := db.AutoMigrate(
		&models.Student{},
		&models.Course{},
		&models.Class{},
		&models.Registration{},
		&models.Payment{},
		&models.ClassPlacement{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	log.Info().Msg("Database auto-migration completed successfully ✅")
	return nil
}

// CloseDB gracefully closes the database connection.
func CloseDB(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve generic database instance during shutdown")
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
		return
	}
	log.Info().Msg("Database connection closed gracefully 🔒")
}
