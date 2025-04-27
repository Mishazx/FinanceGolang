package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"errors"
	"fmt"
	"time"
)

type AccountOperationService interface {
	Deposit(accountID uint, amount float64, description string) error
	Withdraw(accountID uint, amount float64, description string) error
	Transfer(fromAccountID, toAccountID uint, amount float64, description string) error
	GetTransactions(accountID uint) ([]model.Transaction, error)
}

type accountOperationService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewAccountOperationService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) AccountOperationService {
	return &accountOperationService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *accountOperationService) Deposit(accountID uint, amount float64, description string) error {
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

func (s *accountOperationService) Withdraw(accountID uint, amount float64, description string) error {
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

func (s *accountOperationService) Transfer(fromAccountID, toAccountID uint, amount float64, description string) error {
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

func (s *accountOperationService) GetTransactions(accountID uint) ([]model.Transaction, error) {
	// Получаем транзакции за последние 30 дней
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	return s.transactionRepo.GetTransactionsByAccountID(accountID, startDate, endDate)
} 