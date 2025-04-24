package model

import (
	"time"
)

type Account struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	// ExpiredAt time.Time `json:"expired_at"`

	Cards []Card `gorm:"foreignKey:AccountID" json:"cards"`
}
