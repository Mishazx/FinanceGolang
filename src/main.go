package main

import (
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	if err := database.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.CloseDB()

	// Инициализация репозиториев
	roleRepo := repository.NewRoleRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	accountRepo := repository.NewAccountRepository(database.DB)
	cardRepo := repository.NewCardRepository(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	creditRepo := repository.NewCreditRepository(database.DB)

	// Инициализация сервисов
	roleService := service.NewRoleService(roleRepo)
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	cardService := service.NewCardService(cardRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	creditService := service.NewCreditService(creditRepo)

	// Инициализация контроллеров
	router := controller.NewRouter(
		userService,
		roleService,
		accountService,
		cardService,
		transactionService,
		creditService,
	)

	// Настройка Gin
	r := gin.Default()

	// Регистрация маршрутов
	router.RegisterAuthRoutes(r.Group("/api/auth"))
	router.RegisterAccountRoutes(r.Group("/api/accounts"))
	router.RegisterCardRoutes(r.Group("/api/cards"))
	router.RegisterTransactionRoutes(r.Group("/api/transactions"))
	router.RegisterCreditRoutes(r.Group("/api/credits"))
	router.RegisterAdminRoutes(r.Group("/api/admin"))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
