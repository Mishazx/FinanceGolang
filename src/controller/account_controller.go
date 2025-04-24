package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

func (h *AccountController) CreateAccount(c *gin.Context) {
	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if err := h.accountService.CreateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": err.Error(),
			"error":  "could not create account",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "account created successfully",
	})
}
func (h *AccountController) GetAccountByUserID(c *gin.Context) {
	userID, exists := c.MustGet("userID").(uint)

	fmt.Println("UserID from context:", userID)

	fmt.Println("UserID exists:", exists)

	account, err := h.accountService.GetAccountByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   "success",
			"message":  "no accounts found",
			"accounts": account,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"accounts": account,
	})
}

func (h *AccountController) GetAccountsAll(c *gin.Context) {
	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	if len(accounts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":   "success",
			"message":  "no accounts found",
			"accounts": accounts,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"accounts": accounts,
	})
}
