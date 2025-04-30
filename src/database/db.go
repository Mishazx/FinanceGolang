package database

import (
	"FinanceGolang/src/model"
	"os"
	"time"

	"fmt"
	"log"

	// "github.com/glebarez/sqlite" // Заменили импорт на pure Go реализацию
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DBType string

const (
	SQLite   DBType = "sqlite"
	Postgres DBType = "postgres"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// InitDB - инициализирует соединение с базой данных
func InitDB() (*gorm.DB, error) {
	// Определяем тип базы данных из переменной окружения
	// По умолчанию используем SQLite
	dbType := DBType(os.Getenv("DB_TYPE"))
	if dbType == "" {
		dbType = SQLite
	}

	// Настройка логгера GORM
	loggerInstance := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	var err error
	var dialector gorm.Dialector

	switch dbType {
	case SQLite:
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "bank.db"
		}
		dialector = sqlite.Open(dbPath)
	case Postgres:
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		if password == "" {
			password = "postgres"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "bank"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
		dialector = postgres.Open(dsn)
	default:
		log.Fatalf("Неподдерживаемый тип базы данных: %s", dbType)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{Logger: loggerInstance}) // Укажите вашу базу данных
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	err = CreateTables(DB)
	if err != nil {
		log.Fatalf("Ошибка при создании таблиц: %v", err)
	}

	log.Printf("Успешное подключение к базе данных типа: %s", dbType)

	return DB, nil
}

// CloseDB - закрывает соединение с базой данных
func CloseDB() {
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			log.Printf("Ошибка при получении соединения с базой данных: %v", err)
			return
		}
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с базой данных: %v", err)
		}
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
	// Создаем таблицы
	err := db.AutoMigrate(
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

	if err != nil {
		log.Fatalf("Ошибка при создании таблиц: %v", err)
		return err
	}

	createAdmin(db)

	return InitializeRoles(db)
}

func createAdmin(db *gorm.DB) error {
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

// InitializeRoles создает базовые роли в системе, если они еще не существуют
func InitializeRoles(db *gorm.DB) error {
	// Проверяем, существуют ли уже роли
	var count int64
	db.Model(&model.Role{}).Count(&count)

	// Если роли уже существуют, ничего не делаем
	if count > 0 {
		log.Println("Роли уже существуют в базе данных")
		return nil
	}

	// Создаем базовые роли
	roles := []model.Role{
		{Name: model.RoleAdmin, Description: "Администратор"},
		{Name: model.RoleUser, Description: "Пользователь"},
	}

	// Сохраняем роли в базе данных
	result := db.Create(&roles)
	if result.Error != nil {
		log.Printf("Ошибка при создании ролей: %v", result.Error)
		return result.Error
	}

	log.Printf("Создано %d ролей", result.RowsAffected)
	return nil
}
