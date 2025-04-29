package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"errors"
	"fmt"
	"time"
)

type AccountService interface {
	// Базовые операции со счетом
	CreateAccount(account *model.Account, userID uint) error
	GetAccountByID(id uint) (*model.Account, error)
	GetAccountByUserID(id uint) ([]model.Account, error)
	GetAllAccounts() ([]model.Account, error)

	// Операции с балансом
	Deposit(accountID uint, amount float64, description string) error
	Withdraw(accountID uint, amount float64, description string) error
	Transfer(fromAccountID, toAccountID uint, amount float64, description string) error

	// Операции с транзакциями
	GetTransactions(accountID uint) ([]model.Transaction, error)
}

type accountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func AccountServiceInstance(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) AccountService {
	return &accountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// Базовые операции со счетом
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

// Операции с балансом
func (s *accountService) Deposit(accountID uint, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		Type:        model.TransactionTypeDeposit,
		ToAccountID: int(accountID),
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
		Status:      "completed",
	}

	// Обновляем баланс счета
	account.Balance += amount
	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	// Сохраняем транзакцию
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	return nil
}

func (s *accountService) Withdraw(accountID uint, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		Type:          model.TransactionTypeWithdrawal,
		FromAccountID: int(accountID),
		Amount:        amount,
		Description:   description,
		CreatedAt:     time.Now(),
		Status:        "completed",
	}

	// Обновляем баланс счета
	account.Balance -= amount
	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	// Сохраняем транзакцию
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	return nil
}

func (s *accountService) Transfer(fromAccountID, toAccountID uint, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to the same account")
	}

	fromAccount, err := s.accountRepo.GetAccountByID(fromAccountID)
	if err != nil {
		return fmt.Errorf("failed to get source account: %v", err)
	}

	toAccount, err := s.accountRepo.GetAccountByID(toAccountID)
	if err != nil {
		return fmt.Errorf("failed to get destination account: %v", err)
	}

	if fromAccount.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		Type:          model.TransactionTypeTransfer,
		FromAccountID: int(fromAccountID),
		ToAccountID:   int(toAccountID),
		Amount:        amount,
		Description:   description,
		CreatedAt:     time.Now(),
		Status:        "completed",
	}

	// Обновляем балансы счетов
	fromAccount.Balance -= amount
	toAccount.Balance += amount

	if err := s.accountRepo.UpdateAccount(fromAccount); err != nil {
		return fmt.Errorf("failed to update source account balance: %v", err)
	}

	if err := s.accountRepo.UpdateAccount(toAccount); err != nil {
		return fmt.Errorf("failed to update destination account balance: %v", err)
	}

	// Сохраняем транзакцию
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	return nil
}

// Операции с транзакциями
func (s *accountService) GetTransactions(accountID uint) ([]model.Transaction, error) {
	// Получаем транзакции за последние 30 дней
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	return s.transactionRepo.GetTransactionsByAccountID(accountID, startDate, endDate)
}
