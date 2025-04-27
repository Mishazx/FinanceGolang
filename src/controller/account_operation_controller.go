package controller

import (
	"FinanceGolang/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountOperationController struct {
	accountOperationService service.AccountOperationService
}

func NewAccountOperationController(accountOperationService service.AccountOperationService) *AccountOperationController {
	return &AccountOperationController{
		accountOperationService: accountOperationService,
	}
}

type DepositRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type WithdrawRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type TransferRequest struct {
	ToAccountID uint    `json:"to_account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

func (c *AccountOperationController) Deposit(ctx *gin.Context) {
	accountIDStr := ctx.Param("id")
	if accountIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var req DepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.accountOperationService.Deposit(uint(accountID), req.Amount, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deposit successful"})
}

func (c *AccountOperationController) Withdraw(ctx *gin.Context) {
	accountIDStr := ctx.Param("id")
	if accountIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var req WithdrawRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.accountOperationService.Withdraw(uint(accountID), req.Amount, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "withdrawal successful"})
}

func (c *AccountOperationController) Transfer(ctx *gin.Context) {
	fromAccountIDStr := ctx.Param("id")
	if fromAccountIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	fromAccountID, err := strconv.ParseUint(fromAccountIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.accountOperationService.Transfer(uint(fromAccountID), req.ToAccountID, req.Amount, req.Description); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}

func (c *AccountOperationController) GetTransactions(ctx *gin.Context) {
	accountIDStr := ctx.Param("id")
	if accountIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account ID is required"})
		return
	}

	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	transactions, err := c.accountOperationService.GetTransactions(uint(accountID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"transactions": transactions})
} 