package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
	// "FinanceGolang/src/repository"
	// "FinanceGolang/src/database"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) *CardController {
	return &CardController{cardService: cardService}
}

func (cc *CardController) CreateCard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "user not found"})
		return
	}

	var card model.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	unsecureCard, err := cc.cardService.CreateCard(&card, userID.(uint))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "account does not belong to the user" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "card created successfully",
		"card":    unsecureCard,
	})
}

func (cc *CardController) GetCardByID(router *gin.Context) {

}

func (cc *CardController) GetAllCards(c *gin.Context) {
	cards, err := cc.cardService.GetAllCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if len(cards) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "success",
			"message": "no cards found",
			"cards":   cards,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"cards":  cards,
	})
}
