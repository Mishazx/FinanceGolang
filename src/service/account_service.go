package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"fmt"
)

type AccountService interface {
	CreateAccount(account *model.Account) error
	GetAccountByID(id uint) (*model.Account, error)
	GetAllAccounts() ([]model.Account, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo: accountRepo}
}

func (s *accountService) CreateAccount(account *model.Account) error {
	if err := s.accountRepo.CreateAccount(account); err != nil {
		return fmt.Errorf("could not create account: %v", err)
	}
	return nil
}

func (s *accountService) GetAccountByID(id uint) (*model.Account, error) {
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get account by ID: %v", err)
	}
	return account, nil
}

func (s *accountService) GetAllAccounts() ([]model.Account, error) {
	accounts, err := s.accountRepo.GetAllAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not get all accounts: %v", err)
	}
	return accounts, nil
}
