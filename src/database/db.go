package database

import (
	"FinanceGolang/src/model"

	"log"

	"gorm.io/driver/sqlite" // Или другой драйвер для вашей базы данных
	"gorm.io/gorm"
)

var DB *gorm.DB

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

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

func CreateTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Account{}, &model.Card{}, &model.Role{}, &model.UserRole{})
}
