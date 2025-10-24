package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func InitENV() error {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return nil
}
