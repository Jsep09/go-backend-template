package main

import (
	"fmt"
	"log/slog"
	"os"
)

// Helper functions
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// mustGetEnv — required config, exit ถ้าไม่มี
func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("required environment variable not set", "key", key)
		os.Exit(1)
	}
	return val
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	var result int
	if _, err := fmt.Sscan(val, &result); err != nil {
		return defaultVal
	}
	return result
}

func splitEnv(key, defaultVal string) []string {
	val := getEnv(key, defaultVal)
	var result []string
	for _, s := range splitString(val, ",") {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
