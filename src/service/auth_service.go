package service

import (
	"FinanceGolang/src/database"
	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"

	"context"
	"fmt"
	"log"
)

type AuthService interface {
	Register(user *model.User) error
	Login(user *model.User) (string, error)
	AuthStatus(token string) (bool, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error)
	ValidateUserFromToken(tokenString string) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func AuthServiceInstance(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: repository.RoleRepositoryInstance(database.DB),
	}
}

func (s *authService) Register(user *model.User) error {
	// Проверяем, существует ли пользователь с таким именем
	existingUser, err := s.userRepo.GetByUsername(context.Background(), user.Username)
	if err != nil {
		return fmt.Errorf("error checking username: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with username %s already exists", user.Username)
	}

	// Проверяем валидацию
	if err := user.Validate(); err != nil {
		switch err {
		case model.ErrInvalidEmail:
			return fmt.Errorf("invalid email format")
		case model.ErrInvalidUsername:
			return fmt.Errorf("username must be 3-50 characters long and contain only letters, numbers and underscores")
		case model.ErrInvalidPassword:
			return fmt.Errorf("password must be at least 8 characters long")
		case model.ErrInvalidFIO:
			return fmt.Errorf("FIO must be 3-100 characters long and contain only letters, spaces and hyphens")
		default:
			return fmt.Errorf("validation error: %v", err)
		}
	}

	// Сохраняем пользователя в базе данных
	if err := s.userRepo.Create(context.Background(), user); err != nil {
		return err
	}

	// Назначаем роль user новому пользователю
	role, err := s.roleRepo.GetByName(context.Background(), model.RoleUser)
	if err != nil {
		return fmt.Errorf("failed to get default role: %v", err)
	}

	// Создаем связь пользователя с ролью
	if err := s.userRepo.AddRole(context.Background(), user.ID, role.ID); err != nil {
		return fmt.Errorf("failed to assign default role: %v", err)
	}

	return nil
}

func (s *authService) Login(user *model.User) (string, error) {
	// Ищем пользователя по username
	foundUser, err := s.userRepo.GetByUsername(context.Background(), user.Username)
	if err != nil {
		return "", err
	}

	// Проверяем пароль
	if err := foundUser.CheckPassword(user.Password); err != nil {
		return "", err
	}

	// Генерируем токен
	return security.GenerateToken(foundUser.ID, foundUser.Username, foundUser.Email)
}

func (s *authService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) GetUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) AuthStatus(token string) (bool, error) {
	token, _ = security.CutToken(token)

	claims, err := security.ParseToken(token)
	if err != nil {
		return false, fmt.Errorf("invalid token")
	}

	// Проверяем существование пользователя
	user, err := s.userRepo.GetByUsername(context.Background(), claims.Username)
	if err != nil {
		return false, fmt.Errorf("No valid auth token, user not found")
	}

	// Проверяем соответствие ID
	if user.ID != claims.UserID {
		return false, fmt.Errorf("No valid auth token, invalid user")
	}

	return true, nil
}

func (s *authService) ValidateUserFromToken(tokenString string) (*model.User, error) {
	claims, err := security.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	user, err := s.userRepo.GetByUsername(context.Background(), claims.Username)
	if err != nil {
		log.Printf("No valid auth token, user not found: %v", err)
		return nil, fmt.Errorf("No valid auth token, user not found")
	}

	if user.ID != claims.UserID {
		log.Printf("No valid auth token, invalid user: %v", err)
		return nil, fmt.Errorf("No valid auth token, invalid user")
	}

	return user, nil
}
