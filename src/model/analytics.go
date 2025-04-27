package model

import "time"

type AnalyticsPeriod string

const (
	Daily   AnalyticsPeriod = "daily"
	Weekly  AnalyticsPeriod = "weekly"
	Monthly AnalyticsPeriod = "monthly"
	Yearly  AnalyticsPeriod = "yearly"
)

type TransactionCategory string

const (
	Income  TransactionCategory = "income"
	Expense TransactionCategory = "expense"
	Transfer TransactionCategory = "transfer"
)

type Analytics struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	AccountID uint      `json:"account_id" gorm:"index"`
	Period    AnalyticsPeriod `json:"period"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	
	// Общая статистика
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetIncome    float64 `json:"net_income"`
	
	// Статистика по категориям
	Categories map[TransactionCategory]float64 `json:"categories" gorm:"-"`
	
	// Кредитная нагрузка
	CreditPayments float64 `json:"credit_payments"`
	CreditLoad     float64 `json:"credit_load"` // Процент от дохода
	
	// Прогноз
	BalanceForecast []BalanceForecast `json:"balance_forecast" gorm:"-"`
}

type BalanceForecast struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
	Type   string    `json:"type"` // planned, actual
}

type AnalyticsRequest struct {
	AccountID uint           `json:"account_id" binding:"required"`
	Period    AnalyticsPeriod `json:"period" binding:"required"`
	StartDate time.Time      `json:"start_date" binding:"required"`
	EndDate   time.Time      `json:"end_date" binding:"required"`
} 