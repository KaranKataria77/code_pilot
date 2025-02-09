package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading ENV")
	}
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}
