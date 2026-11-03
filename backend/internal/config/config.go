package config

import "os"

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	JWTSecret    string
	Environment  string
	CookieSecure bool
	CookieName   string
}

func Load() Config {
	env := getenv("APP_ENV", "development")
	return Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:  getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/secaudit?sslmode=disable"),
		JWTSecret:    getenv("JWT_SECRET", "dev-only-secaudit-secret-change-me"),
		Environment:  env,
		CookieSecure: env == "production",
		CookieName:   "secaudit_session",
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
