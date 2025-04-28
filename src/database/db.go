package database

import (
	"FinanceGolang/src/model"

	"log"

	"github.com/glebarez/sqlite" // Заменили импорт на pure Go реализацию
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
	return db.AutoMigrate(
		&model.User{}, 
		&model.Account{}, 
		&model.Card{}, 
		&model.Role{}, 
		&model.UserRole{}, 
		&model.Transaction{}, 
		&model.Credit{}, 
		&model.PaymentSchedule{},
		&model.Analytics{},
		&model.BalanceForecast{},
	)
}

func CreateAdmin(db *gorm.DB) error {
	admin := &model.User{
		Username: "admin",
		Password: "admin",
		Email:    "admin@example.com",
	}

	// Создаем роли, если их нет
	adminRole := model.Role{Name: "admin", Description: "Администратор системы"}
	userRole := model.Role{Name: "user", Description: "Обычный пользователь"}

	if err := db.FirstOrCreate(&adminRole, model.Role{Name: "admin"}).Error; err != nil {
		log.Fatalf("Failed to create admin role: %v", err)
	}

	if err := db.FirstOrCreate(&userRole, model.Role{Name: "user"}).Error; err != nil {
		log.Fatalf("Failed to create user role: %v", err)
	}

	return db.Create(admin).Error
}
