package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application settings loaded from environment variables.
type Config struct {
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	AppPort              string
	CheckIntervalMinutes int
	CheckRetentionDays   int
}

// Load reads configuration from environment variables with sensible defaults.
func Load() Config {
	interval, err := strconv.Atoi(getEnv("CHECK_INTERVAL_MINUTES", "144"))
	if err != nil || interval < 1 {
		interval = 144
	}

	retention, err := strconv.Atoi(getEnv("CHECK_RETENTION_DAYS", "7"))
	if err != nil || retention < 1 {
		retention = 7
	}

	return Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBUser:               getEnv("DB_USER", "uptime"),
		DBPassword:           getEnv("DB_PASSWORD", "uptimepassword"),
		DBName:               getEnv("DB_NAME", "uptime"),
		// PORT is set by Railway, Render, and similar hosts.
		AppPort:              getEnv("PORT", getEnv("APP_PORT", "8080")),
		CheckIntervalMinutes: interval,
		CheckRetentionDays:   retention,
	}
}

// DSN returns the MySQL connection string.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
