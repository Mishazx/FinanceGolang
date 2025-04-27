package repository

import (
	"FinanceGolang/src/model"
	"time"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetTransactionsByPeriod(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error)
	GetCreditPaymentsByPeriod(accountID uint, startDate, endDate time.Time) (float64, error)
	GetScheduledPayments(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error)
}

type analyticsRepo struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepo{db: db}
}

func (r *analyticsRepo) GetTransactionsByPeriod(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("(from_account_id = ? OR to_account_id = ?) AND created_at BETWEEN ? AND ?", 
		accountID, accountID, startDate, endDate).
		Order("created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *analyticsRepo) GetCreditPaymentsByPeriod(accountID uint, startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&model.Transaction{}).
		Where("from_account_id = ? AND type = ? AND created_at BETWEEN ? AND ?", 
			accountID, "credit_payment", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *analyticsRepo) GetScheduledPayments(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("from_account_id = ? AND status = ? AND created_at BETWEEN ? AND ?", 
		accountID, "scheduled", startDate, endDate).
		Order("created_at ASC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
} 