package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	"errors"
	"fmt"
	"time"
)

type CardService interface {
	CreateCard(card *model.Card) error
	GetCardByID(id uint) (*model.Card, error)
	GetAllCards() ([]model.Card, error)
}

type cardService struct {
	cardRepo   repository.CardRepository
	publicKey  string
	hmacSecret []byte
}

func NewCardService(cardRepo repository.CardRepository, publicKey string, hmacSecret []byte) CardService {
	return &cardService{
		cardRepo:   cardRepo,
		publicKey:  publicKey,
		hmacSecret: hmacSecret,
	}
}

func (s *cardService) CreateCard(card *model.Card) error {
	fmt.Printf("checking card number: %s\n", card.Number)
	// Проверка валидности номера карты
	if !security.IsValidCardNumber(card.Number) {
		return errors.New("invalid card number")
	}

	fmt.Printf("encrypting card CVV: %s\n", card.CVV)

	// Шифрование номера карты и срока действия
	encryptedNumber, err := security.EncryptData(card.Number, s.publicKey)
	if err != nil {
		return err
	}
	encryptedExpiryDate, err := security.EncryptData(card.ExpiryDate, s.publicKey)
	if err != nil {
		return err
	}

	fmt.Printf("hashing card CVV: %s\n", card.CVV)

	// Хеширование CVV
	hashedCVV, err := security.HashCVV(card.CVV)
	if err != nil {
		return err
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
	if err := s.cardRepo.CreateCard(card, s.publicKey, s.hmacSecret); err != nil {
		return err
	}

	return nil
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
