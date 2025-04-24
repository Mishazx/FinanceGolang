package controller

import (
	"github.com/gin-gonic/gin"
)

type AccountController struct{}

func NewAccountController() *AccountController {
	return &AccountController{}
}

func (h *AccountController) CreateAccount(c *gin.Context) {
	// implementation
}

func (h *AccountController) GetAccounts(c *gin.Context) {
	// implementation
}
