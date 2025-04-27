package repository

import (
	"FinanceGolang/src/model"
	// "errors"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransactionByID(id uint) (*model.Transaction, error)
	GetTransactionsByAccountID(accountID uint) ([]model.Transaction, error)
	GetTransactionsByUserID(userID uint) ([]model.Transaction, error)
}

type transactionRepo struct {
	BaseRepository
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepo{
		BaseRepository: InitializeRepository(db),
	}
}

func (r *transactionRepo) CreateTransaction(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepo) GetTransactionByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepo) GetTransactionsByAccountID(accountID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
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