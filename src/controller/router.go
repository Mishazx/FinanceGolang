package controller

import (
	"FinanceGolang/src/database"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	"FinanceGolang/src/service"
	"net/http"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Константы для API путей
const (
	APIPathRegister     = "/register"
	APIPathLogin        = "/login"
	APIPathMyUser       = "/my"
	APIPathAuthStatus   = "/auth-status"
	APIPathUsers        = "/users"
	APIPathMe           = "/me"
	APIPathAccounts     = "/accounts"
	APIPathAll          = "/all"
	APIPathDeposit      = "/deposit"
	APIPathWithdraw     = "/withdraw"
	APIPathTransfer     = "/transfer"
	APIPathTransactions = "/transactions"
	APIPathCards        = "/cards"
	APIPathCredits      = "/credits"
	APIPathSchedule     = "/schedule"
	APIPathPayment      = "/payment"
	APIPathAnalytics    = "/analytics"
	APIPathForecast     = "/forecast"
)

// Константы для сообщений об ошибках
const (
	ErrInvalidCredentials = "Неверные учетные данные"
	ErrUserNotFound       = "Пользователь не найден"
	ErrAccountNotFound    = "Счет не найден"
	ErrInsufficientFunds  = "Недостаточно средств"
	ErrInvalidAmount      = "Неверная сумма"
	ErrInvalidToken       = "Неверный токен"
	ErrUnauthorized       = "Не авторизован"
	ErrInternalServer     = "Внутренняя ошибка сервера"
)

// Константы для успешных операций
const (
	MsgRegistrationSuccess = "Регистрация успешно завершена"
	MsgLoginSuccess        = "Вход выполнен успешно"
	MsgDepositSuccess      = "Средства успешно зачислены"
	MsgWithdrawSuccess     = "Средства успешно сняты"
	MsgTransferSuccess     = "Перевод выполнен успешно"
	MsgPaymentSuccess      = "Платеж выполнен успешно"
)

type Router struct{}

// NewRouter создает новый экземпляр маршрутизатора
func NewRouter() *Router {
	return &Router{}
}

// createAuthService создает сервис аутентификации
func (r *Router) createAuthService() service.AuthService {
	userRepo := repository.NewUserRepository(database.DB)
	return service.NewAuthService(userRepo)
}

// createUserService создает сервис пользователей
func (r *Router) createUserService() service.UserService {
	userRepo := repository.NewUserRepository(database.DB)
	return service.NewUserService(userRepo)
}

// createAccountService создает сервис счетов
func (r *Router) createAccountService() service.AccountService {
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	return service.NewAccountService(accountRepo, transactionRepo)
}

// createCardService создает сервис карт
func (r *Router) createCardService() service.CardService {
	cardRepo := repository.CardRepositoryInstance(database.DB)
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	return service.NewCardService(cardRepo, accountRepo, "defaultString", []byte("defaultBytes"))
}

// createCreditService создает сервис кредитов
func (r *Router) createCreditService() service.CreditService {
	return service.NewCreditService(
		repository.NewCreditRepository(database.DB),
		repository.AccountRepositoryInstance(database.DB),
		repository.NewTransactionRepository(database.DB),
		service.NewExternalService("", 0, "", "", ""),
	)
}

// createAnalyticsService создает сервис аналитики
func (r *Router) createAnalyticsService() *service.AnalyticsService {
	accountRepo := repository.AccountRepositoryInstance(database.DB)
	transactionRepo := repository.NewTransactionRepository(database.DB)
	creditRepo := repository.NewCreditRepository(database.DB)
	return service.NewAnalyticsService(transactionRepo, accountRepo, creditRepo)
}

// LoggerMiddleware логирует информацию о запросах
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Продолжаем обработку запроса
		c.Next()

		// Логируем информацию после обработки
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		logrus.WithFields(logrus.Fields{
			"status":  status,
			"latency": latency,
			"client":  clientIP,
			"method":  method,
			"path":    path,
		}).Info("Обработан HTTP запрос")
	}
}

// ErrorHandlerMiddleware обрабатывает ошибки и возвращает стандартизированный ответ
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Проверяем наличие ошибок
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			status := http.StatusInternalServerError
			message := ErrInternalServer

			// Определяем тип ошибки и устанавливаем соответствующий статус
			switch err.(type) {
			case *gin.Error:
				status = http.StatusBadRequest
				message = err.Error()
			case error:
				if err.Error() == "unauthorized" {
					status = http.StatusUnauthorized
					message = ErrUnauthorized
				}
			}

			c.JSON(status, gin.H{
				"error": message,
			})
			c.Abort()
		}
	}
}

