package db

import (
	"log"

	"ToDo/internal/model"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB()(*gorm.DB, error){
	dsn := "host=localhost port=5432 user=postgres password=password sslmode=disable database=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	if err = db.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	return db, nil
}