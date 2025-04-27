package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
	"fmt"
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
	// Печатаем все поля из контекста

	fmt.Println("START !!! ------------")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "user not found"})
		return
	}

	fmt.Println("userID: ", userID)

	fmt.Println("Request URL:", c.Request.URL.String())
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	fmt.Println("END !!! ------------")

	var card model.Card
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	unsecureCard, err := cc.cardService.CreateCard(&card)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  err.Error(),
			"message": "could not create card",
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
