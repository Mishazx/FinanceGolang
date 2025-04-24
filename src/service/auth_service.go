package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(user *model.User) error
	Login(user *model.User) (string, error)
	AuthStatus(token string) (bool, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(user *model.User) error {
	// Проверяем, существует ли пользователь с таким именем
	_, err := s.userRepo.FindUserByUsername(user.Username)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// Возвращаем ошибку, если это не ошибка "запись не найдена"
			return err
		}
	} else {
		// Если пользователь найден, возвращаем ошибку
		return fmt.Errorf("user with username %s already exists", user.Username)
	}
	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Сохраняем пользователя в базе данных
	return s.userRepo.CreateUser(user)
}

func (s *authService) Login(user *model.User) (string, error) {
	// Ищем пользователя по username
	foundUser, err := s.userRepo.FindUserByUsername(user.Username)
	if err != nil {
		return "", err
	}

	// Сравниваем хэши паролей
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return "", err
	}

	// Генерируем токен
	return security.GenerateToken(foundUser.Username)
}

func (s *authService) AuthStatus(token string) (bool, error) {
	token, _ = security.CutToken(token)

	IsTokenValid := security.IsTokenValid(token)

	if !IsTokenValid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
