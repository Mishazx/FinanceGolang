package model

import (
	"time"
)

type CreditStatus string

const (
	CreditStatusActive    CreditStatus = "active"    // Активный кредит
	CreditStatusPaid      CreditStatus = "paid"      // Погашен
	CreditStatusOverdue   CreditStatus = "overdue"   // Просрочен
	CreditStatusCancelled CreditStatus = "cancelled" // Отменен
)

type Credit struct {
	ID              uint        `json:"id" gorm:"primaryKey"`
	UserID          uint        `json:"user_id"`
	AccountID       uint        `json:"account_id"`
	Amount          float64     `json:"amount"`
	InterestRate    float64     `json:"interest_rate"`
	TermMonths      int         `json:"term_months"`
	MonthlyPayment  float64     `json:"monthly_payment"`
	Status          CreditStatus `json:"status" gorm:"type:varchar(20)"`
	StartDate       time.Time   `json:"start_date"`
	EndDate         time.Time   `json:"end_date"`
	Description     string      `json:"description"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Account Account `gorm:"foreignKey:AccountID" json:"account"`
}

type PaymentSchedule struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreditID        uint      `json:"credit_id"`
	PaymentNumber   int       `json:"payment_number"`
	PaymentDate     time.Time `json:"payment_date"`
	PrincipalAmount float64   `json:"principal_amount"`
	InterestAmount  float64   `json:"interest_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status" gorm:"type:varchar(20)"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Credit Credit `gorm:"foreignKey:CreditID" json:"credit"`
} 