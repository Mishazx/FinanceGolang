package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"fmt"
)

type AccountService interface {
	CreateAccount(account *model.Account, userID uint) error
	GetAccountByID(id uint) (*model.Account, error)
	GetAccountByUserID(id uint) ([]model.Account, error)
	GetAllAccounts() ([]model.Account, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo: accountRepo}
}

func (s *accountService) CreateAccount(account *model.Account, userID uint) error {
	fmt.Println("Creating account for user ID:", userID)
	if err := s.accountRepo.CreateAccount(account, userID); err != nil {
		return fmt.Errorf("could not create account: %v", err)
	}
	return nil
}

func (s *accountService) GetAccountByID(id uint) (*model.Account, error) {
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not get account by ID: %v", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found with ID: %d", id)
	}
	return account, nil
}

func (s *accountService) GetAccountByUserID(userID uint) ([]model.Account, error) {
	accounts, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not get account by user ID: %v", err)
	}
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found for user ID: %d", userID)
	}
	return accounts, nil
}

func (s *accountService) GetAllAccounts() ([]model.Account, error) {
	accounts, err := s.accountRepo.GetAllAccounts()
	if err != nil {
		return nil, fmt.Errorf("could not get all accounts: %v", err)
	}
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts found")
	}
	return accounts, nil
}
