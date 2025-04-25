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
	userRepo := repository.NewUserRepository(database.DB)
	authService := service.NewAuthService(userRepo)

	authController := NewAuthController(authService)
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
	// jwt required
	g.GET("/my", security.AuthMiddleware(), authController.MyUser)
	g.GET("/auth-status", security.AuthMiddleware(), authController.AuthStatus)
}

func (r *Router) RegisterAccountRoutes(g *gin.RouterGroup) {
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	accountService := service.NewAccountService(accountRepo)
	accountController := NewAccountController(accountService)
	g.POST("", security.AuthMiddleware(), accountController.CreateAccount)
	g.GET("", security.AuthMiddleware(), accountController.GetAccountByUserID)
	g.GET("/all", security.AuthMiddleware(), accountController.GetAccountsAll)
}

func (r *Router) RegisterCardRoutes(g *gin.RouterGroup) {

	cardRepo := repository.CardRepositoryInstance(database.DB)
	cardService := service.NewCardService(cardRepo, "defaultString", []byte("defaultBytes"))
	cardController := NewCardController(cardService)
	g.POST("", security.AuthMiddleware(), cardController.CreateCard)
	// g.GET(":id", security.AuthMiddleware(), cardController.GetCardByID)
	g.GET("", security.AuthMiddleware(), cardController.GetAllCards)
}
