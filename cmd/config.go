package main

import (
	"time"
)

type config struct {
	AppName         string
	AppPort         string
	AppEnv          string
	DatabaseURL     string
	JWTSecret       string
	AllowedOrigins  []string
	RateLimitMax    int
	RateLimitWindow time.Duration
}

func loadConfig() config {
	return config{
		AppName:         getEnv("APP_NAME", "go-backend-template"),
		AppPort:         getEnv("APP_PORT", "3000"),
		AppEnv:          getEnv("APP_ENV", "development"),
		DatabaseURL:     mustGetEnv("DATABASE_URL"), // required — ถ้าไม่มี exit ทันที
		JWTSecret:       mustGetEnv("SUPABASE_JWT_SECRET"),
		AllowedOrigins:  splitEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		RateLimitMax:    getEnvInt("RATE_LIMIT_MAX", 60),
		RateLimitWindow: time.Duration(getEnvInt("RATE_LIMIT_WINDOW", 1)) * time.Minute,
	}
}
