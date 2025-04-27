package repository

import (
	"FinanceGolang/src/model"
	"gorm.io/gorm"
	"time"
)

type CreditRepository interface {
	CreateCredit(credit *model.Credit) error
	GetCreditByID(id uint) (*model.Credit, error)
	GetCreditsByUserID(userID uint) ([]model.Credit, error)
	UpdateCredit(credit *model.Credit) error
	CreatePaymentSchedule(schedule *model.PaymentSchedule) error
	GetPaymentScheduleByCreditID(creditID uint) ([]model.PaymentSchedule, error)
	UpdatePaymentSchedule(schedule *model.PaymentSchedule) error
	GetOverduePayments() ([]model.PaymentSchedule, error)
}

type creditRepo struct {
	BaseRepository
}

func NewCreditRepository(db *gorm.DB) CreditRepository {
	return &creditRepo{
		BaseRepository: InitializeRepository(db),
	}
}

func (r *creditRepo) CreateCredit(credit *model.Credit) error {
	return r.db.Create(credit).Error
}

func (r *creditRepo) GetCreditByID(id uint) (*model.Credit, error) {
	var credit model.Credit
	if err := r.db.First(&credit, id).Error; err != nil {
		return nil, err
	}
	return &credit, nil
}

func (r *creditRepo) GetCreditsByUserID(userID uint) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Where("user_id = ?", userID).Find(&credits).Error; err != nil {
		return nil, err
	}
	return credits, nil
}

func (r *creditRepo) UpdateCredit(credit *model.Credit) error {
	return r.db.Save(credit).Error
}

func (r *creditRepo) CreatePaymentSchedule(schedule *model.PaymentSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *creditRepo) GetPaymentScheduleByCreditID(creditID uint) ([]model.PaymentSchedule, error) {
	var schedules []model.PaymentSchedule
	if err := r.db.Where("credit_id = ?", creditID).Order("payment_number").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *creditRepo) UpdatePaymentSchedule(schedule *model.PaymentSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *creditRepo) GetOverduePayments() ([]model.PaymentSchedule, error) {
	var schedules []model.PaymentSchedule
	now := time.Now()
	if err := r.db.Where("status = ? AND payment_date < ?", "pending", now).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
} 