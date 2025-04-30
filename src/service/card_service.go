package service

import (
	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	"errors"
	"fmt"
	"regexp"
	"time"
)

type CardService interface {
	CreateCard(card *model.Card, userID uint) (*dto.UnsecureCard, error)
	GetCardByID(id uint) (*model.Card, error)
	GetUserCards(userID uint) ([]model.Card, error)
}

type cardService struct {
	cardRepo    repository.CardRepository
	accountRepo repository.AccountRepository
	publicKey   string
	hmacSecret  []byte
}

func CardServiceInstance(cardRepo repository.CardRepository, accountRepo repository.AccountRepository, publicKey string, hmacSecret []byte) CardService {
	return &cardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
	}
}

func (s *cardService) CreateCard(card *model.Card, userID uint) (*dto.UnsecureCard, error) {
	// Проверяем, что счет принадлежит пользователю
	accounts, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not get user accounts: %v", err)
	}

	accountExists := false
	for _, account := range accounts {
		if account.ID == card.AccountID {
			accountExists = true
			break
		}
	}

	if !accountExists {
		return nil, fmt.Errorf("account does not belong to the user")
	}

	var unsecureCard dto.UnsecureCard
	unsecureCard.Number = security.GenerateCardNumber("4", 16)

	// Проверка валидности номера карты
	if !security.IsValidCardNumber(unsecureCard.Number) {
		return nil, errors.New("invalid card number")
	}

	// Дополнительная валидация формата номера карты
	cardRegex := regexp.MustCompile(`^[0-9]{16}$`)
	if !cardRegex.MatchString(unsecureCard.Number) {
		return nil, errors.New("invalid card number format")
	}

	unsecureCard.CVV = security.GenerateCVV()

	// Валидация CVV
	cvvRegex := regexp.MustCompile(`^[0-9]{3}$`)
	if !cvvRegex.MatchString(unsecureCard.CVV) {
		return nil, errors.New("invalid CVV format")
	}

	unsecureCard.ExpiryDate = security.GenerateExpiryDate()

	// Валидация даты истечения срока действия
	expiryRegex := regexp.MustCompile(`^(0[1-9]|1[0-2])\/([0-9]{2})$`)
	if !expiryRegex.MatchString(unsecureCard.ExpiryDate) {
		return nil, errors.New("invalid expiry date format")
	}

	// Шифрование номера карты и срока действия
	encryptedNumber, err := security.EncryptData(unsecureCard.Number)
	if err != nil {
		return nil, err
	}
	encryptedExpiryDate, err := security.EncryptData(unsecureCard.ExpiryDate)
	if err != nil {
		return nil, err
	}

	// Хеширование CVV
	hashedCVV, err := security.HashCVV(card.CVV)
	if err != nil {
		return nil, err
	}

	// Генерация HMAC для данных карты
	hmacData := encryptedNumber + encryptedExpiryDate + hashedCVV
	card.HMAC = security.GenerateHMAC(hmacData, s.hmacSecret)

	// Сохранение зашифрованных данных в структуру
	card.Number = encryptedNumber
	card.ExpiryDate = encryptedExpiryDate
	card.CVV = hashedCVV
	card.CreatedAt = time.Now()

	// Сохранение карты в базе данных
	card, err = s.cardRepo.CreateCard(card, s.publicKey, s.hmacSecret)
	if err != nil {
		return nil, err
	}

	unsecureCard.ID = card.ID
	unsecureCard.AccountID = card.AccountID

	return &unsecureCard, nil
}

func (s *cardService) GetCardByID(id uint) (*model.Card, error) {
	card, err := s.cardRepo.GetCardByID(id)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (s *cardService) GetUserCards(userID uint) ([]model.Card, error) {
	// Получаем все счета пользователя
	accounts, err := s.accountRepo.GetAccountByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not get user accounts: %v", err)
	}

	// Собираем ID счетов пользователя
	accountIDs := make([]int, len(accounts))
	for i, account := range accounts {
		accountIDs[i] = account.ID
	}

	// Получаем карты только для счетов пользователя
	cards, err := s.cardRepo.GetCardsByAccountIDs(accountIDs)
	if err != nil {
		return nil, err
	}

	return cards, nil
}
