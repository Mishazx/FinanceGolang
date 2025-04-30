package controller

import (
	"FinanceGolang/src/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	analyticsService *service.AnalyticsService
}

func CreateAnalyticsController(analyticsService *service.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}

// GetAnalytics возвращает статистику доходов и расходов
func (c *AnalyticsController) GetAnalytics(ctx *gin.Context) {
	accountID := ctx.GetUint("account_id")
	if accountID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан ID счета"})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания"})
		return
	}

	stats, err := c.analyticsService.GetIncomeExpenseStats(accountID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetBalanceForecast возвращает прогноз баланса
func (c *AnalyticsController) GetBalanceForecast(ctx *gin.Context) {
	accountID := ctx.GetUint("account_id")
	if accountID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан ID счета"})
		return
	}

	months := ctx.GetInt("months")
	if months <= 0 {
		months = 6 // По умолчанию прогноз на 6 месяцев
	}

	forecast, err := c.analyticsService.GetBalanceForecast(accountID, months)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forecast)
}

// GetSpendingCategories возвращает статистику по категориям расходов
func (c *AnalyticsController) GetSpendingCategories(ctx *gin.Context) {
	accountID := ctx.GetUint("account_id")
	if accountID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Не указан ID счета"})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания"})
		return
	}

	categories, err := c.analyticsService.GetSpendingCategories(accountID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