// CORSMiddleware обрабатывает CORS заголовки
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware ограничивает количество запросов
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := make(map[string]time.Time)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		if last, exists := limiter[ip]; exists {
			if now.Sub(last) < time.Second {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Слишком много запросов. Пожалуйста, подождите.",
				})
				c.Abort()
				return
			}
		}

		limiter[ip] = now
		c.Next()
	}
}

// CacheMiddleware создает middleware для кэширования
func CacheMiddleware(store *persistence.InMemoryStore) gin.HandlerFunc {
	return cache.CachePage(store, time.Minute*5, func(c *gin.Context) {
		c.Next()
	})
}

// CompressionMiddleware сжимает ответы
func CompressionMiddleware() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

// RegisterAuthRoutes регистрирует маршруты аутентификации
func (r *Router) RegisterAuthRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	authController := NewAuthController(authService)

	g.POST(APIPathRegister, authController.Register)
	g.POST(APIPathLogin, authController.Login)
	g.GET(APIPathMyUser, security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), authController.MyUser)
	g.GET(APIPathAuthStatus, security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), authController.AuthStatus)
}

// RegisterUserRoutes регистрирует маршруты пользователей
func (r *Router) RegisterUserRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	userService := r.createUserService()
	userController := NewUserController(userService)

	users := g.Group(APIPathUsers)
	users.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}))
	{
		users.GET(APIPathMe, userController.GetCurrentUser)
		users.PUT(APIPathMe, userController.UpdateCurrentUser)
		users.DELETE(APIPathMe, userController.DeleteCurrentUser)
	}
}

// RegisterAccountRoutes регистрирует маршруты счетов
func (r *Router) RegisterAccountRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	accountService := r.createAccountService()
	accountController := NewAccountController(accountService)

	g.POST("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), accountController.CreateAccount)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), accountController.GetAccountByUserID)
	g.GET(APIPathAll, security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), accountController.GetAccountsAll)

	accountGroup := g.Group("/:id")
	accountGroup.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}))
	{
		accountGroup.POST(APIPathDeposit, accountController.Deposit)
		accountGroup.POST(APIPathWithdraw, accountController.Withdraw)
		accountGroup.POST(APIPathTransfer, accountController.Transfer)
		accountGroup.GET(APIPathTransactions, accountController.GetTransactions)
	}
}

// RegisterCardRoutes регистрирует маршруты карт
func (r *Router) RegisterCardRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	cardService := r.createCardService()
	cardController := NewCardController(cardService)

	g.POST("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), cardController.CreateCard)
	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), cardController.GetAllCards)
}

// RegisterKeyRateRoutes регистрирует маршруты ключевой ставки
func (r *Router) RegisterKeyRateRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	externalService := service.NewExternalService("", 0, "", "", "")
	cbrController := NewCbrController(externalService)

	g.GET("", security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}), cbrController.GetKeyRate)
}

// RegisterCreditRoutes регистрирует маршруты кредитов
func (r *Router) RegisterCreditRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	creditService := r.createCreditService()
	creditController := NewCreditController(creditService)

	credits := g.Group(APIPathCredits)
	credits.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}))
	{
		credits.POST("", creditController.CreateCredit)
		credits.GET("", creditController.GetUserCredits)
		credits.GET("/:id", creditController.GetCreditByID)
		credits.GET("/:id"+APIPathSchedule, creditController.GetPaymentSchedule)
		credits.POST("/:id"+APIPathPayment, creditController.ProcessPayment)
	}
}

// RegisterAnalyticsRoutes регистрирует маршруты аналитики
func (r *Router) RegisterAnalyticsRoutes(g *gin.RouterGroup) {
	authService := r.createAuthService()
	analyticsService := r.createAnalyticsService()
	analyticsController := NewAnalyticsController(analyticsService)

	analytics := g.Group(APIPathAnalytics)
	analytics.Use(security.AuthMiddleware(security.AuthMiddlewareDeps{
		ValidateUserFromToken: authService.ValidateUserFromToken,
	}))
	{
		analytics.POST("", analyticsController.GetAnalytics)
		analytics.GET("/accounts/:id"+APIPathForecast, analyticsController.GetBalanceForecast)
	}
}

// InitRoutes инициализирует все маршруты приложения
func (r *Router) InitRoutes() *gin.Engine {
	router := gin.Default()

	// Создаем хранилище для кэша
	store := persistence.NewInMemoryStore(time.Minute)

	// Добавляем middleware
	router.Use(LoggerMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())
	router.Use(CompressionMiddleware())
	router.Use(CacheMiddleware(store))

	api := router.Group("/api")
	{
		r.RegisterAuthRoutes(api)
		r.RegisterUserRoutes(api)
		r.RegisterAccountRoutes(api)
		r.RegisterCardRoutes(api)
		r.RegisterCreditRoutes(api)
		r.RegisterAnalyticsRoutes(api)
	}

	return router
}
