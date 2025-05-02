package service

import (
	"FinanceGolang/src/model"
	"FinanceGolang/src/repository"
	"context"
	"errors"
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

func RoleServiceInstance(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) CreateRole(name, description string) error {
	role := &model.Role{
		Name:        name,
		Description: description,
	}
	return s.roleRepo.Create(context.Background(), role)
}

func (s *roleService) AssignRoleToUser(userID uint, roleName string) error {
	role, err := s.roleRepo.GetByName(context.Background(), roleName)
	if err != nil {
		return err
	}
	return s.roleRepo.UpdatePermissions(context.Background(), role.ID, []string{roleName})
}

func (s *roleService) GetUserRoles(userID uint) ([]model.Role, error) {
	return s.roleRepo.GetByUserID(context.Background(), userID)
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
		_, err := s.roleRepo.GetByName(context.Background(), role.name)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				// Если роль не найдена, создаем её
				if err := s.CreateRole(role.name, role.description); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
