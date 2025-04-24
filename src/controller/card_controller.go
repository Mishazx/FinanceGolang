package controller

import (
	"FinanceGolang/src/service"
)

type CardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) *CardController {
	return &CardController{cardService: cardService}
}
