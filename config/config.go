// Package config memuat konfigurasi aplikasi dari environment variables.
// Padanan dari ConfigModule.forRoot({ isGlobal: true }) di NestJS.
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpires  time.Duration
}

// Load membaca .env (jika ada) lalu mengambil nilai dari environment.
func Load() *Config {
	// .env opsional — di production nilai biasanya di-set langsung di environment.
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpires:  getEnvDuration("JWT_EXPIRES_IN", 24*time.Hour),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return d
}
