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

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"   // Ожидает обработки
	TransactionStatusCompleted TransactionStatus = "completed" // Выполнена
	TransactionStatusFailed    TransactionStatus = "failed"    // Ошибка
	TransactionStatusCancelled TransactionStatus = "cancelled" // Отменена
)

type Transaction struct {
	ID            int               `json:"id" gorm:"primaryKey"`
	Type          TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	FromAccountID int               `json:"from_account_id,omitempty" gorm:"index:idx_transaction_from_account"`
	ToAccountID   int               `json:"to_account_id" gorm:"index:idx_transaction_to_account;not null"`
	Amount        float64           `json:"amount" gorm:"not null;check:amount > 0"`
	Description   string            `json:"description"`
	Status        TransactionStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt     time.Time         `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time         `json:"updated_at" gorm:"not null"`

	FromAccount Account `gorm:"foreignKey:FromAccountID" json:"from_account,omitempty"`
	ToAccount   Account `gorm:"foreignKey:ToAccountID" json:"to_account"`
}

// ValidateAmount проверяет, что сумма транзакции положительная
func (t *Transaction) ValidateAmount() bool {
	return t.Amount > 0
}

// ValidateStatus проверяет корректность статуса
func (t *Transaction) ValidateStatus() bool {
	switch t.Status {
	case TransactionStatusPending,
		TransactionStatusCompleted,
		TransactionStatusFailed,
		TransactionStatusCancelled:
		return true
	default:
		return false
	}
}
