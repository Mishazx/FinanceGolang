package repository

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/security"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CardRepository interface {
	CreateCard(card *model.Card, publicKey string, hmacSecret []byte) error
	GetAllCards() ([]model.Card, error)
	GetCardByID(id uint) (*model.Card, error)
	UpdateCard(card *model.Card) error
	DeleteCard(id uint) error
}

type cardRepo struct {
	BaseRepository
}

func CardRepositoryInstance(db *gorm.DB) CardRepository {
	return &cardRepo{
		BaseRepository: InitializeRepository(db),
	}
}

// CreateCard implements CardRepository.
func (c *cardRepo) CreateCard(card *model.Card, publicKey string, hmacSecret []byte) error {
	// Шифрование номера карты и срока действия
	fmt.Printf("PublicKey: %s\n", publicKey)
	encryptedNumber, err := security.EncryptData(card.Number, publicKey)
	fmt.Printf("Encrypted Number: %s\n", encryptedNumber)
	if err != nil {
		return err
	}
	encryptedExpiryDate, err := security.EncryptData(card.ExpiryDate, publicKey)
	if err != nil {
		return err
	}

	// Генерация HMAC для CVV
	card.CVV = security.GenerateHMAC(card.CVV, hmacSecret)

	// Сохранение зашифрованных данных в структуру
	card.Number = encryptedNumber
	card.ExpiryDate = encryptedExpiryDate
	card.CreatedAt = time.Now()

	// Сохранение карты в базе данных
	if err := c.db.Create(card).Error; err != nil {
		return err
	}

	return nil
}

// func GenerateHMAC(s string, hmacSecret []byte) string {
// 	panic("unimplemented")
// }

// func EncryptData(s, publicKey string) (any, any) {
// 	panic("unimplemented")
// }

// DeleteCard implements CardRepository.
func (c *cardRepo) DeleteCard(id uint) error {
	panic("unimplemented")
}

// GetAllCards implements CardRepository.
func (c *cardRepo) GetAllCards() ([]model.Card, error) {
	var cards []model.Card
	if err := c.db.Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

// GetCardByID implements CardRepository.
func (c *cardRepo) GetCardByID(id uint) (*model.Card, error) {
	panic("unimplemented")
}

// UpdateCard implements CardRepository.
func (c *cardRepo) UpdateCard(card *model.Card) error {
	panic("unimplemented")
}
