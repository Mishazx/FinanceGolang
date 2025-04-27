package controller

import (
	"FinanceGolang/src/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CbrController struct {
	cbrService service.CbrService
}

func NewCbrController(cbrService service.CbrService) *CbrController {
	return &CbrController{cbrService: cbrService}
}

func (cc *CbrController) GetKeyRate(c *gin.Context) {
	keyRate, err := cc.cbrService.GetLastKeyRate()
	fmt.Println("keyRate: ", keyRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keyRate": keyRate})
}
