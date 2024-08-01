package util

import (
	"os"

	"github.com/joho/godotenv"
)

func GodotEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		HandleAppError(err, "GodotEnv", "Load", true)
	}
	return os.Getenv(key)
}
