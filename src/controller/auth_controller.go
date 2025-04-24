package controller

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (h *AuthController) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	if err := h.authService.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user created successfully",
	})
}

func (h *AuthController) Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthController) MyUser(c *gin.Context) {
	username := c.MustGet("username").(string)

	// Получите информацию о пользователе из базы данных или другого источника
	user, err := h.authService.GetUserByUsernameWithoutPassword(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hello " + username, "user": user})
}

func (h *AuthController) AuthStatus(c *gin.Context) {
	token := c.GetHeader("Authorization")

	isValid, err := h.authService.AuthStatus(token)
	if err != nil || !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "success",
			"message": "invalid token",
			"isValid": isValid,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "token is valid",
		"isValid": isValid,
	})
}
