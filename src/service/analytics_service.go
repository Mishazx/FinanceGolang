package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"time"
)

type AnalyticsService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	creditRepo      repository.CreditRepository
}

func NewAnalyticsService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	creditRepo repository.CreditRepository,
) *AnalyticsService {
	return &AnalyticsService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		creditRepo:      creditRepo,
	}
}

// GetIncomeExpenseStats возвращает статистику доходов и расходов
func (s *AnalyticsService) GetIncomeExpenseStats(accountID uint, startDate, endDate time.Time) (*model.IncomeExpenseStats, error) {
	transactions, err := s.transactionRepo.GetTransactionsByAccountID(accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats := &model.IncomeExpenseStats{
		TotalIncome:  0,
		TotalExpense: 0,
		Categories:   make(map[string]float64),
	}

	for _, t := range transactions {
		if t.Type == model.TransactionTypeDeposit {
			stats.TotalIncome += t.Amount
		} else {
			stats.TotalExpense += t.Amount
		}
		stats.Categories[string(t.Type)] += t.Amount
	}

	return stats, nil
}

// GetBalanceForecast возвращает прогноз баланса на указанный период
func (s *AnalyticsService) GetBalanceForecast(accountID uint, months int) (*model.BalanceForecast, error) {
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	forecast := &model.BalanceForecast{
		CurrentBalance: account.Balance,
		MonthlyForecast: make([]model.MonthlyForecast, months),
	}

	now := time.Now()
	for i := 0; i < months; i++ {
		monthStart := time.Date(now.Year(), now.Month()+time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		// Получаем статистику за предыдущий месяц для прогноза
		stats, err := s.GetIncomeExpenseStats(accountID, monthStart.AddDate(0, -1, 0), monthEnd.AddDate(0, -1, 0))
		if err != nil {
			return nil, err
		}

		// Получаем предстоящие платежи по кредитам
		credits, err := s.creditRepo.GetCreditsByAccountID(accountID)
		if err != nil {
			return nil, err
		}

		var creditPayments float64
		for _, credit := range credits {
			if credit.Status == model.CreditStatusActive {
				schedule, err := s.creditRepo.GetPaymentSchedule(credit.ID)
				if err != nil {
					return nil, err
				}

				for _, payment := range schedule {
					if payment.DueDate.After(monthStart) && payment.DueDate.Before(monthEnd) {
						creditPayments += payment.TotalAmount
					}
				}
			}
		}

		forecast.MonthlyForecast[i] = model.MonthlyForecast{
			Month:     monthStart.Format("January 2006"),
			Income:    stats.TotalIncome,
			Expense:   stats.TotalExpense + creditPayments,
			Balance:   forecast.CurrentBalance + (stats.TotalIncome - stats.TotalExpense - creditPayments),
		}

		forecast.CurrentBalance = forecast.MonthlyForecast[i].Balance
	}

	return forecast, nil
}

// GetSpendingCategories возвращает статистику по категориям расходов
func (s *AnalyticsService) GetSpendingCategories(accountID uint, startDate, endDate time.Time) (map[string]float64, error) {
	transactions, err := s.transactionRepo.GetTransactionsByAccountID(accountID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	categories := make(map[string]float64)
	for _, t := range transactions {
		if t.Type == model.TransactionTypeWithdrawal || t.Type == model.TransactionTypeTransfer {
			categories[string(t.Type)] += t.Amount
		}
	}

	return categories, nil
} 