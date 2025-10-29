package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")

	dsn := fmt.Sprintf("host=%s port=%s database=%s password=%s user=%s sslmode=disable",
		dbHost, dbPort, dbName, dbPass, dbUser)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	return db, nil
}
