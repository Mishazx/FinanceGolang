package repository

import (
	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

type CardRepository interface {
	CreateCard(card *model.Card) error
	GetAllCards() ([]model.Card, error)
	GetCardByID(id uint) (*model.Card, error)
	UpdateCard(card *model.Card) error
	DeleteCard(id uint) error
}

type cardRepository struct {
	db *gorm.DB
}
