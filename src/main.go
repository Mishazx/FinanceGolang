package main

import (
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	database.CreateTables(database.DB)

	// Инициализация ролей
	roleRepo := repository.NewRoleRepository(database.DB)
	roleService := service.NewRoleService(roleRepo)
	if err := roleService.InitializeDefaultRoles(); err != nil {
		log.Fatalf("Failed to initialize roles: %v", err)
	}

	r := gin.Default()

	router := controller.NewRouter()

	router.RegisterAuthRoutes(r.Group("/api/auth"))
	router.RegisterAccountRoutes(r.Group("/api/accounts"))
	router.RegisterCardRoutes(r.Group("/api/cards"))
	router.RegisterKeyRateRoutes(r.Group("/api/key-rate"))

	// router.

	r.Run(":8080")
}
