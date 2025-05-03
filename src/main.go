package main

import (
	"FinanceGolang/src/config"
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"
	"fmt"
	"log"
)

func main() {
	// Загрузка конфигурации
	if err := config.Init(); err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	cfg := config.Get()

	// Инициализация базы данных
	db, err := database.InitDB()
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
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Сервер запускается на %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
