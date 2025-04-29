package model

import (
	"fmt"
	"regexp"
	"time"
)

// Card представляет модель данных банковской карты.
type Card struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	AccountID int    `json:"account_id" gorm:"index:idx_card_account_id"`
	Number    string `json:"number" gorm:"check:number ~ '^[0-9]{16}$'"` // Номер карты (зашифрован)
	// ExpiryDate string    `json:"expiry_date" gorm:"check:expiry_date ~ '^(0[1-9]|1[0-2])\/([0-9]{2})$'"` // Срок действия карты (зашифрован)
	ExpiryDate string    `json:"expiry_date" gorm:"type:varchar(5);check:expiry_date ~ '^(0[1-9]|1[0-2])\\/([0-9]{2})$'"` // Срок действия карты (зашифрован)
	CVV        string    `json:"-" gorm:"check:length(cvv) = 3"`                                                          // CVV код (хеширован)
	HMAC       string    `json:"hmac" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`

	Account Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account"`
}

// ValidateCardNumber проверяет корректность номера карты
func (c *Card) ValidateCardNumber() bool {
	cardRegex := regexp.MustCompile(`^[0-9]{16}$`)
	return cardRegex.MatchString(c.Number)
}

// ValidateExpiryDate проверяет корректность срока действия
func (c *Card) ValidateExpiryDate() bool {
	expiryRegex := regexp.MustCompile(`^(0[1-9]|1[0-2])\/([0-9]{2})$`)
	if !expiryRegex.MatchString(c.ExpiryDate) {
		return false
	}

	// Проверяем, что срок действия не истек
	now := time.Now()
	month := now.Month()
	year := now.Year() % 100

	expiryMonth := c.ExpiryDate[:2]
	expiryYear := c.ExpiryDate[3:]

	if expiryYear < fmt.Sprintf("%02d", year) {
		return false
	}
	if expiryYear == fmt.Sprintf("%02d", year) && expiryMonth < fmt.Sprintf("%02d", month) {
		return false
	}

	return true
}

// ValidateCVV проверяет корректность CVV
func (c *Card) ValidateCVV() bool {
	cvvRegex := regexp.MustCompile(`^[0-9]{3}$`)
	return cvvRegex.MatchString(c.CVV)
}
