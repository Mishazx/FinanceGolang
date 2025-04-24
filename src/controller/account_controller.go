package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.accountService.CreateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "account created successfully"})
}
func (h *AccountController) GetAccountByID(c *gin.Context) {
	// implementation
}

func (h *AccountController) GetAccounts(c *gin.Context) {
	accounts, err := h.accountService.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
