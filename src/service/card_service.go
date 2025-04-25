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
	// GenerateCardNumber("4", 16),  // Visa
	// GenerateCardNumber("5", 16),  // MasterCard
	// GenerateCardNumber("37", 15), // American Express
	// GenerateCardNumber("6", 16),  // Discover
	card.Number = security.GenerateCardNumber("4", 16)
	fmt.Println("Generated card number:", card.Number)

	// Проверка валидности номера карты
	if !security.IsValidCardNumber(card.Number) {
		return errors.New("invalid card number")
	}

	card.CVV = security.GenerateCVV()
	card.ExpiryDate = security.GenerateExpiryDate()

	fmt.Printf("encrypting card CVV: %s\n", card.CVV)
	fmt.Printf(("encrypting card expiry date: %s\n"), card.ExpiryDate)

	// Шифрование номера карты и срока действия
	encryptedNumber, err := security.EncryptData(card.Number)
	if err != nil {
		return err
	}
	fmt.Println("Encrypted Number:", encryptedNumber)
	encryptedExpiryDate, err := security.EncryptData(card.ExpiryDate)
	if err != nil {
		return err
	}
	fmt.Println("Encrypted Expiry Date:", encryptedExpiryDate)

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
