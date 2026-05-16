package config

import (
	"os"
	"time"
)

type Config struct {
	AppEnv      string
	DatabaseURL string
	HTTPPort    string

	JWTSecret    string
	JWTExpiresIn time.Duration

	WhatsAppProvider     string
	WhatsAppAdminEnabled bool
	AdminEmail           string
}

func Load() *Config {
	jwtExpiry := parseDuration(getEnv("JWT_EXPIRES_IN", "24h"), 24*time.Hour)

	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: mustEnv("DATABASE_URL"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),

		JWTSecret:    mustEnv("JWT_SECRET"),
		JWTExpiresIn: jwtExpiry,

		WhatsAppProvider:     getEnv("WHATSAPP_PROVIDER", "mock"),
		WhatsAppAdminEnabled: getEnv("WHATSAPP_ADMIN_ENABLED", "true") == "true",
		AdminEmail:           getEnv("ADMIN_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("variável de ambiente obrigatória não definida: " + key)
	}
	return v
}

func parseDuration(s string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback
	}
	return d
}
