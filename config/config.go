// Package config handles environment variable loading and database initialization.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// AppConfig holds core application configurations.
type AppConfig struct {
	Env  string
	Port string
	Name string
}

// DBConfig holds PostgreSQL connection credentials loaded from environment variables.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

// LogConfig holds logging configuration parameters.
type LogConfig struct {
	Level  string
	Pretty bool
}

// Config is the root configuration container.
type Config struct {
	App AppConfig
	DB  DBConfig
	Log LogConfig
}

// LoadConfig loads environment variables from a .env file or system environment.
func LoadConfig() (*Config, error) {
	// 1. Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg(".env file not found, loading from system environment variables")
	}

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Name: getEnv("APP_NAME", "english-course-api"),
		},
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Pretty: os.Getenv("LOG_PRETTY") == "true",
		},
	}

	// 2. Validate mandatory database environment variables
	if err := cfg.validateDBConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validateDBConfig validates that all required database credentials are present.
func (c *Config) validateDBConfig() error {
	missing := []string{}

	if c.DB.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.DB.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.DB.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.DB.Name == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required database environment variables: %v", missing)
	}

	return nil
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
