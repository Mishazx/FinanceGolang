package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsController(analyticsService service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{analyticsService: analyticsService}
}

func (ac *AnalyticsController) GetAnalytics(c *gin.Context) {
	var request model.AnalyticsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Неверный формат запроса",
		})
		return
	}

	analytics, err := ac.analyticsService.GetAnalytics(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка при получении аналитики",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"analytics": analytics,
	})
}

func (ac *AnalyticsController) GetBalanceForecast(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Неверный ID счета",
		})
		return
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil || days < 1 || days > 365 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Неверное количество дней",
		})
		return
	}

	forecast, err := ac.analyticsService.GetBalanceForecast(uint(accountID), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Ошибка при получении прогноза",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"forecast": forecast,
	})
} 