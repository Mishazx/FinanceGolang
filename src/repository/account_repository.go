package repository

import (
	"context"
	"time"

	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

// AccountRepository интерфейс репозитория счетов
type AccountRepository interface {
	Repository[model.Account]
	GetByNumber(ctx context.Context, number string) (*model.Account, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Account, error)
	GetWithTransactions(ctx context.Context, id uint) (*model.Account, error)
	UpdateBalance(ctx context.Context, id uint, amount float64) error
	GetByType(ctx context.Context, accountType model.AccountType) ([]model.Account, error)
	GetOverdueCredits(ctx context.Context) ([]model.Account, error)
	GetDailyTransactions(ctx context.Context, id uint, date time.Time) ([]model.Transaction, error)
	GetMonthlyTransactions(ctx context.Context, id uint, year int, month time.Month) ([]model.Transaction, error)
}

// accountRepository реализация репозитория счетов
type accountRepository struct {
	*BaseRepository[model.Account]
}

// AccountRepositoryInstance создает новый репозиторий счетов
func AccountRepositoryInstance(db *gorm.DB) AccountRepository {
	return &accountRepository{
		BaseRepository: NewBaseRepository[model.Account](db),
	}
}

// Create создает новый счет
func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := account.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Create(account).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// GetByID получает счет по ID
func (r *accountRepository) GetByID(ctx context.Context, id uint) (*model.Account, error) {
	var account model.Account
	if err := r.db.First(&account, id).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &account, nil
}

// GetByNumber получает счет по номеру
func (r *accountRepository) GetByNumber(ctx context.Context, number string) (*model.Account, error) {
	var account model.Account
	if err := r.db.Where("number = ?", number).First(&account).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &account, nil
}

// GetByUserID получает счета пользователя
func (r *accountRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return accounts, nil
}

// GetWithTransactions получает счет с транзакциями
func (r *accountRepository) GetWithTransactions(ctx context.Context, id uint) (*model.Account, error) {
	var account model.Account
	if err := r.db.Preload("Transactions").First(&account, id).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &account, nil
}

// Update обновляет счет
func (r *accountRepository) Update(ctx context.Context, account *model.Account) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := account.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Save(account).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// UpdateBalance обновляет баланс счета
func (r *accountRepository) UpdateBalance(ctx context.Context, id uint, amount float64) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Account{}).Where("id = ?", id).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// Delete удаляет счет
func (r *accountRepository) Delete(ctx context.Context, id uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Account{}, id).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// List получает список счетов
func (r *accountRepository) List(ctx context.Context, offset, limit int) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return accounts, nil
}

// Count возвращает количество счетов
func (r *accountRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Account{}).Count(&count).Error; err != nil {
		return 0, r.HandleError(err)
	}
	return count, nil
}

// GetByType получает счета по типу
func (r *accountRepository) GetByType(ctx context.Context, accountType model.AccountType) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("type = ?", accountType).Find(&accounts).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return accounts, nil
}

// GetOverdueCredits получает просроченные кредитные счета
func (r *accountRepository) GetOverdueCredits(ctx context.Context) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("type = ? AND balance < 0", model.AccountTypeCredit).Find(&accounts).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return accounts, nil
}

// GetDailyTransactions получает транзакции за день
func (r *accountRepository) GetDailyTransactions(ctx context.Context, id uint, date time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	if err := r.db.Where("(from_account_id = ? OR to_account_id = ?) AND created_at BETWEEN ? AND ?",
		id, id, startOfDay, endOfDay).Find(&transactions).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return transactions, nil
}

// GetMonthlyTransactions получает транзакции за месяц
func (r *accountRepository) GetMonthlyTransactions(ctx context.Context, id uint, year int, month time.Month) ([]model.Transaction, error) {
	var transactions []model.Transaction
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	if err := r.db.Where("(from_account_id = ? OR to_account_id = ?) AND created_at BETWEEN ? AND ?",
		id, id, startOfMonth, endOfMonth).Find(&transactions).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return transactions, nil
}
