package security

import (
	"FinanceGolang/src/model"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your_secret_key")

var ErrTokenExpired = errors.New("token expired")

func GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	fmt.Printf("Parsing token: %s\n", tokenString)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return nil, err
	}

	if !token.Valid {
		fmt.Println("Token is invalid")
		return nil, fmt.Errorf("invalid token")
	}

	fmt.Printf("Parsed claims: %+v\n", claims)
	return claims, nil
}

func GetToken(c *gin.Context) (string, error) {
	// Попытка извлечь токен из заголовка Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		// Если токен отсутствует в заголовке, ищем его в параметрах запроса
		tokenString = c.Query("token")
		if tokenString == "" {
			return "", fmt.Errorf("token not found in either Authorization header or query parameters")
		}
	}

	// Удаляем префикс "Bearer ", если он есть
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	return tokenString, nil
}

func IsTokenValid(tokenString string) bool {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	// Если ошибка или токен недействителен, возвращаем false
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return false
	}

	return true
}

func CutToken(tokenString string) (string, error) {
	// Удаляем префикс "Bearer ", если он есть
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"status":  "error",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString, _ = CutToken(tokenString)

		log.Printf("Received token: %s", tokenString)

		claims, err := ParseToken(tokenString)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				log.Printf("Token expired: %v", err)
				c.JSON(401, gin.H{"status": "error", "message": "token expired"})
			} else {
				log.Printf("Token parsing error: %v", err)
				c.JSON(401, gin.H{"status": "error", "message": "invalid token"})
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
