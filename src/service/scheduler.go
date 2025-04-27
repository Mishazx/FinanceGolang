package service

import (
    "fmt"
    "time"

    "FinanceGolang/src/model"
    "FinanceGolang/src/repository"
)

type Scheduler struct {
    creditRepo repository.CreditRepository
    accountRepo repository.AccountRepository
    transactionRepo repository.TransactionRepository
    externalService *ExternalService
}

func NewScheduler(
    creditRepo repository.CreditRepository,
    accountRepo repository.AccountRepository,
    transactionRepo repository.TransactionRepository,
    externalService *ExternalService,
) *Scheduler {
    return &Scheduler{
        creditRepo: creditRepo,
        accountRepo: accountRepo,
        transactionRepo: transactionRepo,
        externalService: externalService,
    }
}

// Start запускает шедулер
func (s *Scheduler) Start() {
    // Проверка платежей каждые 12 часов
    go s.checkPayments()
}

// checkPayments проверяет и обрабатывает платежи по кредитам
func (s *Scheduler) checkPayments() {
    ticker := time.NewTicker(12 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        credits, err := s.creditRepo.GetAllCredits()
        if err != nil {
            fmt.Printf("Ошибка при получении кредитов: %v\n", err)
            continue
        }

        for _, credit := range credits {
            if credit.Status == model.CreditStatusActive {
                s.processCreditPayment(&credit)
            }
        }
    }
}

// processCreditPayment обрабатывает платеж по кредиту
func (s *Scheduler) processCreditPayment(credit *model.Credit) {
    // Получаем график платежей
    schedule, err := s.creditRepo.GetPaymentSchedule(credit.ID)
    if err != nil {
        fmt.Printf("Ошибка при получении графика платежей: %v\n", err)
        return
    }

    now := time.Now()
    for _, payment := range schedule {
        // Если платеж просрочен и не оплачен
        if payment.DueDate.Before(now) && payment.Status == model.PaymentStatusPending {
            // Проверяем баланс счета
            account, err := s.accountRepo.GetAccountByID(credit.AccountID)
            if err != nil {
                fmt.Printf("Ошибка при получении счета: %v\n", err)
                continue
            }

            // Если на счету достаточно средств
            if account.Balance >= payment.TotalAmount {
                // Списание средств
                err = s.accountRepo.UpdateBalance(uint(account.ID), -payment.TotalAmount)
                if err != nil {
                    fmt.Printf("Ошибка при списании средств: %v\n", err)
                    continue
                }

                // Обновление статуса платежа
                payment.Status = model.PaymentStatusPaid
                err = s.creditRepo.UpdatePaymentStatus(payment.ID, string(payment.Status))
                if err != nil {
                    fmt.Printf("Ошибка при обновлении статуса платежа: %v\n", err)
                    continue
                }

                // Отправка уведомления
                user, err := s.accountRepo.GetUserByAccountID(uint(account.ID))
                if err == nil {
                    s.externalService.SendPaymentNotification(
                        user.Email,
                        "Платеж по кредиту",
                        payment.TotalAmount,
                    )
                }
            } else {
                // Начисление штрафа за просрочку
                penalty := payment.TotalAmount * 0.1
                payment.TotalAmount += penalty
                payment.Status = model.PaymentStatusOverdue
                
                err = s.creditRepo.UpdatePaymentStatus(payment.ID, string(payment.Status))
                if err != nil {
                    fmt.Printf("Ошибка при обновлении статуса платежа: %v\n", err)
                    continue
                }

                // Отправка уведомления о просрочке
                user, err := s.accountRepo.GetUserByAccountID(uint(account.ID))
                if err == nil {
                    s.externalService.SendPaymentNotification(
                        user.Email,
                        "Просрочка платежа по кредиту",
                        payment.TotalAmount,
                    )
                }
            }
        }
    }
} 