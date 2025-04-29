package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"errors"
)

type UserService interface {
	GetUserByID(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func UserServiceInstance(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) UpdateUser(user *model.User) error {
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}
	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(id uint) error {
	existingUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}
	return s.userRepo.DeleteUser(id)
}
