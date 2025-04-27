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

func (r *Router) RegisterUserRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	userService := service.NewUserService(userRepo)
	userController := NewUserController(userService)

	users := g.Group("/users")
	users.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}))
	{
		users.GET("/me", userController.GetCurrentUser)
		users.PUT("/me", userController.UpdateCurrentUser)
		users.DELETE("/me", userController.DeleteCurrentUser)
	}
}

func (r *Router) RegisterAccountRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	accountService := service.NewAccountService(accountRepo)
	accountController := NewAccountController(accountService)

	// Создаем сервис и контроллер для операций со счетами
	transactionRepo := repository.NewTransactionRepository(database.DB)
	accountOperationService := service.NewAccountOperationService(accountRepo, transactionRepo)
	accountOperationController := NewAccountOperationController(accountOperationService)

	// Базовые операции со счетами
	g.POST("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.CreateAccount)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.GetAccountByUserID)
	g.GET("/all", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), accountController.GetAccountsAll)

	// Операции с конкретным счетом
	accountGroup := g.Group("/:id")
	accountGroup.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}))
	{
		accountGroup.POST("/deposit", accountOperationController.Deposit)
		accountGroup.POST("/withdraw", accountOperationController.Withdraw)
		accountGroup.POST("/transfer", accountOperationController.Transfer)
		accountGroup.GET("/transactions", accountOperationController.GetTransactions)
	}
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
	externalService := service.NewExternalService("", 0, "", "", "") // Здесь нужно передать реальные параметры SMTP
	cbrController := NewCbrController(externalService)

	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}), cbrController.GetKeyRate)
}

// Кредитные операции
func (r *Router) RegisterCreditRoutes(g *gin.RouterGroup) {
	creditService := service.NewCreditService(
		repository.NewCreditRepository(database.DB),
		repository.AccountRepositoryInstance(database.DB),
		repository.NewTransactionRepository(database.DB),
		service.NewExternalService("", 0, "", "", ""), // Заменяем NewCbrService на NewExternalService
	)
	creditController := NewCreditController(creditService)

	credits := g.Group("/credits")
	credits.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: service.NewAuthCheckService(repository.NewUserRepository(database.DB)).ValidateUserFromToken,
	}))
	{
		credits.POST("", creditController.CreateCredit)
		credits.GET("", creditController.GetUserCredits)
		credits.GET("/:id", creditController.GetCreditByID)
		credits.GET("/:id/schedule", creditController.GetPaymentSchedule)
		credits.POST("/:id/payment", creditController.ProcessPayment)
	}
}

func (r *Router) RegisterAnalyticsRoutes(g *gin.RouterGroup) {
	// analyticsRepo := repository.NewAnalyticsRepository(database.DB)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	creditRepo := repository.NewCreditRepository(database.DB)
	analyticsService := service.NewAnalyticsService(transactionRepo, accountRepo, creditRepo)
	analyticsController := NewAnalyticsController(analyticsService)

	analytics := g.Group("/analytics")
	analytics.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: service.NewAuthCheckService(repository.NewUserRepository(database.DB)).ValidateUserFromToken,
	}))
	{
		analytics.POST("", analyticsController.GetAnalytics)
		analytics.GET("/accounts/:id/forecast", analyticsController.GetBalanceForecast)
	}
}

func (r *Router) RegisterAccountOperationRoutes(g *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(database.DB)
	authCheckService := service.NewAuthCheckService(userRepo)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	accountOperationService := service.NewAccountOperationService(accountRepo, transactionRepo)
	accountOperationController := NewAccountOperationController(accountOperationService)

	accounts := g.Group("/accounts")
	accounts.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authCheckService.ValidateUserFromToken,
	}))
	{
		accounts.POST("/:id/deposit", accountOperationController.Deposit)
		accounts.POST("/:id/withdraw", accountOperationController.Withdraw)
		accounts.POST("/:id/transfer", accountOperationController.Transfer)
		accounts.GET("/:id/transactions", accountOperationController.GetTransactions)
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.Default()
	
	api := router.Group("/api")
	{
		r.RegisterAuthRoutes(api)
		r.RegisterUserRoutes(api)
		r.RegisterAccountRoutes(api)
		r.RegisterCardRoutes(api)
		r.RegisterAccountOperationRoutes(api)
		r.RegisterCreditRoutes(api)
		r.RegisterAnalyticsRoutes(api)
	}
	
	return router
}
