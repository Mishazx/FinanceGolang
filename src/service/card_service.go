package service

import (
	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	"errors"
	"fmt"
	"time"
)

type CardService interface {
	CreateCard(card *model.Card, userID uint) (*dto.UnsecureCard, error)
	GetCardByID(id uint) (*model.Card, error)
	GetAllCards() ([]model.Card, error)
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

	unsecureCard.CVV = security.GenerateCVV()
	unsecureCard.ExpiryDate = security.GenerateExpiryDate()

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

func (s *cardService) GetAllCards() ([]model.Card, error) {
	cards, err := s.cardRepo.GetAllCards()
	if err != nil {
		return nil, err
	}
	return cards, nil
}
