package repository

import (
	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uint) (*model.Account, error)
	GetAllAccounts() ([]model.Account, error)
	UpdateAccount(account *model.Account) error
	DeleteAccount(id uint) error
	GetAccountByUserID(userID uint) ([]model.Account, error)
	GetAccountByName(name string) (*model.Account, error)
	GetAccountByType(accountType string) ([]model.Account, error)
	GetAccountByBalanceRange(minBalance, maxBalance float64) ([]model.Account, error)
}

type accountRepo struct {
	BaseRepository
}

func AccountRepositoryInstance(db *gorm.DB) AccountRepository {
	return &accountRepo{
		BaseRepository: InitializeRepository(db),
	}
}

func (r *accountRepo) CreateAccount(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepo) GetAccountByID(id uint) (*model.Account, error) {
	var account model.Account
	if err := r.db.Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *accountRepo) GetAllAccounts() ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (r *accountRepo) UpdateAccount(account *model.Account) error {
	return r.db.Save(account).Error
}
func (r *accountRepo) DeleteAccount(id uint) error {
	return r.db.Delete(&model.Account{}, id).Error
}
func (r *accountRepo) GetAccountByUserID(userID uint) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (r *accountRepo) GetAccountByName(name string) (*model.Account, error) {
	var account model.Account
	if err := r.db.Where("name = ?", name).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *accountRepo) GetAccountByType(accountType string) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("type = ?", accountType).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (r *accountRepo) GetAccountByBalanceRange(minBalance, maxBalance float64) ([]model.Account, error) {
	var accounts []model.Account
	if err := r.db.Where("balance >= ? AND balance <= ?", minBalance, maxBalance).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
