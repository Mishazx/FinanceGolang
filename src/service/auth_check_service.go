package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"FinanceGolang/src/security"
	// "FinanceGolang/src/token"
	"fmt"
	"log"
)

type AuthCheckService interface {
	ValidateUserFromToken(tokenString string) (*model.User, error)
}

type authCheckService struct {
	userRepo repository.UserRepository
}

func NewAuthCheckService(userRepo repository.UserRepository) AuthCheckService {
	return &authCheckService{userRepo: userRepo}
}

func (s *authCheckService) ValidateUserFromToken(tokenString string) (*model.User, error) {
	claims, err := security.ParseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	user, err := s.userRepo.FindUserByUsername(claims.Username)
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