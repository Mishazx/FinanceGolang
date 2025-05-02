package repository

import (
	"context"

	"FinanceGolang/src/dto"
	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

// UserRepository интерфейс репозитория пользователей
type UserRepository interface {
	Repository[model.User]
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetWithRoles(ctx context.Context, id uint) (*model.User, error)
	UpdatePassword(ctx context.Context, id uint, password string) error
	AddRole(ctx context.Context, userID, roleID uint) error
	RemoveRole(ctx context.Context, userID, roleID uint) error
	GetByRole(ctx context.Context, roleName string, offset, limit int) ([]model.User, error)
}

// userRepository реализация репозитория пользователей
type userRepository struct {
	*BaseRepository[model.User]
}

// UserRepositoryInstance создает новый репозиторий пользователей
func UserRepositoryInstance(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository[model.User](db),
	}
}

// Create создает нового пользователя
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := user.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Create(user).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// GetByID получает пользователя по ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &user, nil
}

// GetByUsername получает пользователя по username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.HandleError(err)
	}
	return &user, nil
}

// GetWithRoles получает пользователя с ролями
func (r *userRepository) GetWithRoles(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Roles").First(&user, id).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return &user, nil
}

// Update обновляет пользователя
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := user.Validate(); err != nil {
			return ErrInvalidData
		}

		if err := tx.Save(user).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// UpdatePassword обновляет пароль пользователя
func (r *userRepository) UpdatePassword(ctx context.Context, id uint, password string) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// Delete удаляет пользователя
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Delete(&model.User{}, id).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// List получает список пользователей
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]model.User, error) {
	var users []model.User
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return users, nil
}

// Count возвращает количество пользователей
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, r.HandleError(err)
	}
	return count, nil
}

// AddRole добавляет роль пользователю
func (r *userRepository) AddRole(ctx context.Context, userID, roleID uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		userRole := model.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		if err := tx.Create(&userRole).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// RemoveRole удаляет роль у пользователя
func (r *userRepository) RemoveRole(ctx context.Context, userID, roleID uint) error {
	return r.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{}).Error; err != nil {
			return r.HandleError(err)
		}
		return nil
	})
}

// GetByRole получает пользователей по роли
func (r *userRepository) GetByRole(ctx context.Context, roleName string, offset, limit int) ([]model.User, error) {
	var users []model.User
	if err := r.db.Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", roleName).
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, r.HandleError(err)
	}
	return users, nil
}

// FindUserByUsername получает пользователя по username (устаревший метод)
func (r *userRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByUsernameWithoutPassword получает пользователя по username без пароля
func (r *userRepository) FindUserByUsernameWithoutPassword(username string) (*dto.UserResponse, error) {
	var user model.User
	if err := r.db.Select("id, username, email, created_at").Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
