package main

import (
	"FinanceGolang/src/controller"
	"FinanceGolang/src/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDatabase()

	database.CreateTables(database.DB)

	r := gin.Default()

	router := controller.NewRouter()

	router.RegisterAuthRoutes(r.Group("/api/auth"))
	router.RegisterAccountRoutes(r.Group("/api/accounts"))
	router.RegisterCardRoutes(r.Group("/api/cards"))
	router.RegisterKeyRateRoutes(r.Group("/api/key-rate"))

	// router.

	r.Run(":8080")
}
