package repository

import (
	"context"
	"time"

	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

// AnalyticsRepository интерфейс репозитория аналитики
type AnalyticsRepository interface {
	GetDailyTransactions(ctx context.Context, date time.Time) ([]model.Transaction, error)
	GetMonthlyTransactions(ctx context.Context, year int, month time.Month) ([]model.Transaction, error)
	GetTransactionStats(ctx context.Context, startDate, endDate time.Time) (*model.TransactionStats, error)
	GetUserStats(ctx context.Context) (*model.UserStats, error)
	GetAccountStats(ctx context.Context) (*model.AccountStats, error)
	GetCreditStats(ctx context.Context) (*model.CreditStats, error)
	GetCardStats(ctx context.Context) (*model.CardStats, error)
	GetRoleStats(ctx context.Context) (*model.RoleStats, error)
}

// analyticsRepository реализация репозитория аналитики
type analyticsRepository struct {
	db *gorm.DB
}

// AnalyticsRepositoryInstance создает новый репозиторий аналитики
func AnalyticsRepositoryInstance(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

// GetDailyTransactions получает транзакции за день
func (r *analyticsRepository) GetDailyTransactions(ctx context.Context, date time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.db.Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetMonthlyTransactions получает транзакции за месяц
func (r *analyticsRepository) GetMonthlyTransactions(ctx context.Context, year int, month time.Month) ([]model.Transaction, error) {
	var transactions []model.Transaction
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	if err := r.db.Where("created_at BETWEEN ? AND ?", startOfMonth, endOfMonth).
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

// GetTransactionStats получает статистику по транзакциям
func (r *analyticsRepository) GetTransactionStats(ctx context.Context, startDate, endDate time.Time) (*model.TransactionStats, error) {
	var stats model.TransactionStats

	// Общее количество транзакций
	if err := r.db.Model(&model.Transaction{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.TotalTransactions).Error; err != nil {
		return nil, err
	}

	// Общая сумма транзакций
	if err := r.db.Model(&model.Transaction{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalAmount).Error; err != nil {
		return nil, err
	}

	// Средняя сумма транзакции
	if stats.TotalTransactions > 0 {
		stats.AverageAmount = stats.TotalAmount / float64(stats.TotalTransactions)
	}

	// Количество транзакций по типам
	if err := r.db.Model(&model.Transaction{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&stats.TransactionsByType).Error; err != nil {
		return nil, err
	}

	// Количество транзакций по статусам
	if err := r.db.Model(&model.Transaction{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&stats.TransactionsByStatus).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetUserStats получает статистику по пользователям
func (r *analyticsRepository) GetUserStats(ctx context.Context) (*model.UserStats, error) {
	var stats model.UserStats

	// Общее количество пользователей
	if err := r.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Количество активных пользователей
	if err := r.db.Model(&model.User{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Количество пользователей по ролям
	if err := r.db.Model(&model.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Select("roles.name, COUNT(*) as count").
		Group("roles.name").
		Scan(&stats.UsersByRole).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetAccountStats получает статистику по счетам
func (r *analyticsRepository) GetAccountStats(ctx context.Context) (*model.AccountStats, error) {
	var stats model.AccountStats

	// Общее количество счетов
	if err := r.db.Model(&model.Account{}).Count(&stats.TotalAccounts).Error; err != nil {
		return nil, err
	}

	// Общий баланс
	if err := r.db.Model(&model.Account{}).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&stats.TotalBalance).Error; err != nil {
		return nil, err
	}

	// Количество счетов по типам
	if err := r.db.Model(&model.Account{}).
		Select("type, COUNT(*) as count").
		Group("type").
		Scan(&stats.AccountsByType).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetCreditStats получает статистику по кредитам
func (r *analyticsRepository) GetCreditStats(ctx context.Context) (*model.CreditStats, error) {
	var stats model.CreditStats

	// Общее количество кредитов
	if err := r.db.Model(&model.Credit{}).Count(&stats.TotalCredits).Error; err != nil {
		return nil, err
	}

	// Общая сумма кредитов
	if err := r.db.Model(&model.Credit{}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalAmount).Error; err != nil {
		return nil, err
	}

	// Общая сумма выплат
	if err := r.db.Model(&model.Credit{}).
		Select("COALESCE(SUM(total_paid), 0)").
		Scan(&stats.TotalPaid).Error; err != nil {
		return nil, err
	}

	// Количество кредитов по статусам
	if err := r.db.Model(&model.Credit{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&stats.CreditsByStatus).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetCardStats получает статистику по картам
func (r *analyticsRepository) GetCardStats(ctx context.Context) (*model.CardStats, error) {
	var stats model.CardStats

	// Общее количество карт
	if err := r.db.Model(&model.Card{}).Count(&stats.TotalCards).Error; err != nil {
		return nil, err
	}

	// Количество активных карт
	if err := r.db.Model(&model.Card{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveCards).Error; err != nil {
		return nil, err
	}

	// Количество просроченных карт
	now := time.Now().Format("01/06")
	if err := r.db.Model(&model.Card{}).
		Where("expiry_date < ?", now).
		Count(&stats.ExpiredCards).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetRoleStats получает статистику по ролям
func (r *analyticsRepository) GetRoleStats(ctx context.Context) (*model.RoleStats, error) {
	var stats model.RoleStats

	// Общее количество ролей
	if err := r.db.Model(&model.Role{}).Count(&stats.TotalRoles).Error; err != nil {
		return nil, err
	}

	// Количество активных ролей
	if err := r.db.Model(&model.Role{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveRoles).Error; err != nil {
		return nil, err
	}

	// Количество пользователей по ролям
	if err := r.db.Model(&model.Role{}).
		Select("roles.name, COUNT(DISTINCT user_roles.user_id) as count").
		Joins("LEFT JOIN user_roles ON user_roles.role_id = roles.id").
		Group("roles.name").
		Scan(&stats.UsersByRole).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
