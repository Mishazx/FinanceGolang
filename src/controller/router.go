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
	authCheckService := service.NewAuthCheckService(userRepo)

	authController := NewAuthController(authService)
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)
	// jwt required
	g.GET("/my", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), authController.MyUser)
	g.GET("/auth-status", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), authController.AuthStatus)
}

func (r *Router) RegisterAccountRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	accountService := service.NewAccountService(accountRepo)
	accountController := NewAccountController(accountService)
	g.POST("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.CreateAccount)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.GetAccountByUserID)
	g.GET("/all", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.GetAccountsAll)
}

func (r *Router) RegisterCardRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	cardRepo := repository.CardRepositoryInstance(database.DB)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	cardService := service.NewCardService(cardRepo, accountRepo, "defaultString", []byte("defaultBytes"))
	cardController := NewCardController(cardService)
	g.POST("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), cardController.CreateCard)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), cardController.GetAllCards)
}

func (r *Router) RegisterKeyRateRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	cbrService := service.NewCbrService()
	cbrController := NewCbrController(cbrService)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), cbrController.GetKeyRate)
}
