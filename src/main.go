package main

import (
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"
	"log"
	"os"
)

func main() {
	// Инициализация базы данных
	if err := database.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.CloseDB()

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
