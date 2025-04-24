package model

import (
	"time"
)

// Card представляет модель данных банковской карты.
type Card struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	AccountID  int       `json:"account_id"`
	Number     string    `json:"number"`      // Номер карты (зашифрован).
	ExpiryDate string    `json:"expiry_date"` // Срок действия карты (зашифрован).
	CVV        string    `json:"-"`           // CVV код (хеширован).
	CreatedAt  time.Time `json:"created_at"`

	Account Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"account"`
}
