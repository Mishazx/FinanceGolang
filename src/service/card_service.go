package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	// "FinanceGolang/src/"
)

type CardService interface {
	CreateCard(card *model.Card) error
	GetCardByID(id uint) (*model.Card, error)
	GetAllCards() ([]model.Card, error)
}

type cardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(cardRepo repository.CardRepository) CardService {
	return &cardService{cardRepo: cardRepo}
}

func (s *cardService) CreateCard(card *model.Card) error {
	return s.cardRepo.CreateCard(card)
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
