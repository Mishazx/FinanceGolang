package service

import (
	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	// "FinanceGolang/src/token"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(user *model.User) error
	Login(user *model.User) (string, error)
	AuthStatus(token string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error)
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
	return security.GenerateToken(foundUser.ID, foundUser.Username, foundUser.Email)
}

func (s *authService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) GetUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error) {
	userResponse, err := s.userRepo.FindUserByUsernameWithoutPassword(username)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (s *authService) AuthStatus(token string) (bool, error) {
	token, _ = security.CutToken(token)

	claims, err := security.ParseToken(token)
	if err != nil {
		return false, fmt.Errorf("invalid token")
	}

	// Проверяем существование пользователя
	user, err := s.userRepo.FindUserByUsername(claims.Username)
	if err != nil {
		return false, fmt.Errorf("No valid auth token, user not found")
	}

	// Проверяем соответствие ID
	if user.ID != claims.UserID {
		return false, fmt.Errorf("No valid auth token, invalid user")
	}

	return true, nil
}
