package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"time"
)

type AnalyticsService interface {
	GetAnalytics(request *model.AnalyticsRequest) (*model.Analytics, error)
	GetBalanceForecast(accountID uint, days int) ([]model.BalanceForecast, error)
}

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
	accountRepo   repository.AccountRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository, accountRepo repository.AccountRepository) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
		accountRepo:   accountRepo,
	}
}

func (s *analyticsService) GetAnalytics(request *model.AnalyticsRequest) (*model.Analytics, error) {
	// Получаем транзакции за период
	transactions, err := s.analyticsRepo.GetTransactionsByPeriod(request.AccountID, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	// Получаем кредитные платежи
	creditPayments, err := s.analyticsRepo.GetCreditPaymentsByPeriod(request.AccountID, request.StartDate, request.EndDate)
	if err != nil {
		return nil, err
	}

	// Создаем объект аналитики
	analytics := &model.Analytics{
		AccountID:     request.AccountID,
		Period:        request.Period,
		StartDate:     request.StartDate,
		EndDate:       request.EndDate,
		Categories:    make(map[model.TransactionCategory]float64),
		CreditPayments: creditPayments,
	}

	// Анализируем транзакции
	for _, t := range transactions {
		switch t.Type {
		case "deposit":
			analytics.TotalIncome += t.Amount
			analytics.Categories[model.Income] += t.Amount
		case "withdrawal":
			analytics.TotalExpense += t.Amount
			analytics.Categories[model.Expense] += t.Amount
		case "transfer":
			if uint(t.FromAccountID) == request.AccountID {
				analytics.TotalExpense += t.Amount
				analytics.Categories[model.Transfer] += t.Amount
			} else {
				analytics.TotalIncome += t.Amount
				analytics.Categories[model.Transfer] += t.Amount
			}
		}
	}

	// Рассчитываем чистый доход
	analytics.NetIncome = analytics.TotalIncome - analytics.TotalExpense

	// Рассчитываем кредитную нагрузку
	if analytics.TotalIncome > 0 {
		analytics.CreditLoad = (creditPayments / analytics.TotalIncome) * 100
	}

	return analytics, nil
}

func (s *analyticsService) GetBalanceForecast(accountID uint, days int) ([]model.BalanceForecast, error) {
	// Получаем текущий баланс
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Получаем запланированные платежи
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, days)
	scheduledPayments, err := s.analyticsRepo.GetScheduledPayments(accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Создаем прогноз
	forecast := make([]model.BalanceForecast, 0)
	currentBalance := account.Balance

	// Добавляем текущий баланс
	forecast = append(forecast, model.BalanceForecast{
		Date:   startDate,
		Amount: currentBalance,
		Type:   "actual",
	})

	// Добавляем запланированные платежи
	for _, payment := range scheduledPayments {
		currentBalance -= payment.Amount
		forecast = append(forecast, model.BalanceForecast{
			Date:   payment.CreatedAt,
			Amount: currentBalance,
			Type:   "planned",
		})
	}

	return forecast, nil
} 