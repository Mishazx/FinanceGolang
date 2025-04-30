package controller

import (
	"FinanceGolang/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreditController struct {
	creditService service.CreditService
}

func CreateCreditController(creditService service.CreditService) *CreditController {
	return &CreditController{creditService: creditService}
}

type CreateCreditRequest struct {
	AccountID   uint    `json:"account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	TermMonths  int     `json:"term_months" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type ProcessPaymentRequest struct {
	PaymentNumber int `json:"payment_number" binding:"required,gt=0"`
}

func (c *CreditController) CreateCredit(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	var req CreateCreditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Дополнительная обработка - убеждаемся, что description всегда строка
	description := req.Description
	if description == "" {
		description = "Потребительский кредит"
	}

	credit, err := c.creditService.CreateCredit(
		userID.(uint),
		req.AccountID,
		req.Amount,
		req.TermMonths,
		description,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "credit created successfully",
		"credit":  credit,
	})
}

func (c *CreditController) GetCreditByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit ID"})
		return
	}

	credit, err := c.creditService.GetCreditByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"credit": credit})
}

func (c *CreditController) GetUserCredits(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	credits, err := c.creditService.GetUserCredits(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"credits": credits})
}

func (c *CreditController) GetPaymentSchedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit ID"})
		return
	}

	schedule, err := c.creditService.GetPaymentSchedule(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

func (c *CreditController) ProcessPayment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit ID"})
		return
	}

	var req ProcessPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.creditService.ProcessPayment(uint(id), req.PaymentNumber); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "payment processed successfully"})
}
