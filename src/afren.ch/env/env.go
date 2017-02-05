package env

import (
	"os"
	"strconv"
	"log"
)

func RequiredInt(key string) int {
	envStr := os.Getenv(key)
	if len(envStr) == 0 { fatal(key) }

	envInt, err := strconv.Atoi(envStr)
	if err != nil { fatal(key) }

	return envInt
}

func OptionalInt(key string, fallback int) int {
	envStr := os.Getenv(key)
	if len(envStr) == 0 {
		warn(key)
		return fallback
	}

	envInt, err := strconv.Atoi(envStr)
	if err != nil {
		warn(key)
		return fallback
	}

	return envInt
}

func RequiredString(key string) string {
	envStr := os.Getenv(key)
	if len(envStr) == 0 { fatal(key) }

	return envStr
}

func OptionalString(key, fallback string) string {
	envStr := os.Getenv(key)
	if len(envStr) == 0 {
		warn(key)
		return fallback
	}

	return envStr
}

func fatal(key string) {
	log.Fatalf("Failed to parse %s, exiting", key)
}

func warn(key string) {
	log.Printf("Failed to parse: %s, continuing", key)
}
