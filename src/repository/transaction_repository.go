package repository

import (
	"FinanceGolang/src/model"
	// "errors"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransactionsByAccountID(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error)
	GetTransactionByID(id uint) (*model.Transaction, error)
	GetTransactionsByUserID(userID uint) ([]model.Transaction, error)
}

type transactionRepo struct {
	BaseRepository
}

func TransactionRepositoryInstance(db *gorm.DB) TransactionRepository {
	return &transactionRepo{
		BaseRepository: InitializeRepository(db),
	}
}

func (r *transactionRepo) CreateTransaction(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepo) GetTransactionsByAccountID(accountID uint, startDate, endDate time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("(from_account_id = ? OR to_account_id = ?) AND created_at BETWEEN ? AND ?",
		accountID, accountID, startDate, endDate).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepo) GetTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

func (r *transactionRepo) GetTransactionsByUserID(userID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Joins("JOIN accounts ON accounts.id = transactions.to_account_id OR accounts.id = transactions.from_account_id").
		Where("accounts.user_id = ?", userID).
		Order("transactions.created_at desc").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
