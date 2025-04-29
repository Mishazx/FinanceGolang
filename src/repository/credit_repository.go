package repository

import (
	"FinanceGolang/src/model"
	"time"

	"gorm.io/gorm"
)

type CreditRepository interface {
	CreateCredit(credit *model.Credit) error
	GetCreditByID(id uint) (*model.Credit, error)
	GetCreditsByUserID(userID uint) ([]model.Credit, error)
	GetCreditsByAccountID(accountID uint) ([]model.Credit, error)
	GetAllCredits() ([]model.Credit, error)
	UpdateCredit(credit *model.Credit) error
	CreatePaymentSchedule(schedule *model.PaymentSchedule) error
	GetPaymentSchedule(creditID uint) ([]model.PaymentSchedule, error)
	UpdatePaymentSchedule(schedule *model.PaymentSchedule) error
	GetOverduePayments() ([]model.PaymentSchedule, error)
	UpdatePaymentStatus(paymentID uint, status string) error
}

type creditRepo struct {
	BaseRepository
}

func CreditRepositoryInstance(db *gorm.DB) CreditRepository {
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

func (r *creditRepo) GetCreditsByAccountID(accountID uint) ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Where("account_id = ?", accountID).Find(&credits).Error; err != nil {
		return nil, err
	}
	return credits, nil
}

func (r *creditRepo) GetAllCredits() ([]model.Credit, error) {
	var credits []model.Credit
	if err := r.db.Find(&credits).Error; err != nil {
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

func (r *creditRepo) GetPaymentSchedule(creditID uint) ([]model.PaymentSchedule, error) {
	var schedules []model.PaymentSchedule
	if err := r.db.Where("credit_id = ?", creditID).Find(&schedules).Error; err != nil {
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
	if err := r.db.Where("status = ? AND due_date < ?", model.PaymentStatusPending, now).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *creditRepo) UpdatePaymentStatus(paymentID uint, status string) error {
	return r.db.Model(&model.PaymentSchedule{}).Where("id = ?", paymentID).Update("status", status).Error
}
