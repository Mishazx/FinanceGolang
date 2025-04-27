package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

type CreditService interface {
	CreateCredit(userID uint, accountID uint, amount float64, termMonths int, description string) (*model.Credit, error)
	GetCreditByID(id uint) (*model.Credit, error)
	GetUserCredits(userID uint) ([]model.Credit, error)
	GetPaymentSchedule(creditID uint) ([]model.PaymentSchedule, error)
	ProcessPayment(creditID uint, paymentNumber int) error
	ProcessOverduePayments() error
}

type creditService struct {
	creditRepo     repository.CreditRepository
	accountRepo    repository.AccountRepository
	transactionRepo repository.TransactionRepository
	keyRateService CbrService
}

func NewCreditService(
	creditRepo repository.CreditRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	keyRateService CbrService,
) CreditService {
	return &creditService{
		creditRepo:     creditRepo,
		accountRepo:    accountRepo,
		transactionRepo: transactionRepo,
		keyRateService: keyRateService,
	}
}

func (s *creditService) CreateCredit(userID uint, accountID uint, amount float64, termMonths int, description string) (*model.Credit, error) {
	// Проверяем, что счет принадлежит пользователю
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}
	if uint(account.UserID) != userID {
		return nil, errors.New("account does not belong to the user")
	}

	// Получаем текущую ключевую ставку
	keyRateData, err := s.keyRateService.GetLastKeyRate()
	if err != nil {
		return nil, fmt.Errorf("failed to get key rate: %v", err)
	}

	// Преобразуем строковую ставку в число
	keyRate, err := strconv.ParseFloat(keyRateData.Rate, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key rate: %v", err)
	}

	// Рассчитываем процентную ставку (ключевая ставка + 5%)
	interestRate := keyRate + 5.0

	// Рассчитываем ежемесячный платеж по формуле аннуитета
	monthlyRate := interestRate / 12 / 100
	monthlyPayment := amount * (monthlyRate * math.Pow(1+monthlyRate, float64(termMonths))) / (math.Pow(1+monthlyRate, float64(termMonths)) - 1)

	// Создаем кредит
	credit := &model.Credit{
		UserID:         userID,
		AccountID:      accountID,
		Amount:         amount,
		InterestRate:   interestRate,
		TermMonths:     termMonths,
		MonthlyPayment: monthlyPayment,
		Status:         model.CreditStatusActive,
		StartDate:      time.Now(),
		EndDate:        time.Now().AddDate(0, termMonths, 0),
		Description:    description,
	}

	// Сохраняем кредит
	if err := s.creditRepo.CreateCredit(credit); err != nil {
		return nil, fmt.Errorf("failed to create credit: %v", err)
	}

	// Создаем график платежей
	if err := s.createPaymentSchedule(credit); err != nil {
		return nil, fmt.Errorf("failed to create payment schedule: %v", err)
	}

	return credit, nil
}

func (s *creditService) createPaymentSchedule(credit *model.Credit) error {
	remainingAmount := credit.Amount
	monthlyRate := credit.InterestRate / 12 / 100

	for i := 1; i <= credit.TermMonths; i++ {
		interestAmount := remainingAmount * monthlyRate
		principalAmount := credit.MonthlyPayment - interestAmount
		remainingAmount -= principalAmount

		schedule := &model.PaymentSchedule{
			CreditID:        credit.ID,
			PaymentNumber:   i,
			PaymentDate:     credit.StartDate.AddDate(0, i, 0),
			PrincipalAmount: principalAmount,
			InterestAmount:  interestAmount,
			TotalAmount:     credit.MonthlyPayment,
			Status:          "pending",
		}

		if err := s.creditRepo.CreatePaymentSchedule(schedule); err != nil {
			return err
		}
	}

	return nil
}

func (s *creditService) GetCreditByID(id uint) (*model.Credit, error) {
	return s.creditRepo.GetCreditByID(id)
}

func (s *creditService) GetUserCredits(userID uint) ([]model.Credit, error) {
	return s.creditRepo.GetCreditsByUserID(userID)
}

func (s *creditService) GetPaymentSchedule(creditID uint) ([]model.PaymentSchedule, error) {
	return s.creditRepo.GetPaymentScheduleByCreditID(creditID)
}

func (s *creditService) ProcessPayment(creditID uint, paymentNumber int) error {
	// Получаем график платежей
	schedules, err := s.creditRepo.GetPaymentScheduleByCreditID(creditID)
	if err != nil {
		return fmt.Errorf("failed to get payment schedule: %v", err)
	}

	// Находим нужный платеж
	var payment *model.PaymentSchedule
	for _, s := range schedules {
		if s.PaymentNumber == paymentNumber {
			payment = &s
			break
		}
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	// Проверяем статус платежа
	if payment.Status != "pending" {
		return errors.New("payment already processed")
	}

	// Получаем кредит
	credit, err := s.creditRepo.GetCreditByID(creditID)
	if err != nil {
		return fmt.Errorf("failed to get credit: %v", err)
	}

	// Проверяем баланс счета
	account, err := s.accountRepo.GetAccountByID(credit.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	if account.Balance < payment.TotalAmount {
		// Если средств недостаточно, начисляем штраф
		penalty := payment.TotalAmount * 0.1
		payment.TotalAmount += penalty
		credit.Status = model.CreditStatusOverdue
		if err := s.creditRepo.UpdateCredit(credit); err != nil {
			return fmt.Errorf("failed to update credit status: %v", err)
		}
		return errors.New("insufficient funds, penalty applied")
	}

	// Списываем средства
	account.Balance -= payment.TotalAmount
	if err := s.accountRepo.UpdateAccount(account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		Type:          model.TransactionTypeWithdrawal,
		FromAccountID: int(credit.AccountID),
		Amount:        payment.TotalAmount,
		Description:   fmt.Sprintf("Кредитный платеж #%d", paymentNumber),
		CreatedAt:     time.Now(),
		Status:        "completed",
	}
	if err := s.transactionRepo.CreateTransaction(transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Обновляем статус платежа
	now := time.Now()
	payment.Status = "paid"
	payment.PaidAt = &now
	if err := s.creditRepo.UpdatePaymentSchedule(payment); err != nil {
		return fmt.Errorf("failed to update payment schedule: %v", err)
	}

	// Проверяем, все ли платежи оплачены
	allPaid := true
	for _, s := range schedules {
		if s.Status != "paid" {
			allPaid = false
			break
		}
	}

	if allPaid {
		credit.Status = model.CreditStatusPaid
		if err := s.creditRepo.UpdateCredit(credit); err != nil {
			return fmt.Errorf("failed to update credit status: %v", err)
		}
	}

	return nil
}

func (s *creditService) ProcessOverduePayments() error {
	// Получаем просроченные платежи
	overduePayments, err := s.creditRepo.GetOverduePayments()
	if err != nil {
		return fmt.Errorf("failed to get overdue payments: %v", err)
	}

	for _, payment := range overduePayments {
		// Начисляем штраф
		penalty := payment.TotalAmount * 0.1
		payment.TotalAmount += penalty

		// Обновляем статус кредита
		credit, err := s.creditRepo.GetCreditByID(payment.CreditID)
		if err != nil {
			continue
		}
		credit.Status = model.CreditStatusOverdue
		if err := s.creditRepo.UpdateCredit(credit); err != nil {
			continue
		}

		// Обновляем платеж
		if err := s.creditRepo.UpdatePaymentSchedule(&payment); err != nil {
			continue
		}
	}

	return nil
} 