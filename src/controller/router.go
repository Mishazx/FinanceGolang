package controller

import (
	"FinanceGolang/src/database"
	"FinanceGolang/src/repository" // Импорт пакета для account_controller
	"FinanceGolang/src/security"
	"FinanceGolang/src/service" // Импорт пакета для auth_service

	"github.com/gin-gonic/gin"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) RegisterAuthRoutes(g *gin.RouterGroup) {
	db, err := database.ConnectDB()
	if err != nil {
		panic("Failed to connect to database")
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)

	authController := NewAuthController(authService)
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
	// jwt required
	g.GET("/protected", security.AuthMiddleware(), authController.Protected)
	g.GET("/auth-status", security.AuthMiddleware(), authController.AuthStatus)
}

func (r *Router) RegisterAccountRoutes(g *gin.RouterGroup) {
	db, err := database.ConnectDB()
	if err != nil {
		panic("Failed to connect to database")
	}
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo)
	accountController := NewAccountController(accountService)
	g.POST("/", security.AuthMiddleware(), accountController.CreateAccount)
	g.GET("/:id", security.AuthMiddleware(), accountController.GetAccountByID)
	g.GET("/", security.AuthMiddleware(), accountController.GetAccounts)
	// g.PUT("/accounts/:id", accountController.UpdateAccount)
	// g.DELETE("/accounts/:id", accountController.DeleteAccount)
}
