package database

import (
	"log"

	"gorm.io/driver/sqlite" // Или другой драйвер для вашей базы данных
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) // Укажите вашу базу данных
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
}

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
