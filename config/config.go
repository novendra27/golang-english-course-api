package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// AppConfig menampung konfigurasi aplikasi
type AppConfig struct {
	Env  string
	Port string
	Name string
}

// DBConfig menampung parameter koneksi database PostgreSQL murni dari .env
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

// LogConfig menampung parameter konfigurasi logging
type LogConfig struct {
	Level  string
	Pretty bool
}

// Config adalah root container konfigurasi
type Config struct {
	App AppConfig
	DB  DBConfig
	Log LogConfig
}

// LoadConfig memuat environment variables murni dari file .env
func LoadConfig() (*Config, error) {
	// 1. Muat file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("File .env tidak ditemukan, membaca dari environment variables sistem")
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

	// 2. Validasi field database wajib (wajib ada di .env)
	if err := cfg.validateDBConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validateDBConfig memvalidasi apakah semua kredensial database yang dibutuhkan sudah diset di .env
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
		return fmt.Errorf("variabel database wajib belum diset di .env: %v", missing)
	}

	return nil
}

// getEnv membaca env var atau mengembalikan fallback jika tidak diset (khusus non-kredensial)
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
