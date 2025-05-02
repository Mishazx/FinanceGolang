package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"context"
	"errors"
	"fmt"

	// "strconv"
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
	creditRepo      repository.CreditRepository
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	keyRateService  *ExternalService
}

func CreditServiceInstance(
	creditRepo repository.CreditRepository,
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	keyRateService *ExternalService,
) CreditService {
	return &creditService{
		creditRepo:      creditRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		keyRateService:  keyRateService,
	}
}

func (s *creditService) CreateCredit(userID uint, accountID uint, amount float64, termMonths int, description string) (*model.Credit, error) {
	// Проверяем, что счет принадлежит пользователю
	account, err := s.accountRepo.GetByID(context.Background(), accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}
	if account.UserID != userID {
		return nil, errors.New("account does not belong to the user")
	}

	// Получаем текущую ключевую ставку
	keyRate, err := s.keyRateService.GetKeyRate()
	if err != nil {
		return nil, fmt.Errorf("failed to get key rate: %v", err)
	}

	// Рассчитываем процентную ставку (ключевая ставка + 5%)
	interestRate := keyRate + 5.0

	// Создаем кредит
	credit := &model.Credit{
		AccountID:     accountID,
		Amount:        amount,
		Term:          termMonths,
		InterestRate:  interestRate,
		Status:        model.CreditStatusActive,
		StartDate:     time.Now(),
		EndDate:       time.Now().AddDate(0, termMonths, 0),
		PaymentDay:    time.Now().Day(),
		NextPayment:   time.Now().AddDate(0, 1, 0),
		RemainingDebt: amount,
	}

	// Сохраняем кредит
	if err := s.creditRepo.Create(context.Background(), credit); err != nil {
		return nil, fmt.Errorf("failed to create credit: %v", err)
	}

	// Зачисляем сумму кредита на счет пользователя
	account.Balance += amount
	if err := s.accountRepo.Update(context.Background(), account); err != nil {
		return nil, fmt.Errorf("failed to update account balance: %v", err)
	}

	// Создаем транзакцию о зачислении кредита
	transaction := &model.Transaction{
		Type:        model.TransactionTypeCredit,
		ToAccountID: accountID,
		Amount:      amount,
		Description: fmt.Sprintf("Зачисление по кредиту #%d: %s", credit.ID, description),
		Status:      model.TransactionStatusCompleted,
	}
	if err := s.transactionRepo.Create(context.Background(), transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	return credit, nil
}

func (s *creditService) GetCreditByID(id uint) (*model.Credit, error) {
	return s.creditRepo.GetByID(context.Background(), id)
}

func (s *creditService) GetUserCredits(userID uint) ([]model.Credit, error) {
	return s.creditRepo.GetCreditsByUserID(context.Background(), userID)
}

func (s *creditService) GetPaymentSchedule(creditID uint) ([]model.PaymentSchedule, error) {
	credit, err := s.creditRepo.GetByID(context.Background(), creditID)
	if err != nil {
		return nil, err
	}

	var schedule []model.PaymentSchedule
	remainingAmount := credit.Amount
	monthlyRate := credit.InterestRate / 12 / 100

	for i := 1; i <= credit.Term; i++ {
		interestAmount := remainingAmount * monthlyRate
		principalAmount := credit.CalculateMonthlyPayment() - interestAmount
		remainingAmount -= principalAmount

		schedule = append(schedule, model.PaymentSchedule{
			CreditID:      credit.ID,
			PaymentNumber: i,
			DueDate:       credit.StartDate.AddDate(0, i, 0),
			Amount:        credit.CalculateMonthlyPayment(),
			Interest:      interestAmount,
			Principal:     principalAmount,
			TotalAmount:   credit.CalculateMonthlyPayment(),
			Status:        model.PaymentStatusPending,
		})
	}

	return schedule, nil
}

func (s *creditService) ProcessPayment(creditID uint, paymentNumber int) error {
	// Получаем кредит
	credit, err := s.creditRepo.GetByID(context.Background(), creditID)
	if err != nil {
		return fmt.Errorf("failed to get credit: %v", err)
	}

	// Получаем график платежей
	schedule, err := s.GetPaymentSchedule(creditID)
	if err != nil {
		return fmt.Errorf("failed to get payment schedule: %v", err)
	}

	// Находим нужный платеж
	var payment *model.PaymentSchedule
	for i := range schedule {
		if schedule[i].PaymentNumber == paymentNumber {
			payment = &schedule[i]
			break
		}
	}
	if payment == nil {
		return errors.New("payment not found")
	}

	// Проверяем статус платежа
	if payment.Status != model.PaymentStatusPending {
		return errors.New("payment already processed")
	}

	// Проверяем баланс счета
	account, err := s.accountRepo.GetByID(context.Background(), credit.AccountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}

	if account.Balance < payment.TotalAmount {
		// Если средств недостаточно, начисляем штраф
		penalty := payment.TotalAmount * 0.1
		payment.TotalAmount += penalty
		credit.Status = model.CreditStatusOverdue
		if err := s.creditRepo.Update(context.Background(), credit); err != nil {
			return fmt.Errorf("failed to update credit status: %v", err)
		}
		return errors.New("insufficient funds, penalty applied")
	}

	// Списываем средства со счета
	account.Balance -= payment.TotalAmount
	if err := s.accountRepo.Update(context.Background(), account); err != nil {
		return fmt.Errorf("failed to update account balance: %v", err)
	}

	// Создаем транзакцию о платеже
	transaction := &model.Transaction{
		Type:          model.TransactionTypePayment,
		FromAccountID: credit.AccountID,
		Amount:        payment.TotalAmount,
		Description:   fmt.Sprintf("Платеж по кредиту #%d, платеж #%d", credit.ID, paymentNumber),
		Status:        model.TransactionStatusCompleted,
	}
	if err := s.transactionRepo.Create(context.Background(), transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Обновляем статус платежа
	payment.Status = model.PaymentStatusPaid
	now := time.Now()
	payment.PaidAt = &now

	// Обновляем общую сумму выплат
	if err := s.creditRepo.UpdateTotalPaid(context.Background(), credit.ID, payment.TotalAmount); err != nil {
		return fmt.Errorf("failed to update total paid: %v", err)
	}

	// Если это последний платеж, закрываем кредит
	if paymentNumber == credit.Term {
		credit.Status = model.CreditStatusPaid
		if err := s.creditRepo.Update(context.Background(), credit); err != nil {
			return fmt.Errorf("failed to update credit status: %v", err)
		}
	} else {
		// Обновляем дату следующего платежа
		nextPayment := credit.StartDate.AddDate(0, paymentNumber+1, 0)
		if err := s.creditRepo.UpdateNextPayment(context.Background(), credit.ID, nextPayment); err != nil {
			return fmt.Errorf("failed to update next payment date: %v", err)
		}
	}

	return nil
}

func (s *creditService) ProcessOverduePayments() error {
	// Получаем просроченные кредиты
	overdueCredits, err := s.creditRepo.GetOverdueCredits(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get overdue credits: %v", err)
	}

	for _, credit := range overdueCredits {
		// Получаем график платежей
		schedule, err := s.GetPaymentSchedule(credit.ID)
		if err != nil {
			return fmt.Errorf("failed to get payment schedule for credit %d: %v", credit.ID, err)
		}

		// Находим просроченный платеж
		for _, payment := range schedule {
			if payment.Status == model.PaymentStatusPending && payment.DueDate.Before(time.Now()) {
				// Пытаемся обработать платеж
				if err := s.ProcessPayment(credit.ID, payment.PaymentNumber); err != nil {
					// Если не удалось обработать платеж, обновляем статус кредита
					credit.Status = model.CreditStatusOverdue
					if err := s.creditRepo.Update(context.Background(), &credit); err != nil {
						return fmt.Errorf("failed to update credit status: %v", err)
					}
				}
				break
			}
		}
	}

	return nil
}
