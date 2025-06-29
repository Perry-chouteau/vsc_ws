package internal

import (
	"log"
	"os"
)

func getEnv(s string) string {
	jsonPath := os.Getenv(s)
	if jsonPath == "" {
		log.Fatal("Environment variable " + s + " is not set")
	}
	return jsonPath
}
