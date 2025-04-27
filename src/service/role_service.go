package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	// "errors"
)

type RoleService interface {
	CreateRole(name, description string) error
	AssignRoleToUser(userID uint, roleName string) error
	GetUserRoles(userID uint) ([]model.Role, error)
	InitializeDefaultRoles() error
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) CreateRole(name, description string) error {
	role := &model.Role{
		Name:        name,
		Description: description,
	}
	return s.roleRepo.CreateRole(role)
}

func (s *roleService) AssignRoleToUser(userID uint, roleName string) error {
	return s.roleRepo.AssignRoleToUser(userID, roleName)
}

func (s *roleService) GetUserRoles(userID uint) ([]model.Role, error) {
	return s.roleRepo.GetUserRoles(userID)
}

func (s *roleService) InitializeDefaultRoles() error {
	// Создаем стандартные роли, если они не существуют
	roles := []struct {
		name        string
		description string
	}{
		{model.RoleAdmin, "Администратор системы"},
		{model.RoleUser, "Обычный пользователь"},
		{model.RoleManager, "Менеджер"},
	}

	for _, role := range roles {
		_, err := s.roleRepo.GetRoleByName(role.name)
		if err != nil {
			// Если роль не найдена, создаем её
			if err := s.CreateRole(role.name, role.description); err != nil {
				return err
			}
		}
	}

	return nil
} 