package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	valueToInt, err := strconv.Atoi(value)

	if !exists || err != nil {
		return fallback
	}
	return valueToInt
}
