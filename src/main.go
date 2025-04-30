package main

import (
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"

	// "FinanceGolang/src/model"
	"log"
	"os"
)

func main() {
	// Инициализация базы данных
	db, err := database.InitDB() // Get both the DB and error
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.CloseDB()

	// Auto migrate models
	err = database.CreateTables(db)
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	// Инициализация контроллеров напрямую через Router
	// Router создает все необходимые репозитории и сервисы внутри себя

	// Инициализация контроллеров
	router := controller.NewRouter()

	// Настройка Gin и middleware
	r := router.InitRoutes()

	// Запуск сервера

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
