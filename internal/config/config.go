package config

import (
	"fmt"
	"os"
) // baru

// Config menampung semua pengaturan aplikasi yang dibaca dari environment.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
}

// Load mengisi Config dari environment variable.
func Load() Config { // baru
	return Config{
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "account_user"),
		DBPassword: getEnv("DB_PASSWORD", "account_pass"),
		DBName:     getEnv("DB_NAME", "account_db"),
		AppPort:    getEnv("APP_PORT", "8080"),
	}
}

// getEnv mengembalikan nilai environment, atau fallback jika kosong.
func getEnv(key, fallback string) string { // baru
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// DSN menyusun string koneksi PostgreSQL untuk GORM.
func (c Config) DSN() string { // baru
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}