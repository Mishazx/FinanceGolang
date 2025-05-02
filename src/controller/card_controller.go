package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"

	// "FinanceGolang/src/repository"
	// "FinanceGolang/src/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CardController struct {
	cardService service.CardService
}

func CreateCardController(cardService service.CardService) *CardController {
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
			"message": "invalid request body",
		})
		return
	}

	if card.AccountID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "account_id is required",
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
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "user not found",
		})
		return
	}

	cards, err := cc.cardService.GetUserCards(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Преобразуем карты в DTO
	var cardDTOs []map[string]interface{}
	for _, card := range cards {
		dto := card.ToDTO()
		// Логируем каждую карту для отладки
		fmt.Printf("Card DTO: %+v\n", dto)
		cardDTOs = append(cardDTOs, dto)
	}

	response := gin.H{
		"status": "success",
		"cards":  cardDTOs,
	}
	// Логируем финальный ответ
	fmt.Printf("Response: %+v\n", response)

	c.JSON(http.StatusOK, response)
}
