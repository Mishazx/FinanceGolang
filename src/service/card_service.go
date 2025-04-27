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
	CreateCard(card *model.Card) (*dto.UnsecureCard, error) //error
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

func (s *cardService) CreateCard(card *model.Card) (*dto.UnsecureCard, error) {
	// GenerateCardNumber("4", 16),  // Visa
	// GenerateCardNumber("5", 16),  // MasterCard
	// GenerateCardNumber("37", 15), // American Express
	// GenerateCardNumber("6", 16),  // Discover

	// fmt.Println("CARD : ", card)
	// fmt.Println("Card ID: ", card.ID)
	// fmt.Println("Card Account: ", card.Account)
	// fmt.Println("Card Account ID: ", card.AccountID)

	// # Проверка возможности привязки карты к счету

	var unsecureCard dto.UnsecureCard
	unsecureCard.Number = security.GenerateCardNumber("4", 16)

	// fmt.Println("Generated card number:", unsecureCard.Number)

	// Проверка валидности номера карты
	if !security.IsValidCardNumber(unsecureCard.Number) {
		return nil, errors.New("invalid card number")
	}

	unsecureCard.CVV = security.GenerateCVV()
	unsecureCard.ExpiryDate = security.GenerateExpiryDate()

	// fmt.Printf("encrypting card CVV: %s\n", card.CVV)
	// fmt.Printf(("encrypting card expiry date: %s\n"), card.ExpiryDate)

	// Шифрование номера карты и срока действия
	encryptedNumber, err := security.EncryptData(unsecureCard.Number)
	if err != nil {
		return nil, err
	}
	encryptedExpiryDate, err := security.EncryptData(unsecureCard.ExpiryDate)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Encrypted Number:", encryptedNumber)
	// fmt.Println("Encrypted Expiry Date:", encryptedExpiryDate)
	// fmt.Printf("hashing card CVV: %s\n", card.CVV)

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

	fmt.Println("Card created successfully: ", unsecureCard)

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
