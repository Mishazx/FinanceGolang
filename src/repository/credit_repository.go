package repository

import (
	"context"
	"time"

	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

// CreditRepository интерфейс репозитория кредитов
type CreditRepository interface {
	Repository[model.Credit]
	GetByAccountID(ctx context.Context, accountID uint) (*model.Credit, error)
	GetActiveCredits(ctx context.Context) ([]model.Credit, error)
	GetOverdueCredits(ctx context.Context) ([]model.Credit, error)
	GetCreditsByUserID(ctx context.Context, userID uint) ([]model.Credit, error)
	UpdateStatus(ctx context.Context, id uint, status model.CreditStatus) error
	UpdateNextPayment(ctx context.Context, id uint, nextPayment time.Time) error
	UpdateTotalPaid(ctx context.Context, id uint, amount float64) error
	GetCreditsByStatus(ctx context.Context, status model.CreditStatus) ([]model.Credit, error)
	GetCreditsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]model.Credit, error)
	GetPaymentSchedule(ctx context.Context, creditID uint) ([]model.PaymentSchedule, error)
	UpdatePaymentSchedule(ctx context.Context, payment *model.PaymentSchedule) error
}

// creditRepository реализация репозитория кредитов
type creditRepository struct {
	BaseRepository[model.Credit]
}

// CreditRepositoryInstance создает новый репозиторий кредитов
func CreditRepositoryInstance(db *gorm.DB) CreditRepository {
	return &creditRepository{
		BaseRepository: *NewBaseRepository[model.Credit](db),
	}
}

// Create создает новый кредит
func (r *creditRepository) Create(ctx context.Context, credit *model.Credit) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := credit.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Create(credit).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// GetByID получает кредит по ID
func (r *creditRepository) GetByID(ctx context.Context, id uint) (*model.Credit, error) {
	var credit model.Credit
	if err := r.db.First(&credit, id).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &credit, nil
}

// GetByAccountID получает кредит по ID счета
func (r *creditRepository) GetByAccountID(ctx context.Context, accountID uint) (*model.Credit, error) {
	var credit model.Credit
	if err := r.db.Where("account_id = ?", accountID).First(&credit).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &credit, nil
}

// GetActiveCredits получает активные кредиты
func (r *creditRepository) GetActiveCredits(ctx context.Context) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Where("status = ?", model.CreditStatusActive).Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// GetOverdueCredits получает просроченные кредиты
func (r *creditRepository) GetOverdueCredits(ctx context.Context) ([]model.Credit, error) {
	var credits []model.Credit
	now := time.Now()
	if err := r.db.Where("status = ? AND next_payment < ?", model.CreditStatusActive, now).Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// GetCreditsByUserID получает кредиты пользователя
func (r *creditRepository) GetCreditsByUserID(ctx context.Context, userID uint) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Joins("JOIN accounts ON accounts.id = credits.account_id").
		Where("accounts.user_id = ?", userID).
		Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// Update обновляет кредит
func (r *creditRepository) Update(ctx context.Context, credit *model.Credit) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := credit.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Save(credit).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// UpdateStatus обновляет статус кредита
func (r *creditRepository) UpdateStatus(ctx context.Context, id uint, status model.CreditStatus) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Credit{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// UpdateNextPayment обновляет дату следующего платежа
func (r *creditRepository) UpdateNextPayment(ctx context.Context, id uint, nextPayment time.Time) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Credit{}).Where("id = ?", id).Update("next_payment", nextPayment).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// UpdateTotalPaid обновляет общую сумму выплат
func (r *creditRepository) UpdateTotalPaid(ctx context.Context, id uint, amount float64) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Credit{}).Where("id = ?", id).
			Update("total_paid", gorm.Expr("total_paid + ?", amount)).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// Delete удаляет кредит
func (r *creditRepository) Delete(ctx context.Context, id uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Credit{}, id).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// List получает список кредитов
func (r *creditRepository) List(ctx context.Context, offset, limit int) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Offset(offset).Limit(limit).Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// Count возвращает количество кредитов
func (r *creditRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Credit{}).Count(&count).Error; err != nil {
		return 0, r.HandleError(err)
	}
	return count, nil
}

// GetCreditsByStatus получает кредиты по статусу
func (r *creditRepository) GetCreditsByStatus(ctx context.Context, status model.CreditStatus) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Where("status = ?", status).Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// GetCreditsByDateRange получает кредиты в указанном диапазоне дат
func (r *creditRepository) GetCreditsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&credits).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return credits, nil
}

// GetPaymentSchedule получает график платежей по кредиту
func (r *creditRepository) GetPaymentSchedule(ctx context.Context, creditID uint) ([]model.PaymentSchedule, error) {
	var schedule []model.PaymentSchedule
	if err := r.db.Where("credit_id = ?", creditID).Find(&schedule).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return schedule, nil
}

// UpdatePaymentSchedule обновляет график платежей
func (r *creditRepository) UpdatePaymentSchedule(ctx context.Context, payment *model.PaymentSchedule) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Save(payment).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}
