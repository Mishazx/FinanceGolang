package model

import (
	"time"
)

type Account struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id" gorm:"index:idx_account_user_id"`
	Name      string    `json:"name" gorm:"not null"`
	Balance   float64   `json:"balance" gorm:"check:balance >= 0"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	// ExpiredAt time.Time `json:"expired_at"`

	Cards []Card `gorm:"foreignKey:AccountID" json:"cards"`
}

// ValidateBalance проверяет, что баланс не отрицательный
func (a *Account) ValidateBalance() bool {
	return a.Balance >= 0
}
