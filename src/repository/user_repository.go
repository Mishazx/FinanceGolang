package repository

import (
	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindUserByUsername(username string) (*model.User, error)
	FindUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindUserByUsername(username string) (*model.User, error) {
	var user model.User

	// if err := r.db.Select("id, username").Where("username = ?", username).First(&user).Error; err != nil {
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println("Error finding user:", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error) {
	var user model.User
	// Явно указываем только те поля, которые хотим получить
	if err := r.db.Select("id, username, email, created_at").Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println("Error finding user:", err)
		return nil, err
	}
	// return &user, nil
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
