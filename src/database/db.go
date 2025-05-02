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
		&model.Role{},
		&model.User{},
		&model.UserRole{},
		&model.Account{},
		&model.Card{},
		&model.Transaction{},
		&model.Credit{},
		&model.PaymentSchedule{},
		&model.Analytics{},
		&model.BalanceForecast{},
	)

	if err != nil {
		return fmt.Errorf("ошибка при создании таблиц: %v", err)
	}

	// Инициализируем роли после создания таблиц
	if err := InitializeRoles(db); err != nil {
		return fmt.Errorf("ошибка при инициализации ролей: %v", err)
	}

	// Создаем админа после создания всех таблиц и инициализации ролей
	if err := createAdmin(db); err != nil {
		return fmt.Errorf("ошибка при создании админа: %v", err)
	}

	return nil
}

func createAdmin(db *gorm.DB) error {
	// Проверяем, существует ли уже админ
	var existingAdmin model.User
	if err := db.Where("email = ?", "admin@example.com").First(&existingAdmin).Error; err == nil {
		// Админ уже существует, ничего не делаем
		return nil
	}

	admin := &model.User{
		Username: "admin",
		Password: "admin",
		Email:    "admin@example.com",
	}

	// Создаем роли, если их нет
	adminRole := model.Role{Name: model.RoleAdmin, Description: "Администратор системы"}
	userRole := model.Role{Name: model.RoleUser, Description: "Обычный пользователь"}

	if err := db.FirstOrCreate(&adminRole, model.Role{Name: model.RoleAdmin}).Error; err != nil {
		return fmt.Errorf("ошибка при создании роли админа: %v", err)
	}

	if err := db.FirstOrCreate(&userRole, model.Role{Name: model.RoleUser}).Error; err != nil {
		return fmt.Errorf("ошибка при создании роли пользователя: %v", err)
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("ошибка при создании админа: %v", err)
	}

	return nil
}

// InitializeRoles создает базовые роли в системе
func InitializeRoles(db *gorm.DB) error {
	// Создаем базовые роли
	roles := []model.Role{
		{Name: model.RoleAdmin, Description: "Администратор"},
		{Name: model.RoleUser, Description: "Пользователь"},
	}

	// Сохраняем роли в базе данных
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			return fmt.Errorf("ошибка при создании роли %s: %v", role.Name, err)
		}
	}

	return nil
}

func addNumberField(db *gorm.DB) error {
	// Обновляем существующие записи
	var accounts []model.Account
	if err := db.Find(&accounts).Error; err != nil {
		return fmt.Errorf("ошибка при получении счетов: %v", err)
	}

	for _, account := range accounts {
		if account.Number == "" {
			account.Number = model.GenerateAccountNumber()
			if err := db.Save(&account).Error; err != nil {
				return fmt.Errorf("ошибка при обновлении счета %d: %v", account.ID, err)
			}
		}
	}

	return nil
}
