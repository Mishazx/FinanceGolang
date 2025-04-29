package repository

import (
	"FinanceGolang/src/model"

	"gorm.io/gorm"
	// "errors"
)

type RoleRepository interface {
	CreateRole(role *model.Role) error
	GetRoleByName(name string) (*model.Role, error)
	AssignRoleToUser(userID uint, roleName string) error
	GetUserRoles(userID uint) ([]model.Role, error)
}

type roleRepository struct {
	BaseRepository
}

func RoleRepositoryInstance(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: InitializeRepository(db),
	}
}

func (r *roleRepository) CreateRole(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetRoleByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) AssignRoleToUser(userID uint, roleName string) error {
	role, err := r.GetRoleByName(roleName)
	if err != nil {
		return err
	}

	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&user).Association("Roles").Append(role)
}

func (r *roleRepository) GetUserRoles(userID uint) ([]model.Role, error) {
	var user model.User
	if err := r.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}
