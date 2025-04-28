package model

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"    // Пополнение
	TransactionTypeWithdrawal TransactionType = "withdrawal" // Списание
	TransactionTypeTransfer   TransactionType = "transfer"   // Перевод
	TransactionTypeCredit     TransactionType = "credit"     // Оформление кредита
)

type Transaction struct {
	ID            int             `json:"id" gorm:"primaryKey"`
	Type          TransactionType `json:"type" gorm:"type:varchar(20)"`
	FromAccountID int             `json:"from_account_id,omitempty"`
	ToAccountID   int             `json:"to_account_id"`
	Amount        float64         `json:"amount"`
	Description   string          `json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
	Status        string          `json:"status" gorm:"type:varchar(20)"`

	FromAccount Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID" json:"to_account"`
}
