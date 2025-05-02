package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var SecretKey []byte

func InitConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables")
	}

	// Get secret from env
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatal("SECRET_KEY is not set in environment")
	}
	SecretKey = []byte(secret)
}
